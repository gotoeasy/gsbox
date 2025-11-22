package gsplat

import (
	"errors"
	"gsbox/cmn"
	"math"
	"os"
)

func ReadSpxOpenV2(spxFile string, header *SpxHeader) (*SpxHeader, []*SplatData) {
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

		i32 := cmn.BytesToInt32(bts)
		isCompress := i32 < 0
		u32 := uint32(math.Abs(float64(i32)))
		compressType := uint8((u32 << 1) >> 29)
		blockSize := int64((u32 << 4) >> 4)

		// 块数据读取
		offset += 4
		blockBytes := make([]byte, blockSize)
		_, err = file.ReadAt(blockBytes, offset)
		cmn.ExitOnError(err)
		offset += blockSize

		// 块数据解压
		var blockBts []byte
		if isCompress {
			switch compressType {
			case CT_GZIP:
				blockBts, err = cmn.DecompressGzip(blockBytes)
				cmn.ExitOnError(err)
			case CT_XZ:
				blockBts, err = cmn.DecompressXZ(blockBytes)
				cmn.ExitOnError(err)
			default:
				cmn.ExitOnError(errors.New("unsupported compression type"))
			}
		} else {
			blockBts = blockBytes
		}

		// 块数据格式
		i32BlockSplatCount := int32(cmn.BytesToUint32(blockBts[0:4]))
		blkSplatCnt := int(i32BlockSplatCount)       // 数量
		formatId := cmn.BytesToUint32(blockBts[4:8]) // 格式ID
		switch formatId {
		case BF_SPLAT19:
			readSpxBF_SPLAT19(blockBts, blkSplatCnt, &datas)
		case BF_SPLAT10019:
			readSpxBF_SPLAT10019(blockBts, blkSplatCnt, &datas)
		case BF_SPLAT190_WEBP:
			readSpxBF_SPLAT190_WEBP(blockBts, blkSplatCnt, &datas)
		case BF_SPLAT10190_WEBP:
			readSpxBF_SPLAT10190_WEBP(blockBts, blkSplatCnt, &datas)
		case BF_SPLAT20:
			readSpxBF_SPLAT20(blockBts, blkSplatCnt, &datas)
		case BF_SH1:
			readSpxBF_SH1(blockBts, blkSplatCnt, &datas, n1)
			n1 += blkSplatCnt
		case BF_SH2:
			readSpxBF_SH2(blockBts, blkSplatCnt, &datas, n2)
			n2 += blkSplatCnt
		case BF_SH3:
			readSpxBF_SH3(blockBts, blkSplatCnt, &datas, n3)
			n3 += blkSplatCnt
		case BF_SH3_WEBP:
			readSpxBF_SH3_WEBP(blockBts, blkSplatCnt, &datas, n3, header)
			n1 += blkSplatCnt
			n2 += blkSplatCnt
			n3 += blkSplatCnt
		default:
			// 存在无法识别读取的专有格式数据
			cmn.ExitOnError(errors.New("unknow block data format exists: " + cmn.Uint32ToString(formatId)))
		}

	}

	return header, datas
}

func readSpxBF_SPLAT19(blockBts []byte, blkSplatCnt int, datas *[]*SplatData) {
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
		*datas = append(*datas, data)
	}

}

func readSpxBF_SPLAT10019(blockBts []byte, blkSplatCnt int, datas *[]*SplatData) {
	// splat10019
	logTimes := int(blockBts[8]) // log编码次数
	bts := blockBts[12:]         // 除去前12字节（数量，格式，log编码次数）
	for n := range blkSplatCnt {
		data := &SplatData{}
		data.PositionX = cmn.DecodeLog(cmn.DecodeSpxPositionUint24(bts[blkSplatCnt*0+n], bts[blkSplatCnt*3+n], bts[blkSplatCnt*6+n]), logTimes)
		data.PositionY = cmn.DecodeLog(cmn.DecodeSpxPositionUint24(bts[blkSplatCnt*1+n], bts[blkSplatCnt*4+n], bts[blkSplatCnt*7+n]), logTimes)
		data.PositionZ = cmn.DecodeLog(cmn.DecodeSpxPositionUint24(bts[blkSplatCnt*2+n], bts[blkSplatCnt*5+n], bts[blkSplatCnt*8+n]), logTimes)
		data.ScaleX = cmn.DecodeSpxScale(bts[blkSplatCnt*9+n])
		data.ScaleY = cmn.DecodeSpxScale(bts[blkSplatCnt*10+n])
		data.ScaleZ = cmn.DecodeSpxScale(bts[blkSplatCnt*11+n])
		data.ColorR = bts[blkSplatCnt*12+n]
		data.ColorG = bts[blkSplatCnt*13+n]
		data.ColorB = bts[blkSplatCnt*14+n]
		data.ColorA = bts[blkSplatCnt*15+n]
		data.RotationW, data.RotationX, data.RotationY, data.RotationZ = cmn.DecodeSpxRotations(bts[blkSplatCnt*16+n], bts[blkSplatCnt*17+n], bts[blkSplatCnt*18+n])
		*datas = append(*datas, data)
	}
}

