package gsplat

import (
	"bytes"
	"container/heap"
	"gsbox/cmn"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"sync"
)

func ReWriteShByKmeans(rows []*SplatData) (shN_centroids []uint8, shN_labels []uint8) {
	shDegreeOutput := GetArgShDegree()
	if shDegreeOutput == 0 {
		return // 不输出球谐系数时跳过
	}

	dataCnt := len(rows)
	var nSh45 [][]float32
	for n := range dataCnt {
		nSh45 = append(nSh45, GetSh45Float32ForKmeans(rows[n]))
	}
	dims := []int{0, 9, 24, 45}
	palettes, indexes, counts := kmeansSh45(nSh45, dims[min(shDegreeFrom, shDegreeOutput)], oArg.KI, oArg.KN)
	palettes, indexes = sortCentroidsByCounts(palettes, indexes, counts) // 按质心点数量倒序排序,提高压缩效果稳定输出

	if IsOutputSpz() {
		for n := range dataCnt {
			data := rows[n]
			shs := palettes[indexes[n]]
			switch shDegreeOutput {
			case 1:
				data.SH1 = shs[:9]
				data.SH2 = []uint8(nil)
				data.SH3 = []uint8(nil)
			case 2:
				data.SH1 = []uint8(nil)
				data.SH2 = shs[:24]
				data.SH3 = []uint8(nil)
			case 3:
				data.SH1 = []uint8(nil)
				data.SH2 = shs[:24]
				data.SH3 = shs[24:]
			default:
				data.SH1 = []uint8(nil)
				data.SH2 = []uint8(nil)
				data.SH3 = []uint8(nil)
			}
		}
		return
	}

	paletteSize := len(palettes)
	shN_centroids = make([]uint8, paletteSize*15*4)
	for i := range paletteSize {
		shs := palettes[i]
		for k := range 15 {
			shN_centroids[i*15*4+k*4+0] = shs[k*3+0]
			shN_centroids[i*15*4+k*4+1] = shs[k*3+1]
			shN_centroids[i*15*4+k*4+2] = shs[k*3+2]
			shN_centroids[i*15*4+k*4+3] = 255
		}
	}

	if IsOutputSog() {
		w, h := cmn.ComputeWidthHeight(dataCnt)
		labelsPixCnt := w * h
		shN_labels = make([]uint8, labelsPixCnt*4)
		maxDataIdx := dataCnt - 1
		for i := range labelsPixCnt {
			idx := indexes[min(i, maxDataIdx)]
			shN_labels[i*4+0] = uint8(idx & 0xFF)
			shN_labels[i*4+1] = uint8(idx >> 8)
			shN_labels[i*4+2] = 0
			shN_labels[i*4+3] = 255
		}
	} else if IsOutputSpx() {
		for i, d := range rows {
			d.PaletteIdx = uint16(indexes[i])
		}
	}

	return
}

func sortCentroidsByCounts(centroids [][]uint8, indexes []int32, counts []int32) (sortedCentroids [][]uint8, sortedIndexes []int32) {
	type centroidInfo struct {
		idx   int32
		count int32
	}
	centroidMap := make([]centroidInfo, len(counts))
	for i := range counts {
		centroidMap[i] = centroidInfo{idx: int32(i), count: counts[i]}
	}

	sort.Slice(centroidMap, func(i, j int) bool {
		return centroidMap[i].count > centroidMap[j].count
	})

	sortedCentroids = make([][]uint8, len(centroids))
	sortedIndexes = make([]int32, len(indexes))
	for i, info := range centroidMap {
		sortedCentroids[i] = centroids[info.idx]
		for j, idx := range indexes {
			if idx == info.idx {
				sortedIndexes[j] = int32(i)
			}
		}
	}

	return sortedCentroids, sortedIndexes
}

func kmeansSh45(nSh45 [][]float32, dim int, maxIters int, maxBBFNodes int) (centroids [][]uint8, labels []int32, counts []int32) {
	n := len(nSh45)
	paletteSize := int(math.Min(64, math.Pow(2, math.Floor(math.Log2(float64(n)/1024.0)))) * 1024)
	labels = make([]int32, n)

	// 0 维快速返回
	if dim == 0 {
		centroids = make([][]uint8, paletteSize)
		for i := range centroids {
			centroids[i] = bytes.Repeat([]uint8{128}, 45) // 置零，45个128（浮点数0.0编码后为128）
		}
		return
	}

	// 1. 随机选 k 个中心
	f32Centroids := make([][]float32, paletteSize)
	used := make([]bool, n)
	for i := 0; i < paletteSize; {
		idx := rand.Intn(n)
		if !used[idx] {
			used[idx] = true
			f32Centroids[i] = make([]float32, 45)
			copy(f32Centroids[i], nSh45[idx])
			i++
		}
	}

	for range maxIters {
		// 2. 建 KD-Tree
		tree := buildKDTree(f32Centroids)
		buildCentSoA(f32Centroids)

		// 3. 分配（只算前 dim 维）
		parAssignDim(n, nSh45, labels, tree, dim, maxBBFNodes)

		// 4. 累加新中心（全 45 维，不累加 0 维即可）
		newCents := make([][]float32, paletteSize)
		counts = make([]int32, paletteSize)
		for i := range paletteSize {
			newCents[i] = make([]float32, 45)
		}
		for i := range n {
			c := labels[i]
			for d := range 45 {
				newCents[c][d] += nSh45[i][d]
			}
			counts[c]++
		}
		// 空簇重投
		for c := range paletteSize {
			if counts[c] == 0 {
				copy(newCents[c], nSh45[rand.Intn(n)])
			} else {
				for d := range 45 {
					newCents[c][d] /= float32(counts[c])
				}
			}
		}
		f32Centroids = newCents
	}

	centroids = make([][]uint8, paletteSize)
	for i, v := range f32Centroids {
		centroids[i] = ToSh45(v)
		for j := dim; j < 45; j++ {
			centroids[i][j] = 128 // 超有效维度的部分都置零（浮点数0.0编码后为128）
		}
	}
	return
}

