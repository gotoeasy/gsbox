package gsplat

import (
	"bufio"
	"gsbox/cmn"
	"log"
	"math"
	"os"
	"sort"
)

// Deprecated
func WriteSpxOpenV1(spxFile string, rows []*SplatData, comment string, shDegree uint8) {
	file, err := os.Create(spxFile)
	cmn.ExitOnError(err)
	defer file.Close()

	writer := bufio.NewWriter(file)

	inputBlockSize := Args.GetArgIgnorecase("-bs", "--block-size")
	blockSize := cmn.StringToInt(inputBlockSize, DefaultBlockSize)
	if cmn.EqualsIngoreCase(inputBlockSize, "max") {
		blockSize = MaxBlockSize // 支持 -bs max 写法
	} else {
		blockSize = max(MinBlockSize, min(blockSize, MaxBlockSize)) // 超出范围时限定为边界值
	}

	header := genSpxHeader(rows, comment, shDegree, 0, 0, 0)
	_, err = writer.Write(header.ToBytes())
	cmn.ExitOnError(err)

	log.Println("[Info] (parameter) bf:", BF_SPLAT20, BlockFormatDesc(BF_SPLAT20))
	log.Println("[Info] (parameter) bs:", blockSize, "(block size)")

	var blockDatasList [][]*SplatData
	blockCnt := (int(header.SplatCount) + blockSize - 1) / blockSize
	for i := range blockCnt {
		blockDatas := make([]*SplatData, 0)
		max := min(i*blockSize+blockSize, int(header.SplatCount))
		for n := i * blockSize; n < max; n++ {
			blockDatas = append(blockDatas, rows[n])
		}
		writeSpxBlockSplat20(writer, blockDatas, len(blockDatas), 0)
		blockDatasList = append(blockDatasList, blockDatas)
	}

	switch shDegree {
	case 1:
		for i := range blockDatasList {
			writeSpxBlockSH1(writer, blockDatasList[i], 0)
		}
	case 2:
		for i := range blockDatasList {
			writeSpxBlockSH2(writer, blockDatasList[i], 0)
		}
	case 3:
		for i := range blockDatasList {
			writeSpxBlockSH2(writer, blockDatasList[i], 0)
			writeSpxBlockSH3(writer, blockDatasList[i], 0)
		}
	}

	err = writer.Flush()
	cmn.ExitOnError(err)
}

func genSpxHeader(datas []*SplatData, comment string, shDegree uint8, flag1 uint8, flag2 uint8, flag3 uint8) *SpxHeader {

	header := new(SpxHeader)
	header.Fixed = "spx"
	header.Version = 1
	header.SplatCount = int32(len(datas))

	header.CreateDate = cmn.GetSystemDateYYYYMMDD() // 创建日期
	header.CreaterId = GetOutputCreaterId()         // 0:官方默认识别号，（这里参考阿佩里常数1.202056903159594…以示区分，此常数由瑞士数学家罗杰·阿佩里在1978年证明其无理数性质而闻名）
	header.ExclusiveId = GetOutputExclusiveId()     // 0:官方开放格式的识别号
	header.ShDegree = uint8(shDegree)
	header.Flag1 = flag1
	header.Flag2 = flag2
	header.Flag3 = flag3
	header.Reserve1 = 0
	header.Reserve2 = 0
	del, comment := cmn.RemoveNonASCII(comment)
	if del {
		log.Println("[Warn] The existing non-ASCII characters in the comment have been removed!")
	}
	header.Comment = comment // 注释
	if header.Comment == "" {
		header.Comment = DefaultSpxComment()
	}

	if len(datas) > 0 {
		minX := float64(datas[0].PositionX)
		minY := float64(datas[0].PositionY)
		minZ := float64(datas[0].PositionZ)
		maxX := float64(datas[0].PositionX)
		maxY := float64(datas[0].PositionY)
		maxZ := float64(datas[0].PositionZ)

		for i := 1; i < len(datas); i++ {
			minX = math.Min(minX, float64(datas[i].PositionX))
			minY = math.Min(minY, float64(datas[i].PositionY))
			minZ = math.Min(minZ, float64(datas[i].PositionZ))
			maxX = math.Max(maxX, float64(datas[i].PositionX))
			maxY = math.Max(maxY, float64(datas[i].PositionY))
			maxZ = math.Max(maxZ, float64(datas[i].PositionZ))
		}
		header.MinX = cmn.ToFloat32(minX)
		header.MaxX = cmn.ToFloat32(maxX)
		header.MinY = cmn.ToFloat32(minY)
		header.MaxY = cmn.ToFloat32(maxY)
		header.MinZ = cmn.ToFloat32(minZ)
		header.MaxZ = cmn.ToFloat32(maxZ)

		// TopY
		centerX := cmn.ToFloat32((maxX + minX) / 2)
		centerY := cmn.ToFloat32((maxY + minY) / 2)
		centerZ := cmn.ToFloat32((maxZ + minZ) / 2)
		radius10 := math.Sqrt(float64((centerX-header.MaxX)*(centerX-header.MaxX)+
			(centerY-header.MaxY)*(centerY-header.MaxY)+
			(centerZ-header.MaxZ)*(centerZ-header.MaxZ))) * 0.1

		minTopY := float64(header.MaxY)
		maxTopY := float64(header.MinY)
		for i := range datas {
			if math.Abs(float64(datas[i].PositionY)) < 30 && math.Sqrt(float64(datas[i].PositionX)*float64(datas[i].PositionX)+float64(datas[i].PositionZ)*float64(datas[i].PositionZ)) <= radius10 {
				minTopY = math.Min(minTopY, float64(datas[i].PositionY))
				maxTopY = math.Max(maxTopY, float64(datas[i].PositionY))
			}
		}
		header.MinTopY = cmn.ToFloat32(minTopY)
		header.MaxTopY = cmn.ToFloat32(maxTopY)
	}

	return header
}