func readSpxBF_SPLAT190_WEBP(blockBts []byte, blkSplatCnt int, datas *[]*SplatData) {
	// WEBP190
	bts := blockBts[8:] // 除去前8字节（数量，格式）
	size := cmn.BytesToUint32(bts[:4])
	webps := bts[4 : size+4]
	btsPositions, _, _, err := cmn.DecompressWebp(webps)
	cmn.ExitOnError(err)

	bts = bts[size+4:]
	size = cmn.BytesToUint32(bts[:4])
	webps = bts[4 : size+4]
	btsScales, _, _, err := cmn.DecompressWebp(webps)
	cmn.ExitOnError(err)

	bts = bts[size+4:]
	size = cmn.BytesToUint32(bts[:4])
	webps = bts[4 : size+4]
	btsColors, _, _, err := cmn.DecompressWebp(webps)
	cmn.ExitOnError(err)

	bts = bts[size+4:]
	size = cmn.BytesToUint32(bts[:4])
	webps = bts[4 : size+4]
	btsRotations, _, _, err := cmn.DecompressWebp(webps)
	cmn.ExitOnError(err)

	for n := range blkSplatCnt {
		x0 := btsPositions[n*4+0]
		y0 := btsPositions[n*4+1]
		z0 := btsPositions[n*4+2]
		x1 := btsPositions[blkSplatCnt*4+n*4+0]
		y1 := btsPositions[blkSplatCnt*4+n*4+1]
		z1 := btsPositions[blkSplatCnt*4+n*4+2]
		x2 := btsPositions[blkSplatCnt*8+n*4+0]
		y2 := btsPositions[blkSplatCnt*8+n*4+1]
		z2 := btsPositions[blkSplatCnt*8+n*4+2]
		rx := btsRotations[n*4+0]
		ry := btsRotations[n*4+1]
		rz := btsRotations[n*4+2]

		data := &SplatData{}
		data.PositionX = cmn.DecodeSpxPositionUint24(x0, x1, x2)
		data.PositionY = cmn.DecodeSpxPositionUint24(y0, y1, y2)
		data.PositionZ = cmn.DecodeSpxPositionUint24(z0, z1, z2)
		data.ScaleX = cmn.DecodeSpxScale(btsScales[n*4+0])
		data.ScaleY = cmn.DecodeSpxScale(btsScales[n*4+1])
		data.ScaleZ = cmn.DecodeSpxScale(btsScales[n*4+2])
		data.ColorR = btsColors[n*4+0]
		data.ColorG = btsColors[n*4+1]
		data.ColorB = btsColors[n*4+2]
		data.ColorA = btsColors[n*4+3]
		data.RotationW, data.RotationX, data.RotationY, data.RotationZ = cmn.DecodeSpxRotations(rx, ry, rz)

		*datas = append(*datas, data)
	}
}

