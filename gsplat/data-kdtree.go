package gsplat

import (
	"github.com/kyroy/kdtree"
)

type KdPoint struct {
	Position [3]float64
	Data     *SplatData
}

// Dimensions 返回点的维度
func (p KdPoint) Dimensions() int {
	return 3
}

// Dimension 返回点的第 d 维的值
func (p KdPoint) Dimension(i int) float64 {
	return p.Position[i]
}

func SortKdTree(datas []*SplatData) {

	tree := kdtree.New([]kdtree.Point{})

	for _, data := range datas {
		tree.Insert(KdPoint{
			Position: [3]float64{float64(data.PositionX), float64(data.PositionY), float64(data.PositionZ)},
			Data:     data,
		})
	}

	ps := tree.Points()
	for i, p := range ps {
		datas[i] = p.(KdPoint).Data
	}
}