func writeSpxBlockSplat20(writer *bufio.Writer, blockDatas []*SplatData, blockSplatCount int, compressType uint8) {
	sort.Slice(blockDatas, func(i, j int) bool {
		return blockDatas[i].PositionY < blockDatas[j].PositionY // 坐标分别占3字节，按其中任一排序以更利于压缩
	})

	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(blockSplatCount))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(20)...)                      // 开放的块数据格式 20

	for n := range blockSplatCount {
		bts = append(bts, cmn.EncodeSpxPositionUint24(blockDatas[n].PositionX)...)
	}
	for n := range blockSplatCount {
		bts = append(bts, cmn.EncodeSpxPositionUint24(blockDatas[n].PositionY)...)
	}
	for n := range blockSplatCount {
		bts = append(bts, cmn.EncodeSpxPositionUint24(blockDatas[n].PositionZ)...)
	}
	for n := range blockSplatCount {
		bts = append(bts, cmn.EncodeSpxScale(blockDatas[n].ScaleX))
	}
	for n := range blockSplatCount {
		bts = append(bts, cmn.EncodeSpxScale(blockDatas[n].ScaleY))
	}
	for n := range blockSplatCount {
		bts = append(bts, cmn.EncodeSpxScale(blockDatas[n].ScaleZ))
	}
	for n := range blockSplatCount {
		bts = append(bts, blockDatas[n].ColorR)
	}
	for n := range blockSplatCount {
		bts = append(bts, blockDatas[n].ColorG)
	}
	for n := range blockSplatCount {
		bts = append(bts, blockDatas[n].ColorB)
	}
	for n := range blockSplatCount {
		bts = append(bts, blockDatas[n].ColorA)
	}
	for n := range blockSplatCount {
		bts = append(bts, blockDatas[n].RotationW)
	}
	for n := range blockSplatCount {
		bts = append(bts, blockDatas[n].RotationX)
	}
	for n := range blockSplatCount {
		bts = append(bts, blockDatas[n].RotationY)
	}
	for n := range blockSplatCount {
		bts = append(bts, blockDatas[n].RotationZ)
	}

	if blockSplatCount >= MinCompressBlockSize {
		var err error
		if compressType == CT_XZ {
			bts, err = cmn.CompressXZ(bts)
		} else {
			bts, err = cmn.CompressGzip(bts)
		}
		cmn.ExitOnError(err)
		blockByteLength := -((int32(compressType) << 28) | int32(len(bts)))
		_, err = writer.Write(cmn.Int32ToBytes(blockByteLength))
		cmn.ExitOnError(err)
		_, err = writer.Write(bts)
		cmn.ExitOnError(err)
	} else {
		blockByteLength := int32(len(bts))
		_, err := writer.Write(cmn.Int32ToBytes(blockByteLength))
		cmn.ExitOnError(err)
		_, err = writer.Write(bts)
		cmn.ExitOnError(err)
	}
}

