package gsplat

import (
	"fmt"
	"gsbox/cmn"
	"math"
	"sort"
)

const SPLAT_DATA_SIZE = 3*4 + 3*4 + 4 + 4

type SplatData struct {
	PositionX float32
	PositionY float32
	PositionZ float32
	ScaleX    float32
	ScaleY    float32
	ScaleZ    float32
	ColorR    uint8
	ColorG    uint8
	ColorB    uint8
	ColorA    uint8
	RotationW uint8
	RotationX uint8
	RotationY uint8
	RotationZ uint8
	SH1       []uint8 // sh1 only
	SH2       []uint8 // sh1 + sh2
	SH3       []uint8 // sh3 only
}

func (s *SplatData) ToString() string {
	return fmt.Sprintf("%v, %v, %v; %v, %v, %v; %v, %v, %v, %v; %v, %v, %v, %v",
		s.PositionX, s.PositionY, s.PositionZ, s.ScaleX, s.ScaleY, s.ScaleZ, s.ColorR, s.ColorG, s.ColorB, s.ColorA, s.RotationW, s.RotationX, s.RotationY, s.RotationZ)
}

func Sort(rows []*SplatData) {
	sort.Slice(rows, func(i, j int) bool {
		return math.Exp(float64(cmn.EncodeSplatScale(rows[i].ScaleX)+cmn.EncodeSplatScale(rows[i].ScaleY)+cmn.EncodeSplatScale(rows[i].ScaleZ)))*float64(rows[i].ColorA) <
			math.Exp(float64(cmn.EncodeSplatScale(rows[j].ScaleX)+cmn.EncodeSplatScale(rows[j].ScaleY)+cmn.EncodeSplatScale(rows[j].ScaleZ)))*float64(rows[j].ColorA)
	})
}
