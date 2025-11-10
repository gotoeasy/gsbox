package gsplat

import (
	"errors"
	"gsbox/cmn"
	"math"
	"os"
)

func ReadSpxOpenV3(spxFile string, header *SpxHeader) (*SpxHeader, []*SplatData) {
	file, err := os.Open(spxFile)
	cmn.ExitOnError(err)
	defer file.Close()

	datas := make([]*SplatData, 0)
	offset := int64(HeaderSizeSpx)
	for {
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
		case BF_SPLAT22:
			readSpxBF22_V3(blockBts, blkSplatCnt, header, &datas)
		case BF_SPLAT220_WEBP:
			readSpxBF220_WEBP_V3(blockBts, blkSplatCnt, header, &datas)
		case BF_SH_PALETTES:
			readSpxPalettes_V3(header, blockBts)
		case BF_SH_PALETTES_WEBP:
			readSpxPalettesWebp_V3(header, blockBts)
		default:
			// 存在无法识别读取的专有格式数据
			cmn.ExitOnError(errors.New("unknow block data format exists: " + cmn.Uint32ToString(formatId)))
		}

	}

	// 按调色板设定球谐系数
	setAllShByPalettes(header, datas)

	return header, datas
}

func readSpxBF22_V3(blockBts []byte, blkSplatCnt int, header *SpxHeader, datas *[]*SplatData) {
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
		data.RotationW = bts[blkSplatCnt*16+n]
		data.RotationX = bts[blkSplatCnt*17+n]
		data.RotationY = bts[blkSplatCnt*18+n]
		data.RotationZ = bts[blkSplatCnt*19+n]
		if header.ShDegree > 0 {
			data.ShPaletteIdx = uint16(bts[blkSplatCnt*20+n]) | (uint16(bts[blkSplatCnt*21+n]) << 8)
		}

		*datas = append(*datas, data)
	}

}

func readSpxBF220_WEBP_V3(blockBts []byte, blkSplatCnt int, header *SpxHeader, datas *[]*SplatData) {
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
		ri := btsRotations[n*4+3]

		data := &SplatData{}
		data.PositionX = cmn.DecodeLog(cmn.DecodeSpxPositionUint24(x0, x1, x2), 1)
		data.PositionY = cmn.DecodeLog(cmn.DecodeSpxPositionUint24(y0, y1, y2), 1)
		data.PositionZ = cmn.DecodeLog(cmn.DecodeSpxPositionUint24(z0, z1, z2), 1)
		data.ScaleX = cmn.DecodeSpxScale(btsScales[n*4+0])
		data.ScaleY = cmn.DecodeSpxScale(btsScales[n*4+1])
		data.ScaleZ = cmn.DecodeSpxScale(btsScales[n*4+2])
		data.ColorR = btsColors[n*4+0]
		data.ColorG = btsColors[n*4+1]
		data.ColorB = btsColors[n*4+2]
		data.ColorA = btsColors[n*4+3]
		data.RotationW, data.RotationX, data.RotationY, data.RotationZ = cmn.SogDecodeRotations(rx, ry, rz, ri)

		if len(btsPaletteIdxs) > 0 {
			b0 := uint16(btsPaletteIdxs[n*4+0])
			b1 := uint16(btsPaletteIdxs[n*4+1])
			data.ShPaletteIdx = b0 & (b1 << 8)
		}

		*datas = append(*datas, data)
	}
}

func readSpxPalettes_V3(header *SpxHeader, blockBts []byte) {
	// 调色板
	centroids := blockBts[8:] // 除去前8字节（数量，格式）
	header.ShPalettes = centroids
}

func readSpxPalettesWebp_V3(header *SpxHeader, blockBts []byte) {
	// 调色板
	dataBytes := blockBts[8:] // 除去前8字节（数量，格式）
	centroids, _, _, err := cmn.DecompressWebp(dataBytes)
	cmn.ExitOnError(err)
	header.ShPalettes = centroids
}

func setAllShByPalettes(header *SpxHeader, rows []*SplatData) {
	if len(header.ShPalettes) > 0 {
		for _, d := range rows {
			setShByPalettes(d, header.ShPalettes, header.ShDegree)
		}
	}
}

func setShByPalettes(d *SplatData, centroids []uint8, shDegree uint8) {
	col := (d.ShPaletteIdx & 63) * 15 // 同 (n % 64) * 15
	row := d.ShPaletteIdx >> 6        // 同 Math.floor(n / 64)
	offset := int(row*960 + col)      // 960 = 64 * 15

	sh1 := make([]uint8, 9)
	sh2 := make([]uint8, 15)
	sh3 := make([]uint8, 21)
	for d := range 3 {
		for k := range 3 {
			sh1[k*3+d] = centroids[(offset+k)*4+d]
		}
		for k := range 5 {
			sh2[k*3+d] = centroids[(offset+3+k)*4+d]
		}
		for k := range 7 {
			sh3[k*3+d] = centroids[(offset+8+k)*4+d]
		}
	}
	var shs []uint8
	shs = append(shs, sh1...)
	shs = append(shs, sh2...)
	shs = append(shs, sh3...)

	switch shDegree {
	case 3:
		d.SH2 = shs[:24]
		d.SH3 = shs[24:]
	case 2:
		d.SH2 = shs[:24]
	case 1:
		d.SH1 = shs[:9]
	}
}
