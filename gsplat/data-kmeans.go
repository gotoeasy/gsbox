package gsplat

import (
	"container/heap"
	"log"
	"math"
	"math/rand"
	"sort"
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

		// 3. 分配
		startTime := time.Now().UnixMilli()
		for i := range n {
			labels[i] = tree.NearestBBF(nSh45[i])
		}
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

type bbEntry struct {
	node *kdNode
	dist float32 // 到分割面的距离
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

const maxBBFNodes = 20 // 想再快就调小，想更准就调大

// 近似最近
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

		// 计算到当前节点的真实距离（SIMD 友好版）
		cent := t.cents[n.idx]
		dist := float32(0)

		// 主循环 5×8
		for d := 0; d < 40; d += 8 {
			for k := 0; k < 8; k++ {
				delta := pt[d+k] - cent[d+k]
				dist += delta * delta
			}
		}
		// 尾处理 45-40=5
		for d := 40; d < 45; d++ {
			delta := pt[d] - cent[d]
			dist += delta * delta
		}

		if dist < bestDist {
			bestDist, bestIdx = dist, n.idx
		}

		// 标准 KD-Tree 左右子入堆
		axis := n.axis
		diff := pt[axis] - cent[axis]
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
