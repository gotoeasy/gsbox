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
		blockSplatCount := int(i32BlockSplatCount)   // 数量
		formatId := cmn.BytesToUint32(blockBts[4:8]) // 格式ID
		if formatId == 20 {
			// splat20
			dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
			for n := range blockSplatCount {
				splatData := &SplatData{}
				splatData.PositionX = cmn.DecodeSpxPositionUint24(dataBytes[n*3], dataBytes[n*3+1], dataBytes[n*3+2])
				splatData.PositionY = cmn.DecodeSpxPositionUint24(dataBytes[blockSplatCount*3+n*3], dataBytes[blockSplatCount*3+n*3+1], dataBytes[blockSplatCount*3+n*3+2])
				splatData.PositionZ = cmn.DecodeSpxPositionUint24(dataBytes[blockSplatCount*6+n*3], dataBytes[blockSplatCount*6+n*3+1], dataBytes[blockSplatCount*6+n*3+2])
				splatData.ScaleX = cmn.DecodeSpxScale(dataBytes[blockSplatCount*9+n])
				splatData.ScaleY = cmn.DecodeSpxScale(dataBytes[blockSplatCount*10+n])
				splatData.ScaleZ = cmn.DecodeSpxScale(dataBytes[blockSplatCount*11+n])
				splatData.ColorR = dataBytes[blockSplatCount*12+n]
				splatData.ColorG = dataBytes[blockSplatCount*13+n]
				splatData.ColorB = dataBytes[blockSplatCount*14+n]
				splatData.ColorA = dataBytes[blockSplatCount*15+n]
				splatData.RotationW = dataBytes[blockSplatCount*16+n]
				splatData.RotationX = dataBytes[blockSplatCount*17+n]
				splatData.RotationY = dataBytes[blockSplatCount*18+n]
				splatData.RotationZ = dataBytes[blockSplatCount*19+n]
				datas = append(datas, splatData)
			}

		} else if formatId == 1 {
			// SH1
			dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
			for n := range blockSplatCount {
				splatData := datas[n1+n]
				splatData.SH1 = dataBytes[n*9 : n*9+9]
			}
			n1 += blockSplatCount
		} else if formatId == 2 {
			// SH2
			dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
			for n := range blockSplatCount {
				splatData := datas[n2+n]
				splatData.SH2 = dataBytes[n*24 : n*24+24]
			}
			n2 += blockSplatCount
		} else if formatId == 3 {
			// SH3
			dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
			for n := range blockSplatCount {
				splatData := datas[n3+n]
				splatData.SH3 = dataBytes[n*21 : n*21+21]
			}
			n3 += blockSplatCount
		} else {
			// 存在无法识别读取的专有格式数据
			cmn.ExitOnError(errors.New("unknow block data format exists: " + cmn.Uint32ToString(formatId)))
		}

	}

	return header, datas
}
