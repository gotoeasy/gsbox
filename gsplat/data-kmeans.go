package gsplat

import (
	"log"
	"math"
	"math/rand"
	"sort"
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
		for i := range n {
			labels[i] = tree.nearest(nSh45[i])
		}

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

func (t *kdTree) nearest(pt []float32) int32 {
	var bestIdx int32 = -1
	var bestDist float32 = math.MaxFloat32
	var dfs func(*kdNode, int32)
	dfs = func(n *kdNode, depth int32) {
		if n == nil {
			return
		}
		cent := t.cents[n.idx]
		dist := float32(0)
		for d := range 45 {
			delta := pt[d] - cent[d]
			dist += delta * delta
		}
		if dist < bestDist {
			bestDist, bestIdx = dist, n.idx
		}
		axis := depth % 45
		diff := pt[axis] - cent[axis]
		var first, second *kdNode
		if diff < 0 {
			first, second = n.left, n.right
		} else {
			first, second = n.right, n.left
		}
		dfs(first, depth+1)
		if diff*diff < bestDist {
			dfs(second, depth+1)
		}
	}
	dfs(t.root, 0)
	return bestIdx
}