// KD-Tree
type kdNode struct {
	idx   int32
	axis  int32
	left  *kdNode
	right *kdNode
}

type kdTree struct {
	cents [][]float32
	root  *kdNode
}

func buildKDTree(cents [][]float32) *kdTree {
	k := int32(len(cents))
	idxs := make([]int32, k)
	for i := range k {
		idxs[i] = i
	}
	var build func([]int32, int32) *kdNode
	build = func(idxs []int32, depth int32) *kdNode {
		if len(idxs) == 0 {
			return nil
		}
		axis := depth % 45
		sort.Slice(idxs, func(i, j int) bool {
			return cents[idxs[i]][axis] < cents[idxs[j]][axis]
		})
		med := len(idxs) / 2
		return &kdNode{
			idx:   idxs[med],
			axis:  axis,
			left:  build(idxs[:med], depth+1),
			right: build(idxs[med+1:], depth+1),
		}
	}
	return &kdTree{cents: cents, root: build(idxs, 0)}
}

// SoA 中心视图
var (
	centSoA  [45][]float32
	soaReady bool
	soaSize  int
)

func buildCentSoA(cents [][]float32) {
	if len(cents) == 0 {
		return
	}
	soaSize = len(cents)
	for d := range 45 {
		centSoA[d] = make([]float32, soaSize)
		for i := 0; i < soaSize; i++ {
			centSoA[d][i] = cents[i][d]
		}
	}
	soaReady = true
}

// 堆结构
type bbEntry struct {
	node *kdNode
	dist float32
}
type bbHeap []bbEntry

func (h bbHeap) Len() int            { return len(h) }
func (h bbHeap) Less(i, j int) bool  { return h[i].dist < h[j].dist }
func (h bbHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *bbHeap) Push(x interface{}) { *h = append(*h, x.(bbEntry)) }
func (h *bbHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// 维度感知最近邻
func (t *kdTree) NearestBBF(pt []float32, dim int, maxBBFNodes int) int32 {
	var bestIdx int32 = -1
	var bestDist float32 = math.MaxFloat32

	h := &bbHeap{}
	heap.Init(h)
	heap.Push(h, bbEntry{node: t.root, dist: 0})

	visited := 0
	for h.Len() > 0 && visited < maxBBFNodes {
		e := heap.Pop(h).(bbEntry)
		n := e.node
		if n == nil {
			continue
		}
		visited++

		// 防呆：SoA 未就绪或越界 → 回退行主序
		if !soaReady || int(n.idx) >= soaSize {
			cent := t.cents[n.idx]
			var dist float32
			for d := range dim { // 只算前 dim 维
				delta := pt[d] - cent[d]
				dist += delta * delta
			}
			if dist < bestDist {
				bestDist, bestIdx = dist, n.idx
			}
		} else {
			// 正常 SoA 路径：只算前 dim 维
			var distShort float32
			for d := range dim {
				delta := pt[d] - centSoA[d][n.idx]
				distShort += delta * delta
			}
			if distShort < bestDist {
				bestDist, bestIdx = distShort, n.idx
			}
		}

		// 标准 KD 左右子入堆
		axis := n.axis
		diff := pt[axis] - centSoA[axis][n.idx]
		var first, second *kdNode
		if diff < 0 {
			first, second = n.left, n.right
		} else {
			first, second = n.right, n.left
		}
		if first != nil {
			heap.Push(h, bbEntry{node: first, dist: 0})
		}
		if second != nil {
			heap.Push(h, bbEntry{node: second, dist: diff * diff})
		}
	}
	return bestIdx
}

// 并行 assign（维度感知）
func parAssignDim(n int, nSh45 [][]float32, labels []int32, tree *kdTree, dim int, maxBBFNodes int) {
	var wg sync.WaitGroup
	stride := (n + runtime.GOMAXPROCS(0) - 1) / runtime.GOMAXPROCS(0)
	for g := 0; g < runtime.GOMAXPROCS(0); g++ {
		wg.Add(1)
		go func(g int) {
			start := g * stride
			end := min(start+stride, n)
			for i := start; i < end; i++ {
				labels[i] = tree.NearestBBF(nSh45[i], dim, maxBBFNodes)
			}
			wg.Done()
		}(g)
	}
	wg.Wait()
}
