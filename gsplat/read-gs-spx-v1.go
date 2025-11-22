package gsplat

import (
	"errors"
	"gsbox/cmn"
	"math"
	"os"
)

func ReadSpxOpenV1(spxFile string, header *SpxHeader) (*SpxHeader, []*SplatData) {
	file, err := os.Open(spxFile)
	cmn.ExitOnError(err)
	defer file.Close()

	datas := make([]*SplatData, 0)
	offset := int64(HeaderSizeSpx)
	var n1, n2, n3 int
	for {
		OnProgress(PhaseRead, len(datas), int(header.SplatCount))
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
			blockBts, err = cmn.DecompressGzip(blockBytes)
			cmn.ExitOnError(err)
		} else {
			blockBts = blockBytes
		}

		// 块数据格式
		i32BlockSplatCount := int32(cmn.BytesToUint32(blockBts[0:4]))
		blkSplatCnt := int(i32BlockSplatCount)       // 数量
		formatId := cmn.BytesToUint32(blockBts[4:8]) // 格式ID
		switch formatId {
		case BF_SPLAT20:
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
		case BF_SH1:
			// SH1
			dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
			for n := range blkSplatCnt {
				splatData := datas[n1+n]
				splatData.SH1 = dataBytes[n*9 : n*9+9]
			}
			n1 += blkSplatCnt
		case BF_SH2:
			// SH2
			dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
			for n := range blkSplatCnt {
				splatData := datas[n2+n]
				splatData.SH2 = dataBytes[n*24 : n*24+24]
			}
			n2 += blkSplatCnt
		case BF_SH3:
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