func writeSpxBlockSH1(writer *bufio.Writer, blockDatas []*SplatData, compressType uint8) {
	blockSplatCount := len(blockDatas)
	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(blockSplatCount))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(BF_SH1)...)                  // 开放的块数据格式 1:sh1

	splatCnt := 0
	for n := range blockSplatCount {
		if blockDatas[n].IsWaterMark {
			continue
		}
		splatCnt++
		if len(blockDatas[n].SH1) > 0 {
			for i := range 9 {
				bts = append(bts, cmn.EncodeSpxSH(blockDatas[n].SH1[i]))
			}
		} else if len(blockDatas[n].SH2) > 0 {
			for i := range 9 {
				bts = append(bts, cmn.EncodeSpxSH(blockDatas[n].SH2[i]))
			}
		} else {
			for range 9 {
				bts = append(bts, cmn.EncodeSplatSH(0.0))
			}
		}
	}
	if splatCnt == 0 {
		return
	} else if splatCnt < blockSplatCount {
		cntBytes := cmn.Uint32ToBytes(uint32(splatCnt))
		bts[0] = cntBytes[0]
		bts[1] = cntBytes[1]
		bts[2] = cntBytes[2]
		bts[3] = cntBytes[3]
	}

	if splatCnt >= MinCompressBlockSize {
		var err error
		switch compressType {
		case CT_XZ:
			bts, err = cmn.CompressXZ(bts)
		default:
			bts, err = cmn.CompressGzip(bts)
		}
		cmn.ExitOnError(err)
		blockByteLength := -((int32(compressType) << 28) | int32(len(bts)))
		_, err = writer.Write(cmn.Int32ToBytes(blockByteLength))
		cmn.ExitOnError(err)
		_, err = writer.Write(bts)
		cmn.ExitOnError(err)
	} else {
		blockByteLength := int32(len(bts))
		_, err := writer.Write(cmn.Int32ToBytes(blockByteLength))
		cmn.ExitOnError(err)
		_, err = writer.Write(bts)
		cmn.ExitOnError(err)
	}
}

func writeSpxBlockSH2(writer *bufio.Writer, blockDatas []*SplatData, compressType uint8) {
	blockSplatCount := len(blockDatas)
	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(blockSplatCount))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(BF_SH2)...)                  // 开放的块数据格式 2:sh2

	splatCnt := 0
	for n := range blockSplatCount {
		if blockDatas[n].IsWaterMark {
			continue
		}
		splatCnt++
		if len(blockDatas[n].SH1) > 0 {
			for i := range 9 {
				bts = append(bts, cmn.EncodeSpxSH(blockDatas[n].SH1[i]))
			}
			for range 15 {
				bts = append(bts, cmn.EncodeSplatSH(0.0))
			}
		} else if len(blockDatas[n].SH2) > 0 {
			for i := range 24 {
				bts = append(bts, cmn.EncodeSpxSH(blockDatas[n].SH2[i]))
			}
		} else {
			for range 24 {
				bts = append(bts, cmn.EncodeSplatSH(0.0))
			}
		}
	}
	if splatCnt == 0 {
		return
	} else if splatCnt < blockSplatCount {
		cntBytes := cmn.Uint32ToBytes(uint32(splatCnt))
		bts[0] = cntBytes[0]
		bts[1] = cntBytes[1]
		bts[2] = cntBytes[2]
		bts[3] = cntBytes[3]
	}

	if splatCnt >= MinCompressBlockSize {
		var err error
		switch compressType {
		case CT_XZ:
			bts, err = cmn.CompressXZ(bts)
		default:
			bts, err = cmn.CompressGzip(bts)
		}
		cmn.ExitOnError(err)
		blockByteLength := -((int32(compressType) << 28) | int32(len(bts)))
		_, err = writer.Write(cmn.Int32ToBytes(blockByteLength))
		cmn.ExitOnError(err)
		_, err = writer.Write(bts)
		cmn.ExitOnError(err)
	} else {
		blockByteLength := int32(len(bts))
		_, err := writer.Write(cmn.Int32ToBytes(blockByteLength))
		cmn.ExitOnError(err)
		_, err = writer.Write(bts)
		cmn.ExitOnError(err)
	}
}

func writeSpxBlockSH3(writer *bufio.Writer, blockDatas []*SplatData, compressType uint8) {
	blockSplatCount := len(blockDatas)
	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(blockSplatCount))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(BF_SH3)...)                  // 开放的块数据格式 3:sh3

	splatCnt := 0
	for n := range blockSplatCount {
		if blockDatas[n].IsWaterMark {
			continue
		}
		splatCnt++
		if len(blockDatas[n].SH3) > 0 {
			for i := range 21 {
				bts = append(bts, cmn.EncodeSpxSH(blockDatas[n].SH3[i]))
			}
		} else {
			for range 21 {
				bts = append(bts, cmn.EncodeSplatSH(0.0))
			}
		}
	}
	if splatCnt == 0 {
		return
	} else if splatCnt < blockSplatCount {
		cntBytes := cmn.Uint32ToBytes(uint32(splatCnt))
		bts[0] = cntBytes[0]
		bts[1] = cntBytes[1]
		bts[2] = cntBytes[2]
		bts[3] = cntBytes[3]
	}

	if splatCnt >= MinCompressBlockSize {
		var err error
		switch compressType {
		case CT_XZ:
			bts, err = cmn.CompressXZ(bts)
		default:
			bts, err = cmn.CompressGzip(bts)
		}
		cmn.ExitOnError(err)
		blockByteLength := -((int32(compressType) << 28) | int32(len(bts)))
		_, err = writer.Write(cmn.Int32ToBytes(blockByteLength))
		cmn.ExitOnError(err)
		_, err = writer.Write(bts)
		cmn.ExitOnError(err)
	} else {
		blockByteLength := int32(len(bts))
		_, err := writer.Write(cmn.Int32ToBytes(blockByteLength))
		cmn.ExitOnError(err)
		_, err = writer.Write(bts)
		cmn.ExitOnError(err)
	}
}