func readSpxBF_SPLAT10190_WEBP(blockBts []byte, blkSplatCnt int, datas *[]*SplatData) {
	// WEBP10190
	logTimes := int(blockBts[8]) // log编码次数
	bts := blockBts[12:]         // 除去前12字节（数量，格式，log编码次数）
	size := cmn.BytesToUint32(bts[:4])
	webps := bts[4 : size+4]
	btsPositions, _, _, err := cmn.DecompressWebp(webps)
	cmn.ExitOnError(err)

	bts = bts[size+4:]
	size = cmn.BytesToUint32(bts[:4])
	webps = bts[4 : size+4]
	btsScales, _, _, err := cmn.DecompressWebp(webps)
	cmn.ExitOnError(err)

	bts = bts[size+4:]
	size = cmn.BytesToUint32(bts[:4])
	webps = bts[4 : size+4]
	btsColors, _, _, err := cmn.DecompressWebp(webps)
	cmn.ExitOnError(err)

	bts = bts[size+4:]
	size = cmn.BytesToUint32(bts[:4])
	webps = bts[4 : size+4]
	btsRotations, _, _, err := cmn.DecompressWebp(webps)
	cmn.ExitOnError(err)

	for n := range blkSplatCnt {
		x0 := btsPositions[n*4+0]
		y0 := btsPositions[n*4+1]
		z0 := btsPositions[n*4+2]
		x1 := btsPositions[blkSplatCnt*4+n*4+0]
		y1 := btsPositions[blkSplatCnt*4+n*4+1]
		z1 := btsPositions[blkSplatCnt*4+n*4+2]
		x2 := btsPositions[blkSplatCnt*8+n*4+0]
		y2 := btsPositions[blkSplatCnt*8+n*4+1]
		z2 := btsPositions[blkSplatCnt*8+n*4+2]
		rx := btsRotations[n*4+0]
		ry := btsRotations[n*4+1]
		rz := btsRotations[n*4+2]

		data := &SplatData{}
		data.PositionX = cmn.DecodeLog(cmn.DecodeSpxPositionUint24(x0, x1, x2), logTimes)
		data.PositionY = cmn.DecodeLog(cmn.DecodeSpxPositionUint24(y0, y1, y2), logTimes)
		data.PositionZ = cmn.DecodeLog(cmn.DecodeSpxPositionUint24(z0, z1, z2), logTimes)
		data.ScaleX = cmn.DecodeSpxScale(btsScales[n*4+0])
		data.ScaleY = cmn.DecodeSpxScale(btsScales[n*4+1])
		data.ScaleZ = cmn.DecodeSpxScale(btsScales[n*4+2])
		data.ColorR = btsColors[n*4+0]
		data.ColorG = btsColors[n*4+1]
		data.ColorB = btsColors[n*4+2]
		data.ColorA = btsColors[n*4+3]
		data.RotationW, data.RotationX, data.RotationY, data.RotationZ = cmn.DecodeSpxRotations(rx, ry, rz)

		*datas = append(*datas, data)
	}
}

func readSpxBF_SPLAT20(blockBts []byte, blkSplatCnt int, datas *[]*SplatData) {
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
		*datas = append(*datas, data)
	}
}

func readSpxBF_SH1(blockBts []byte, blkSplatCnt int, datas *[]*SplatData, n1 int) {
	// SH1
	dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
	for n := range blkSplatCnt {
		splatData := (*datas)[n1+n]
		splatData.SH1 = dataBytes[n*9 : n*9+9]
	}
}

func readSpxBF_SH2(blockBts []byte, blkSplatCnt int, datas *[]*SplatData, n2 int) {
	// SH2
	dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
	for n := range blkSplatCnt {
		splatData := (*datas)[n2+n]
		splatData.SH2 = dataBytes[n*24 : n*24+24]
	}
}

func readSpxBF_SH3(blockBts []byte, blkSplatCnt int, datas *[]*SplatData, n3 int) {
	// SH3
	dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
	for n := range blkSplatCnt {
		splatData := (*datas)[n3+n]
		splatData.SH3 = dataBytes[n*21 : n*21+21]
	}
}

func readSpxBF_SH3_WEBP(blockBts []byte, blkSplatCnt int, datas *[]*SplatData, n3 int, header *SpxHeader) {
	// SH1~SH3
	dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
	rgba, _, _, err := cmn.DecompressWebp(dataBytes)
	cmn.ExitOnError(err)
	dataBytes = RgbaToSh(rgba, blkSplatCnt, 15)

	for n := range blkSplatCnt {
		splatData := (*datas)[n3+n]
		switch header.ShDegree {
		case 1:
			splatData.SH1 = dataBytes[n*45 : n*45+9]
		case 2:
			splatData.SH2 = dataBytes[n*45 : n*45+24]
		case 3:
			splatData.SH2 = dataBytes[n*45 : n*45+24]
			splatData.SH3 = dataBytes[n*45+24 : n*45+45]
		}
	}
}
