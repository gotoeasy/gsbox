package gsplat

import (
	"container/heap"
	"log"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"sync"
	"time"
)

func KmeansSh45(nSh45 [][]float32) (centroids [][]uint8, labels []int32) {
	n := len(nSh45)
	paletteSize := int(math.Min(64, math.Pow(2, math.Floor(math.Log2(float64(n)/1024.0)))) * 1024)
	const maxIters = 10 // 固定10次迭代

	labels = make([]int32, n)

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

	for iter := range maxIters {
		log.Println("iters ", iter+1)
		// 2. 建 KD-Tree
		tree := buildKDTree(f32Centroids)
		buildCentSoA(f32Centroids)

		// 3. 分配
		startTime := time.Now().UnixMilli()
		parAssign(n, nSh45, labels, tree)
		log.Println("tree.nearest 耗时", (time.Now().UnixMilli() - startTime), "MS")

		// 4. 累加新中心
		newCents := make([][]float32, paletteSize)
		counts := make([]int32, paletteSize)
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

	centroids = make([][]uint8, len(f32Centroids))
	for i, v := range f32Centroids {
		centroids[i] = ToSh45(v)
	}
	return
}

/* ---------- KD-Tree ---------- */
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

const maxBBFNodes = 15 // 越小误差越大速度越快，酌情调整

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
	for d := 0; d < 45; d++ {
		centSoA[d] = make([]float32, soaSize)
		for i := 0; i < soaSize; i++ {
			centSoA[d][i] = cents[i][d]
		}
	}
	soaReady = true
}

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

// 近似最近邻（SIMD 友好加载）
func (t *kdTree) NearestBBF(pt []float32) int32 {
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
			for d := 0; d < 45; d++ {
				delta := pt[d] - cent[d]
				dist += delta * delta
			}
			if dist < bestDist {
				bestDist, bestIdx = dist, n.idx
			}
		} else {
			// 正常 SoA 路径
			var dist9 float32
			// 1) 前 9 维：8×1 + 1
			for d := 0; d < 8; d++ {
				delta := pt[d] - centSoA[d][n.idx]
				dist9 += delta * delta
			}
			delta := pt[8] - centSoA[8][n.idx]
			dist9 += delta * delta
			if dist9 < bestDist {
				// 2) 补后 36 维：4×8 + 4
				for d := 9; d < 45; d++ {
					delta := pt[d] - centSoA[d][n.idx]
					dist9 += delta * delta
				}
				if dist9 < bestDist {
					bestDist, bestIdx = dist9, n.idx
				}
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

// 并行 assign
func parAssign(n int, nSh45 [][]float32, labels []int32, tree *kdTree) {
	var wg sync.WaitGroup
	stride := (n + runtime.GOMAXPROCS(0) - 1) / runtime.GOMAXPROCS(0)
	for g := 0; g < runtime.GOMAXPROCS(0); g++ {
		wg.Add(1)
		go func(g int) {
			start := g * stride
			end := start + stride
			if end > n {
				end = n
			}
			for i := start; i < end; i++ {
				labels[i] = tree.NearestBBF(nSh45[i])
			}
			wg.Done()
		}(g)
	}
	wg.Wait()
}
