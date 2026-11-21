package gsplat

import (
	"errors"
	"gsbox/cmn"
	"math"
	"os"
)

func ReadSpxOpenV4(spxFile string, header *SpxHeader) (*SpxHeader, []*SplatData) {
	file, err := os.Open(spxFile)
	cmn.ExitOnError(err)
	defer file.Close()

	datas := make([]*SplatData, 0)
	offset := int64(HeaderSizeSpx)
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
		case BF_SPLAT21:
			readSpxV4(blockBts, blkSplatCnt, header, &datas)
		case BF_SPLAT210_WEBP:
			readSpxWebpV4(blockBts, blkSplatCnt, header, &datas)
		case BF_SH_PALETTES:
			readSpxPalettes(header, blockBts)
		case BF_SH_PALETTES_WEBP:
			readSpxPalettesWebp(header, blockBts)
		default:
			// 存在无法识别读取的专有格式数据
			cmn.ExitOnError(errors.New("unknow block data format exists: " + cmn.Uint32ToString(formatId)))
		}
	}

	// 按调色板设定球谐系数
	if header.ShDegree > 0 {
		setAllShByPalettes(header, datas)
	}

	return header, datas
}

func readSpxV4(blockBts []byte, blkSplatCnt int, header *SpxHeader, datas *[]*SplatData) {
	bts := blockBts[8:] // 除去前8字节（数量，格式）
	for n := range blkSplatCnt {
		data := &SplatData{}
		data.PositionX = cmn.DecodeSpxPositionUint24(bts[blkSplatCnt*0+n], bts[blkSplatCnt*3+n], bts[blkSplatCnt*6+n])
		data.PositionY = cmn.DecodeSpxPositionUint24(bts[blkSplatCnt*1+n], bts[blkSplatCnt*4+n], bts[blkSplatCnt*7+n])
		data.PositionZ = cmn.DecodeSpxPositionUint24(bts[blkSplatCnt*2+n], bts[blkSplatCnt*5+n], bts[blkSplatCnt*8+n])
		data.ScaleX, data.ScaleY, data.ScaleZ = cmn.DecodeSpxScale(bts[blkSplatCnt*9+n]), cmn.DecodeSpxScale(bts[blkSplatCnt*10+n]), cmn.DecodeSpxScale(bts[blkSplatCnt*11+n])
		data.ColorR, data.ColorG, data.ColorB, data.ColorA = bts[blkSplatCnt*12+n], bts[blkSplatCnt*13+n], bts[blkSplatCnt*14+n], bts[blkSplatCnt*15+n]
		data.RotationW, data.RotationX, data.RotationY, data.RotationZ = cmn.DecodeRotationsExpMap(bts[blkSplatCnt*16+n], bts[blkSplatCnt*17+n], bts[blkSplatCnt*18+n])
		if header.ShDegree > 0 {
			data.PaletteIdx = uint16(bts[blkSplatCnt*19+n]) | (uint16(bts[blkSplatCnt*20+n]) << 8)
		}

		*datas = append(*datas, data)
	}

}

func readSpxWebpV4(blockBts []byte, blkSplatCnt int, header *SpxHeader, datas *[]*SplatData) {
	bts := blockBts[8:] // 除去前8字节（数量，格式）

	size := cmn.BytesToUint32(bts[:4])
	webps := bts[4 : size+4]
	btsXyz0, _, _, err := cmn.DecompressWebp(webps)
	cmn.ExitOnError(err)

	bts = bts[size+4:]
	size = cmn.BytesToUint32(bts[:4])
	webps = bts[4 : size+4]
	btsXyz1, _, _, err := cmn.DecompressWebp(webps)
	cmn.ExitOnError(err)

	bts = bts[size+4:]
	size = cmn.BytesToUint32(bts[:4])
	webps = bts[4 : size+4]
	btsXyz2, _, _, err := cmn.DecompressWebp(webps)
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

	var btsPaletteIdxs []byte
	if header.ShDegree > 0 {
		bts = bts[size+4:]
		size = cmn.BytesToUint32(bts[:4])
		webps = bts[4 : size+4]
		paletteIdxBytes, _, _, err := cmn.DecompressWebp(webps)
		cmn.ExitOnError(err)
		btsPaletteIdxs = paletteIdxBytes
	}

	for n := range blkSplatCnt {
		x0 := btsXyz0[n*4+0]
		y0 := btsXyz0[n*4+1]
		z0 := btsXyz0[n*4+2]
		x1 := btsXyz1[n*4+0]
		y1 := btsXyz1[n*4+1]
		z1 := btsXyz1[n*4+2]
		x2 := btsXyz2[n*4+0]
		y2 := btsXyz2[n*4+1]
		z2 := btsXyz2[n*4+2]

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
		data.RotationW, data.RotationX, data.RotationY, data.RotationZ = cmn.DecodeRotationsExpMap(btsRotations[n*4+0], btsRotations[n*4+1], btsRotations[n*4+2])

		if len(btsPaletteIdxs) > 0 {
			data.PaletteIdx = uint16(btsPaletteIdxs[n*4+0]) | (uint16(btsPaletteIdxs[n*4+1]) << 8)
		}

		*datas = append(*datas, data)
	}
}

func readSpxPalettes(header *SpxHeader, blockBts []byte) {
	// 调色板
	header.Palettes = blockBts[8:] // 除去前8字节（数量，格式）
}

func readSpxPalettesWebp(header *SpxHeader, blockBts []byte) {
	// 调色板webp
	dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
	centroids, _, _, err := cmn.DecompressWebp(dataBytes)
	cmn.ExitOnError(err)
	header.Palettes = centroids
}
