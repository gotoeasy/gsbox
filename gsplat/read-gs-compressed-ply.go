package gsplat

import (
	"gsbox/cmn"
	"math"
	"os"
)

func readCompressedPlyDatas(file *os.File, header *PlyHeader, datas []*SplatData) {
	chunks := make([]*Chunk, header.ChunkCount)
	chunkSize := 18 * 4
	for i := 0; i < header.ChunkCount; i++ {
		dataBytes := make([]byte, chunkSize)
		_, err := file.ReadAt(dataBytes, int64(header.HeaderLength+i*chunkSize))
		cmn.ExitOnError(err)

		chunk := &Chunk{}
		chunk.MinX = float64(cmn.BytesToFloat32(dataBytes[0:4]))
		chunk.MinY = float64(cmn.BytesToFloat32(dataBytes[4:8]))
		chunk.MinZ = float64(cmn.BytesToFloat32(dataBytes[8:12]))
		chunk.MaxX = float64(cmn.BytesToFloat32(dataBytes[12:16]))
		chunk.MaxY = float64(cmn.BytesToFloat32(dataBytes[16:20]))
		chunk.MaxZ = float64(cmn.BytesToFloat32(dataBytes[20:24]))
		chunk.MinScaleX = float64(cmn.BytesToFloat32(dataBytes[24:28]))
		chunk.MinScaleY = float64(cmn.BytesToFloat32(dataBytes[28:32]))
		chunk.MinScaleZ = float64(cmn.BytesToFloat32(dataBytes[32:36]))
		chunk.MaxScaleX = float64(cmn.BytesToFloat32(dataBytes[36:40]))
		chunk.MaxScaleY = float64(cmn.BytesToFloat32(dataBytes[40:44]))
		chunk.MaxScaleZ = float64(cmn.BytesToFloat32(dataBytes[44:48]))
		chunk.MinR = float64(cmn.BytesToFloat32(dataBytes[48:52]))
		chunk.MinG = float64(cmn.BytesToFloat32(dataBytes[52:56]))
		chunk.MinB = float64(cmn.BytesToFloat32(dataBytes[56:60]))
		chunk.MaxR = float64(cmn.BytesToFloat32(dataBytes[60:64]))
		chunk.MaxG = float64(cmn.BytesToFloat32(dataBytes[64:68]))
		chunk.MaxB = float64(cmn.BytesToFloat32(dataBytes[68:72]))
		chunks[i] = chunk
	}

	dataCnt := 256
	length := len(chunks)
	offset := header.HeaderLength + header.ChunkCount*chunkSize
	n := 0
	for i := 0; i < length; i++ {
		chunk := chunks[i]
		if i == length-1 {
			dataCnt = header.VertexCount % 256
		}

		dataSize := dataCnt * 16
		dataBytes := make([]byte, dataSize)
		_, err := file.ReadAt(dataBytes, int64(offset))
		cmn.ExitOnError(err)
		offset += dataSize

		for j := 0; j < dataCnt; j++ {
			splat := &SplatData{}
			splat.PositionX, splat.PositionY, splat.PositionZ = unpack111011Xyz(cmn.BytesToUint32(dataBytes[j*16:j*16+4]), chunk)
			splat.RotationW, splat.RotationX, splat.RotationY, splat.RotationZ = unpackRotations(dataBytes[j*16+4 : j*16+8])
			splat.ScaleX, splat.ScaleY, splat.ScaleZ = unpack111011Scale(cmn.BytesToUint32(dataBytes[j*16+8:j*16+12]), chunk)
			splat.ColorR, splat.ColorG, splat.ColorB, splat.ColorA = unpackAndEncodeRgba(cmn.BytesToUint32(dataBytes[j*16+12:j*16+16]), chunk)

			datas[n] = splat
			n++
		}

	}

	shDim := 0
	maxShDegree := header.MaxShDegree()
	switch maxShDegree {
	case 1:
		shDim = 3
	case 2:
		shDim = 8
	case 3:
		shDim = 15
	}

	if shDim > 0 {
		offset = header.HeaderLength + header.ChunkCount*chunkSize + header.VertexCount*16
		for i := 0; i < header.VertexCount; i++ {

			shSize := shDim * 3
			shBytes := make([]byte, shSize)
			_, err := file.ReadAt(shBytes, int64(offset))
			cmn.ExitOnError(err)
			offset += shSize

			shs := make([]byte, 45)
			n := 0
			for j := range shDim {
				for c := range 3 {
					sh := (float64(shBytes[j+c*shDim])/256.0 - 0.5) * 8.0
					shs[n] = cmn.EncodeSplatSH(sh)
					n++
				}
			}
			for ; n < 45; n++ {
				shs[n] = cmn.EncodeSplatSH(0)
			}

			switch maxShDegree {
			case 3:
				datas[i].SH2 = shs[:24]
				datas[i].SH3 = shs[24:]
			case 2:
				datas[i].SH2 = shs[:24]
			case 1:
				datas[i].SH1 = shs[:9]
			}
		}
	}
}

