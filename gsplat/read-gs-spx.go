package gsplat

import (
	"errors"
	"fmt"
	"gsbox/cmn"
	"math"
	"os"
)

func ReadSpx(spxFile string) (*SpxHeader, []*SplatData) {

	header := ParseSpxHeader(spxFile)
	if !header.IsValid() {
		fmt.Println("[WARN] hash check failed: " + spxFile)
	}

	file, err := os.Open(spxFile)
	cmn.ExitOnError(err)
	defer file.Close()

	datas := make([]*SplatData, 0)
	offset := int64(HeaderSizeSpx)
	var x0, x1, x2, y0, y1, y2, z0, z1, z2, i32x, i32y, i32z int32
	var n1, n2, n3 int
	for {
		// 块数据长度、是否压缩
		bts := make([]byte, 4)
		_, err = file.ReadAt(bts, offset)
		if err != nil {
			break
		}

		blockSize := int64(math.Abs(float64(cmn.BytesToInt32(bts[0:]))))
		isGzip := cmn.BytesToInt32(bts[0:]) < 0

		// 块数据读取
		offset += 4
		blockBytes := make([]byte, blockSize)
		_, err = file.ReadAt(blockBytes, offset)
		cmn.ExitOnError(err)
		offset += blockSize

		// 块数据解压
		var blockBts []byte
		if isGzip {
			blockBts, err = cmn.UnGzipBytes(blockBytes)
			cmn.ExitOnError(err)
		} else {
			blockBts = blockBytes
		}

		// 块数据格式
		i32BlockSplatCount := int32(cmn.BytesToUint32(blockBts[0:4]))
		blockSplatCount := int(i32BlockSplatCount) // 数量
		format := cmn.BytesToUint32(blockBts[4:8]) // 格式ID
		if format == 20 {
			// splat20
			dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
			for n := range blockSplatCount {
				splatData := &SplatData{}
				x0 = int32(dataBytes[n*3])
				x1 = int32(dataBytes[n*3+1])
				x2 = int32(dataBytes[n*3+2])
				y0 = int32(dataBytes[blockSplatCount*3+n*3])
				y1 = int32(dataBytes[blockSplatCount*3+n*3+1])
				y2 = int32(dataBytes[blockSplatCount*3+n*3+2])
				z0 = int32(dataBytes[blockSplatCount*6+n*3])
				z1 = int32(dataBytes[blockSplatCount*6+n*3+1])
				z2 = int32(dataBytes[blockSplatCount*6+n*3+2])
				splatData.ScaleX = cmn.DecodeByteToFloat32(dataBytes[blockSplatCount*9+n])
				splatData.ScaleY = cmn.DecodeByteToFloat32(dataBytes[blockSplatCount*10+n])
				splatData.ScaleZ = cmn.DecodeByteToFloat32(dataBytes[blockSplatCount*11+n])
				splatData.ColorR = dataBytes[blockSplatCount*12+n]
				splatData.ColorG = dataBytes[blockSplatCount*13+n]
				splatData.ColorB = dataBytes[blockSplatCount*14+n]
				splatData.ColorA = dataBytes[blockSplatCount*15+n]
				splatData.RotationX = dataBytes[blockSplatCount*16+n]
				splatData.RotationY = dataBytes[blockSplatCount*17+n]
				splatData.RotationZ = dataBytes[blockSplatCount*18+n]
				splatData.RotationW = dataBytes[blockSplatCount*19+n]

				i32x = x0 | (x1 << 8) | (x2 << 16)
				if i32x&0x800000 > 0 {
					i32x |= -0x1000000
				}
				i32y = y0 | (y1 << 8) | (y2 << 16)
				if i32y&0x800000 > 0 {
					i32y |= -0x1000000
				}
				i32z = z0 | (z1 << 8) | (z2 << 16)
				if i32z&0x800000 > 0 {
					i32z |= -0x1000000
				}
				splatData.PositionX = float32(i32x) / 4096.0
				splatData.PositionY = float32(i32y) / 4096.0
				splatData.PositionZ = float32(i32z) / 4096.0

				datas = append(datas, splatData)
			}

		} else if format == 1 {
			// SH1
			dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
			for n := range blockSplatCount {
				splatData := datas[n1+n]
				splatData.SH1 = dataBytes[n*9 : n*9+9]
			}
			n1 += blockSplatCount
		} else if format == 2 {
			// SH2
			dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
			for n := range blockSplatCount {
				splatData := datas[n2+n]
				splatData.SH2 = dataBytes[n*24 : n*24+24]
			}
			n2 += blockSplatCount
		} else if format == 3 {
			// SH3
			dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
			for n := range blockSplatCount {
				splatData := datas[n3+n]
				splatData.SH3 = dataBytes[n*21 : n*21+21]
			}
			n3 += blockSplatCount
		} else {
			// 存在无法识别读取的专有格式数据
			cmn.ExitOnError(errors.New("unreadable proprietary data exists"))
		}

	}

	return header, datas
}