func writeSpxBlockSH3Webp(writer *bufio.Writer, blockDatas []*SplatData, shDegree uint8) {
	blockSplatCount := len(blockDatas)
	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(blockSplatCount))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(BF_SH3_WEBP)...)             // 开放的块数据格式 4:球谐系数3级（共15个）

	splatCnt := 0
	color0 := cmn.EncodeSplatSH(0.0)
	shRgba := make([]byte, 0)
	for n := range blockSplatCount {
		if blockDatas[n].IsWaterMark {
			continue
		}
		splatCnt++
		if len(blockDatas[n].SH1) > 0 {
			// 只有1级数据，输出1、2、3级别都是一样的结果
			for i := range 3 {
				shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH1[i*3+0]))
				shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH1[i*3+1]))
				shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH1[i*3+2]))
				shRgba = append(shRgba, 255)
			}
			for range 12 {
				shRgba = append(shRgba, color0, color0, color0, 255)
			}
		} else if len(blockDatas[n].SH3) > 0 {
			// 有全部3级数据
			switch shDegree {
			case 1:
				for i := range 3 {
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH2[i*3+0]))
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH2[i*3+1]))
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH2[i*3+2]))
					shRgba = append(shRgba, 255)
				}
				for range 12 {
					shRgba = append(shRgba, color0, color0, color0, 255)
				}
			case 2:
				for i := range 8 {
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH2[i*3+0]))
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH2[i*3+1]))
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH2[i*3+2]))
					shRgba = append(shRgba, 255)
				}
				for range 7 {
					shRgba = append(shRgba, color0, color0, color0, 255)
				}
			default:
				for i := range 8 {
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH2[i*3+0]))
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH2[i*3+1]))
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH2[i*3+2]))
					shRgba = append(shRgba, 255)
				}
				for i := range 7 {
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH3[i*3+0]))
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH3[i*3+1]))
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH3[i*3+2]))
					shRgba = append(shRgba, 255)
				}
			}

		} else if len(blockDatas[n].SH2) > 0 {
			// 只有1、2级数据
			switch shDegree {
			case 1:
				for i := range 3 {
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH2[i*3+0]))
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH2[i*3+1]))
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH2[i*3+2]))
					shRgba = append(shRgba, 255)
				}
				for range 12 {
					shRgba = append(shRgba, color0, color0, color0, 255)
				}
			default:
				for i := range 8 {
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH2[i*3+0]))
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH2[i*3+1]))
					shRgba = append(shRgba, cmn.EncodeSpxSH(blockDatas[n].SH2[i*3+2]))
					shRgba = append(shRgba, 255)
				}
				for range 7 {
					shRgba = append(shRgba, color0, color0, color0, 255)
				}
			}
		} else {
			// 无
			for range 15 {
				shRgba = append(shRgba, color0, color0, color0, 255)
			}
		}
	}
	if splatCnt == 0 {
		return
	} else if splatCnt < blockSplatCount {
		cntBytes := cmn.Uint32ToBytes(uint32(splatCnt))
		bts[0] = cntBytes[0]
		bts[1] = cntBytes[1]
		bts[2] = cntBytes[2]
		bts[3] = cntBytes[3]
	}

	webpBts, err := cmn.CompressWebp(shRgba, oArg.webpQuality)
	cmn.ExitOnError(err)
	bts = append(bts, webpBts...)
	cmn.ExitOnError(err)

	blockByteLength := int32(len(bts))
	_, err = writer.Write(cmn.Int32ToBytes(blockByteLength))
	cmn.ExitOnError(err)
	_, err = writer.Write(bts)
	cmn.ExitOnError(err)
}

func RgbaToSh(rgba []byte, splatCnt int, shCount int) []byte {
	var rs []byte

	for i := range splatCnt {
		for j := range shCount {
			rs = append(rs, rgba[i*shCount*4+j*4+0])
			rs = append(rs, rgba[i*shCount*4+j*4+1])
			rs = append(rs, rgba[i*shCount*4+j*4+2])
		}
	}

	return rs
}