func unpack111011Xyz(packedVal uint32, chunk *Chunk) (float32, float32, float32) {
	// 提取 x, y, z 的值
	x := float64((packedVal>>21)&0x7FF) / 2047.0
	y := float64((packedVal>>11)&0x3FF) / 1023.0
	z := float64(packedVal&0x7FF) / 2047.0

	// 反归一化
	x = x*(chunk.MaxX-chunk.MinX) + chunk.MinX
	y = y*(chunk.MaxY-chunk.MinY) + chunk.MinY
	z = z*(chunk.MaxZ-chunk.MinZ) + chunk.MinZ

	return cmn.ClipFloat32(x), cmn.ClipFloat32(y), cmn.ClipFloat32(z)
}

func unpack111011Scale(packedVal uint32, chunk *Chunk) (float32, float32, float32) {
	// 提取 x, y, z 的值
	x := float64((packedVal>>21)&0x7FF) / 2047.0
	y := float64((packedVal>>11)&0x3FF) / 1023.0
	z := float64(packedVal&0x7FF) / 2047.0

	// 反归一化
	x = x*(chunk.MaxScaleX-chunk.MinScaleX) + chunk.MinScaleX
	y = y*(chunk.MaxScaleY-chunk.MinScaleY) + chunk.MinScaleY
	z = z*(chunk.MaxScaleZ-chunk.MinScaleZ) + chunk.MinScaleZ

	return cmn.ClipFloat32(x), cmn.ClipFloat32(y), cmn.ClipFloat32(z)
}

func unpackAndEncodeRgba(packedVal uint32, chunk *Chunk) (uint8, uint8, uint8, uint8) {
	// 提取 r, g, b, a 分量
	r := float64((packedVal>>24)&0xFF) / 255.0
	g := float64((packedVal>>16)&0xFF) / 255.0
	b := float64((packedVal>>8)&0xFF) / 255.0
	a := float64(packedVal&0xFF) / 255.0

	// 反归一化 r, g, b
	r = r*(chunk.MaxR-chunk.MinR) + chunk.MinR
	g = g*(chunk.MaxG-chunk.MinG) + chunk.MinG
	b = b*(chunk.MaxB-chunk.MinB) + chunk.MinB

	r = (r - 0.5) / SH_C0
	g = (g - 0.5) / SH_C0
	b = (b - 0.5) / SH_C0

	// 反 Sigmoid 函数处理 a
	a = -math.Log(1/a - 1)

	return cmn.EncodeSplatColor(r), cmn.EncodeSplatColor(g), cmn.EncodeSplatColor(b), cmn.EncodeSplatOpacity(a)
}

func unpackRotations(bs []byte) (uint8, uint8, uint8, uint8) {
	comp := cmn.BytesToUint32(bs)
	index := int(comp >> 30)
	remaining := comp
	sumSquares := 0.0
	rotation := []float64{0.0, 0.0, 0.0, 0.0}

	maxVal := uint32(0x3FF) // 1023
	for i := 3; i >= 0; i-- {
		if i != index {
			magnitude := float64(remaining & maxVal)
			remaining = remaining >> 10

			rotation[i] = ((magnitude / float64(maxVal)) - 0.5) / cmn.SQRT1_2
			sumSquares += rotation[i] * rotation[i]
		}
	}

	rotation[index] = math.Sqrt(math.Max(1.0-sumSquares, 0))

	r0, r1, r2, r3 := rotation[0], rotation[1], rotation[2], rotation[3]
	return cmn.ClipUint8(r0*128.0 + 128.0), cmn.ClipUint8(r1*128.0 + 128.0), cmn.ClipUint8(r2*128.0 + 128.0), cmn.ClipUint8(r3*128.0 + 128.0)
}
