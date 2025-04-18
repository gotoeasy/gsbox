package gsplat

import (
	"fmt"
	"gsbox/cmn"
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
	RotationX uint8
	RotationY uint8
	RotationZ uint8
	RotationW uint8
	SH1       []uint8
	SH2       []uint8
	SH3       []uint8
}

func (s *SplatData) ToBytes() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, cmn.Float32ToBytes(s.PositionX)...)
	bytes = append(bytes, cmn.Float32ToBytes(s.PositionY)...)
	bytes = append(bytes, cmn.Float32ToBytes(s.PositionZ)...)
	bytes = append(bytes, cmn.Float32ToBytes(s.ScaleX)...)
	bytes = append(bytes, cmn.Float32ToBytes(s.ScaleY)...)
	bytes = append(bytes, cmn.Float32ToBytes(s.ScaleZ)...)
	bytes = append(bytes, s.ColorR, s.ColorG, s.ColorB, s.ColorA)
	bytes = append(bytes, s.RotationX, s.RotationY, s.RotationZ, s.RotationW)
	return bytes
}

func (s *SplatData) ToString() string {
	return fmt.Sprintf("%v, %v, %v; %v, %v, %v; %v, %v, %v, %v; %v, %v, %v, %v",
		s.PositionX, s.PositionY, s.PositionZ, s.ScaleX, s.ScaleY, s.ScaleZ, s.ColorR, s.ColorG, s.ColorB, s.ColorA, s.RotationX, s.RotationY, s.RotationZ, s.RotationW)
}
