package gsplat

import "gsbox/cmn"

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
	bytes = append(bytes, s.RotationW, s.RotationX, s.RotationY, s.RotationZ)
	return bytes
}

func (s *SplatData) ToBytesSplat20() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, cmn.EncodeFloat32ToBytes3(s.PositionX)...)
	bytes = append(bytes, cmn.EncodeFloat32ToBytes3(s.PositionY)...)
	bytes = append(bytes, cmn.EncodeFloat32ToBytes3(s.PositionZ)...)
	bytes = append(bytes, cmn.EncodeFloat32ToByte(s.ScaleX))
	bytes = append(bytes, cmn.EncodeFloat32ToByte(s.ScaleY))
	bytes = append(bytes, cmn.EncodeFloat32ToByte(s.ScaleZ))
	bytes = append(bytes, s.ColorR, s.ColorG, s.ColorB, s.ColorA)
	bytes = append(bytes, s.RotationW, s.RotationX, s.RotationY, s.RotationZ)
	return bytes
}
