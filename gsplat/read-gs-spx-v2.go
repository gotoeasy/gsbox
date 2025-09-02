package gsplat

import (
	"errors"
	"gsbox/cmn"
	"math"
	"os"
)

func ReadSpxV2(spxFile string, header *SpxHeader) (*SpxHeader, []*SplatData) {
	file, err := os.Open(spxFile)
	cmn.ExitOnError(err)
	defer file.Close()

	datas := make([]*SplatData, 0)
	offset := int64(HeaderSizeSpx)
	var n1, n2, n3 int
	for {
		// 块数据长度、是否压缩
		bts := make([]byte, 4)
		_, err = file.ReadAt(bts, offset)
		if err != nil {
			break
		}

		i32 := cmn.BytesToInt32(bts)
		isCompress := i32 < 0
		i64 := int64(math.Abs(float64(i32)))
		compressType := uint8((i64 >> 28) & 0b111)
		blockSize := i64 & ((1 << 29) - 1)

		// 块数据读取
		offset += 4
		blockBytes := make([]byte, blockSize)
		_, err = file.ReadAt(blockBytes, offset)
		cmn.ExitOnError(err)
		offset += blockSize

		// 块数据解压
		var blockBts []byte
		if isCompress {
			if compressType == 0 {
				blockBts, err = cmn.UnGzipBytes(blockBytes)
				cmn.ExitOnError(err)
			} else {
				cmn.ExitOnError(errors.New("unsupported compress type"))
			}
		} else {
			blockBts = blockBytes
		}

		// 块数据格式
		i32BlockSplatCount := int32(cmn.BytesToUint32(blockBts[0:4]))
		blkSplatCnt := int(i32BlockSplatCount)       // 数量
		formatId := cmn.BytesToUint32(blockBts[4:8]) // 格式ID
		switch formatId {
		case 16:
			// splat16
			minX := cmn.BytesToFloat32(blockBts[8:12])
			maxX := cmn.BytesToFloat32(blockBts[12:16])
			minY := cmn.BytesToFloat32(blockBts[16:20])
			maxY := cmn.BytesToFloat32(blockBts[20:24])
			minZ := cmn.BytesToFloat32(blockBts[24:28])
			maxZ := cmn.BytesToFloat32(blockBts[28:32])

			bts := blockBts[32:] // 除去前32字节（数量，格式，包围盒）
			for n := range blkSplatCnt {
				data := &SplatData{}
				x := cmn.BytesToUint16([]byte{bts[blkSplatCnt*0+n], bts[blkSplatCnt*3+n]})
				y := cmn.BytesToUint16([]byte{bts[blkSplatCnt*1+n], bts[blkSplatCnt*4+n]})
				z := cmn.BytesToUint16([]byte{bts[blkSplatCnt*2+n], bts[blkSplatCnt*5+n]})

				data.PositionX = cmn.DecodeSpxPositionUint16(x, minX, maxX)
				data.PositionY = cmn.DecodeSpxPositionUint16(y, minY, maxY)
				data.PositionZ = cmn.DecodeSpxPositionUint16(z, minZ, maxZ)
				data.ScaleX = cmn.DecodeSpxScale(bts[blkSplatCnt*6+n])
				data.ScaleY = cmn.DecodeSpxScale(bts[blkSplatCnt*7+n])
				data.ScaleZ = cmn.DecodeSpxScale(bts[blkSplatCnt*8+n])
				data.ColorR = bts[blkSplatCnt*9+n]
				data.ColorG = bts[blkSplatCnt*10+n]
				data.ColorB = bts[blkSplatCnt*11+n]
				data.ColorA = bts[blkSplatCnt*12+n]
				data.RotationW, data.RotationX, data.RotationY, data.RotationZ = cmn.DecodeSpxRotations(bts[blkSplatCnt*13+n], bts[blkSplatCnt*14+n], bts[blkSplatCnt*15+n])
				datas = append(datas, data)
			}
		case 19:
			// splat19
			bts := blockBts[8:] // 除去前8字节（数量，格式）
			for n := range blkSplatCnt {
				data := &SplatData{}
				data.PositionX = cmn.DecodeSpxPositionUint24(bts[blkSplatCnt*0+n], bts[blkSplatCnt*3+n], bts[blkSplatCnt*6+n])
				data.PositionY = cmn.DecodeSpxPositionUint24(bts[blkSplatCnt*1+n], bts[blkSplatCnt*4+n], bts[blkSplatCnt*7+n])
				data.PositionZ = cmn.DecodeSpxPositionUint24(bts[blkSplatCnt*2+n], bts[blkSplatCnt*5+n], bts[blkSplatCnt*8+n])
				data.ScaleX = cmn.DecodeSpxScale(bts[blkSplatCnt*9+n])
				data.ScaleY = cmn.DecodeSpxScale(bts[blkSplatCnt*10+n])
				data.ScaleZ = cmn.DecodeSpxScale(bts[blkSplatCnt*11+n])
				data.ColorR = bts[blkSplatCnt*12+n]
				data.ColorG = bts[blkSplatCnt*13+n]
				data.ColorB = bts[blkSplatCnt*14+n]
				data.ColorA = bts[blkSplatCnt*15+n]
				data.RotationW, data.RotationX, data.RotationY, data.RotationZ = cmn.DecodeSpxRotations(bts[blkSplatCnt*16+n], bts[blkSplatCnt*17+n], bts[blkSplatCnt*18+n])
				datas = append(datas, data)
			}
		case 20:
			// splat20
			bts := blockBts[8:] // 除去前8字节（数量，格式）
			for n := range blkSplatCnt {
				data := &SplatData{}
				data.PositionX = cmn.DecodeSpxPositionUint24(bts[n*3], bts[n*3+1], bts[n*3+2])
				data.PositionY = cmn.DecodeSpxPositionUint24(bts[blkSplatCnt*3+n*3], bts[blkSplatCnt*3+n*3+1], bts[blkSplatCnt*3+n*3+2])
				data.PositionZ = cmn.DecodeSpxPositionUint24(bts[blkSplatCnt*6+n*3], bts[blkSplatCnt*6+n*3+1], bts[blkSplatCnt*6+n*3+2])
				data.ScaleX = cmn.DecodeSpxScale(bts[blkSplatCnt*9+n])
				data.ScaleY = cmn.DecodeSpxScale(bts[blkSplatCnt*10+n])
				data.ScaleZ = cmn.DecodeSpxScale(bts[blkSplatCnt*11+n])
				data.ColorR = bts[blkSplatCnt*12+n]
				data.ColorG = bts[blkSplatCnt*13+n]
				data.ColorB = bts[blkSplatCnt*14+n]
				data.ColorA = bts[blkSplatCnt*15+n]
				data.RotationW, data.RotationX, data.RotationY, data.RotationZ = cmn.NormalizeRotations(bts[blkSplatCnt*16+n], bts[blkSplatCnt*17+n], bts[blkSplatCnt*18+n], bts[blkSplatCnt*19+n])
				datas = append(datas, data)
			}
		case 1:
			// SH1
			dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
			for n := range blkSplatCnt {
				splatData := datas[n1+n]
				splatData.SH1 = dataBytes[n*9 : n*9+9]
			}
			n1 += blkSplatCnt
		case 2:
			// SH2
			dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
			for n := range blkSplatCnt {
				splatData := datas[n2+n]
				splatData.SH2 = dataBytes[n*24 : n*24+24]
			}
			n2 += blkSplatCnt
		case 3:
			// SH3
			dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
			for n := range blkSplatCnt {
				splatData := datas[n3+n]
				splatData.SH3 = dataBytes[n*21 : n*21+21]
			}
			n3 += blkSplatCnt
		default:
			// 存在无法识别读取的专有格式数据
			cmn.ExitOnError(errors.New("unknow block data format exists: " + cmn.Uint32ToString(formatId)))
		}

	}

	return header, datas
}
