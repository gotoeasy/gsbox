package gsplat

import (
	"log"
	"math"
	"math/rand"
	"sort"
)

const dim = 45

// Kmeans45 对 n×45 维数据做 KD-Tree 加速 Lloyd 迭代
// nSh45 长度必须 > k，且每行 45 元素；返回 k×45 中心矩阵 和 n×1 标签切片
func KmeansSh45(nSh45 [][]float32, k int, iters int) (centroids [][]uint8, labels []int32) {
	n := len(nSh45)
	labels = make([]int32, n)

	// 1. 随机选 k 个中心
	f32Centroids := make([][]float32, k)
	used := make([]bool, n)
	for i := 0; i < k; {
		idx := rand.Intn(n)
		if !used[idx] {
			used[idx] = true
			f32Centroids[i] = make([]float32, dim)
			copy(f32Centroids[i], nSh45[idx])
			i++
		}
	}

	for iter := range iters {
		log.Println("iters ", iter+1)
		// 2. 建 KD-Tree
		tree := buildKDTree(f32Centroids)

		// 3. 分配
		for i := range n {
			labels[i] = tree.nearest(nSh45[i])
		}

		// 4. 累加新中心
		newCents := make([][]float32, k)
		counts := make([]int32, k)
		for i := range k {
			newCents[i] = make([]float32, dim)
		}
		for i := range n {
			c := labels[i]
			for d := 0; d < dim; d++ {
				newCents[c][d] += nSh45[i][d]
			}
			counts[c]++
		}
		// 空簇重投
		for c := range k {
			if counts[c] == 0 {
				copy(newCents[c], nSh45[rand.Intn(n)])
			} else {
				for d := 0; d < dim; d++ {
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
	for i := int32(0); i < k; i++ {
		idxs[i] = i
	}
	var build func([]int32, int32) *kdNode
	build = func(idxs []int32, depth int32) *kdNode {
		if len(idxs) == 0 {
			return nil
		}
		axis := depth % dim
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
		for d := 0; d < dim; d++ {
			delta := pt[d] - cent[d]
			dist += delta * delta
		}
		if dist < bestDist {
			bestDist, bestIdx = dist, n.idx
		}
		axis := depth % dim
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
