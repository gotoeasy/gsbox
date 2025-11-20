package gsplat

import (
	"bufio"
	"gsbox/cmn"
	"log"
	"math"
	"os"
)

// Deprecated，旧版将废弃，推荐使用新版本
func WriteSpxOpenV2(spxFile string, rows []*SplatData, comment string, shDegree uint8) {
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

	header := genSpxHeaderV2(rows, comment, shDegree)
	_, err = writer.Write(header.ToBytes())
	cmn.ExitOnError(err)

	bf := cmn.StringToInt(Args.GetArgIgnorecase("-bf", "--block-format"), BF_SPLAT190_WEBP)
	if bf != BF_SPLAT19 && bf != BF_SPLAT20 && bf != BF_SPLAT190_WEBP && bf != BF_SPLAT10190_WEBP && bf != BF_SPLAT10019 {
		bf = BF_SPLAT190_WEBP // 默认格式
	}
	log.Println("[Info] (parameter) bf:", bf, BlockFormatDesc(bf))
	log.Println("[Info] (parameter) bs:", blockSize, "(block size)")

	var compressType uint8 = CT_XZ // 默认xz
	ct := Args.GetArgIgnorecase("-ct", "--compress-type")
	if cmn.EqualsIngoreCase(ct, "gzip") || ct == "0" {
		compressType = CT_GZIP
		log.Println("[Info] block compress type: gzip")
	} else {
		log.Println("[Info] block compress type: xz")
	}

	logTimes := min(max(0, uint8(cmn.StringToInt(Args.GetArgIgnorecase("-l", "--log-times"), 1))), 9) // [TEST]有效范围0~9，默认1
	if (bf == BF_SPLAT10019 || bf == BF_SPLAT10190_WEBP) && logTimes > 0 {
		log.Println("[Info] log encoding times:", logTimes)
	}

	var blockDatasList [][]*SplatData
	blockCnt := (int(header.SplatCount) + blockSize - 1) / blockSize
	for i := range blockCnt {
		blockDatas := make([]*SplatData, 0)
		max := min(i*blockSize+blockSize, int(header.SplatCount))
		for n := i * blockSize; n < max; n++ {
			blockDatas = append(blockDatas, rows[n])
		}
		blockSplatCount := len(blockDatas)
		switch bf {
		case BF_SPLAT20:
			// splat20 格式，优势不够突出，spx v1版本使用
			writeSpxBlockSplat20(writer, blockDatas, blockSplatCount, compressType)
		case BF_SPLAT190_WEBP:
			if blockSplatCount >= MinWebpBlockSize {
				//  数据够多时才一定使用 webp 编码格式
				writeSpxBlockSplat190Webp(writer, blockDatas, blockSplatCount)
			} else {
				// 数据较少时，切换使用 splat19 格式
				writeSpxBlockSplat19(writer, blockDatas, blockSplatCount, compressType)
			}
		case BF_SPLAT10190_WEBP:
			if blockSplatCount >= MinWebpBlockSize {
				//  数据够多时才一定使用 webp 编码格式
				writeSpxBlockSplat10190Webp(writer, blockDatas, blockSplatCount, logTimes)
			} else {
				// 数据较少时，切换使用 10019 格式
				writeSpxBlockSplat10019(writer, blockDatas, blockSplatCount, compressType, logTimes)
			}
		case BF_SPLAT19:
			// splat19 格式，压缩率等综合表现好
			writeSpxBlockSplat19(writer, blockDatas, blockSplatCount, compressType)
		default:
			// splat10019 格式，压缩率等综合表现比较优秀，默认使用
			writeSpxBlockSplat10019(writer, blockDatas, blockSplatCount, compressType, logTimes)
		}
		blockDatasList = append(blockDatasList, blockDatas)
	}

	if shDegree > 0 && bf == BF_SPLAT190_WEBP {
		for i := range blockDatasList {
			writeSpxBlockSH3Webp(writer, blockDatasList[i], shDegree)
		}
	} else {

		switch shDegree {
		case 1:
			for i := range blockDatasList {
				writeSpxBlockSH1(writer, blockDatasList[i], compressType)
			}
		case 2:
			for i := range blockDatasList {
				writeSpxBlockSH2(writer, blockDatasList[i], compressType)
			}
		case 3:
			for i := range blockDatasList {
				writeSpxBlockSH2(writer, blockDatasList[i], compressType)
				writeSpxBlockSH3(writer, blockDatasList[i], compressType)
			}
		}
	}

	err = writer.Flush()
	cmn.ExitOnError(err)
}

func genSpxHeaderV2(datas []*SplatData, comment string, shDegree uint8) *SpxHeader {
	var f1 uint8 = 0            // 是否Y轴倒立模型
	var f8 uint8 = GetLodFlag() // 是否大场景

	if !Args.HasCmd("join") && inputSpxHeader != nil && inputSpxHeader.Version > 1 {
		if inputSpxHeader.IsInverted() {
			f1 = 1 << 7
		}
	}

	if Args.HasArgIgnorecase("-f1", "--is-inverted") {
		F1 := Args.GetArgIgnorecase("-f1", "--is-inverted")
		if cmn.EqualsIngoreCase(F1, "true") || cmn.EqualsIngoreCase(F1, "yes") || cmn.EqualsIngoreCase(F1, "y") || cmn.EqualsIngoreCase(F1, "1") {
			f1 = 1 << 7
		} else {
			f1 = 0
		}
	}

	header := new(SpxHeader)
	header.Fixed = "spx"
	header.Version = 2
	header.SplatCount = int32(len(datas))

	header.CreateDate = cmn.GetSystemDateYYYYMMDD() // 创建日期
	header.CreaterId = CreaterIdOpen                // 创建者ID，（这里参考阿佩里常数1.202056903159594…以示区分，此常数由瑞士数学家罗杰·阿佩里在1978年证明其无理数性质而闻名）
	header.ExclusiveId = ExclusiveIdOpen            // 0:官方开放格式的识别号
	header.ShDegree = uint8(shDegree)               // 0,1,2,3
	header.Flag = f1 | f8                           // v2
	header.MaxFlagValue = 0                         // v2
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

func writeSpxBlockSplat19(writer *bufio.Writer, blockDatas []*SplatData, blockSplatCount int, compressType uint8) {
	for n := range blockSplatCount {
		blockDatas[n].RotationW, blockDatas[n].RotationX, blockDatas[n].RotationY, blockDatas[n].RotationZ = cmn.NormalizeRotations(blockDatas[n].RotationW, blockDatas[n].RotationX, blockDatas[n].RotationY, blockDatas[n].RotationZ)
	}

	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(blockSplatCount))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(BF_SPLAT19)...)              // 开放的块数据格式 19

	var bs0 []byte
	var bs1 []byte
	var bs2 []byte
	for n := range blockSplatCount {
		b3 := cmn.EncodeSpxPositionUint24(blockDatas[n].PositionX)
		bs0 = append(bs0, b3[0])
		bs1 = append(bs1, b3[1])
		bs2 = append(bs2, b3[2])
	}
	for n := range blockSplatCount {
		b3 := cmn.EncodeSpxPositionUint24(blockDatas[n].PositionY)
		bs0 = append(bs0, b3[0])
		bs1 = append(bs1, b3[1])
		bs2 = append(bs2, b3[2])
	}
	for n := range blockSplatCount {
		b3 := cmn.EncodeSpxPositionUint24(blockDatas[n].PositionZ)
		bs0 = append(bs0, b3[0])
		bs1 = append(bs1, b3[1])
		bs2 = append(bs2, b3[2])
	}
	bts = append(bts, bs0...)
	bts = append(bts, bs1...)
	bts = append(bts, bs2...)

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

func writeSpxBlockSplat10019(writer *bufio.Writer, blockDatas []*SplatData, blockSplatCount int, compressType uint8, logTimes uint8) {
	for n := range blockSplatCount {
		blockDatas[n].RotationW, blockDatas[n].RotationX, blockDatas[n].RotationY, blockDatas[n].RotationZ = cmn.NormalizeRotations(blockDatas[n].RotationW, blockDatas[n].RotationX, blockDatas[n].RotationY, blockDatas[n].RotationZ)
	}

	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(blockSplatCount))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(BF_SPLAT10019)...)           // 开放的块数据格式 60019
	bts = append(bts, logTimes, 0, 0, 0)                             // log编码次数(通常0~9)

	var bs0 []byte
	var bs1 []byte
	var bs2 []byte
	for n := range blockSplatCount {
		b3 := cmn.EncodeSpxPositionUint24(cmn.EncodeLog(blockDatas[n].PositionX, int(logTimes)))
		bs0 = append(bs0, b3[0])
		bs1 = append(bs1, b3[1])
		bs2 = append(bs2, b3[2])
	}
	for n := range blockSplatCount {
		b3 := cmn.EncodeSpxPositionUint24(cmn.EncodeLog(blockDatas[n].PositionY, int(logTimes)))
		bs0 = append(bs0, b3[0])
		bs1 = append(bs1, b3[1])
		bs2 = append(bs2, b3[2])
	}
	for n := range blockSplatCount {
		b3 := cmn.EncodeSpxPositionUint24(cmn.EncodeLog(blockDatas[n].PositionZ, int(logTimes)))
		bs0 = append(bs0, b3[0])
		bs1 = append(bs1, b3[1])
		bs2 = append(bs2, b3[2])
	}
	bts = append(bts, bs0...)
	bts = append(bts, bs1...)
	bts = append(bts, bs2...)

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

func writeSpxBlockSplat190Webp(writer *bufio.Writer, blockDatas []*SplatData, blockSplatCount int) {
	for n := range blockSplatCount {
		blockDatas[n].RotationW, blockDatas[n].RotationX, blockDatas[n].RotationY, blockDatas[n].RotationZ = cmn.NormalizeRotations(blockDatas[n].RotationW, blockDatas[n].RotationX, blockDatas[n].RotationY, blockDatas[n].RotationZ)
	}

	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(blockSplatCount))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(BF_SPLAT190_WEBP)...)        // 开放的块数据格式 190

	bsTmp := make([]byte, 0)
	bs1 := make([]byte, 0)
	bs2 := make([]byte, 0)
	bs3 := make([]byte, 0)
	for n := range blockSplatCount {
		b3x := cmn.EncodeSpxPositionUint24(blockDatas[n].PositionX)
		b3y := cmn.EncodeSpxPositionUint24(blockDatas[n].PositionY)
		b3z := cmn.EncodeSpxPositionUint24(blockDatas[n].PositionZ)
		bs1 = append(bs1, b3x[0], b3y[0], b3z[0], 255)
		bs2 = append(bs2, b3x[1], b3y[1], b3z[1], 255)
		bs3 = append(bs3, b3x[2], b3y[2], b3z[2], 255)
	}
	bsTmp = append(bsTmp, bs1...)
	bsTmp = append(bsTmp, bs2...)
	bsTmp = append(bsTmp, bs3...)
	bsTmp, err := cmn.CompressWebp(bsTmp, oArg.webpQuality)
	cmn.ExitOnError(err)
	bts = append(bts, cmn.Uint32ToBytes(uint32(len(bsTmp)))...)
	bts = append(bts, bsTmp...)

	bsTmp = make([]byte, 0)
	for n := range blockSplatCount {
		bsTmp = append(bsTmp, cmn.EncodeSpxScale(blockDatas[n].ScaleX), cmn.EncodeSpxScale(blockDatas[n].ScaleY), cmn.EncodeSpxScale(blockDatas[n].ScaleZ), 255)
	}
	bsTmp, err = cmn.CompressWebp(bsTmp, oArg.webpQuality)
	cmn.ExitOnError(err)
	bts = append(bts, cmn.Uint32ToBytes(uint32(len(bsTmp)))...)
	bts = append(bts, bsTmp...)

	bsTmp = make([]byte, 0)
	for n := range blockSplatCount {
		bsTmp = append(bsTmp, blockDatas[n].ColorR, blockDatas[n].ColorG, blockDatas[n].ColorB, blockDatas[n].ColorA)
	}
	bsTmp, err = cmn.CompressWebp(bsTmp, oArg.webpQuality)
	cmn.ExitOnError(err)
	bts = append(bts, cmn.Uint32ToBytes(uint32(len(bsTmp)))...)
	bts = append(bts, bsTmp...)

	bsTmp = make([]byte, 0)
	for n := range blockSplatCount {
		bsTmp = append(bsTmp, blockDatas[n].RotationX, blockDatas[n].RotationY, blockDatas[n].RotationZ, 255)
	}
	bsTmp, err = cmn.CompressWebp(bsTmp, oArg.webpQuality)
	cmn.ExitOnError(err)
	bts = append(bts, cmn.Uint32ToBytes(uint32(len(bsTmp)))...)
	bts = append(bts, bsTmp...)

	blockByteLength := int32(len(bts))
	_, err = writer.Write(cmn.Int32ToBytes(blockByteLength))
	cmn.ExitOnError(err)
	_, err = writer.Write(bts)
	cmn.ExitOnError(err)
}

func writeSpxBlockSplat10190Webp(writer *bufio.Writer, blockDatas []*SplatData, blockSplatCount int, logTimes uint8) {
	for n := range blockSplatCount {
		blockDatas[n].RotationW, blockDatas[n].RotationX, blockDatas[n].RotationY, blockDatas[n].RotationZ = cmn.NormalizeRotations(blockDatas[n].RotationW, blockDatas[n].RotationX, blockDatas[n].RotationY, blockDatas[n].RotationZ)
	}

	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(blockSplatCount))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(BF_SPLAT10190_WEBP)...)      // 开放的块数据格式 60190
	bts = append(bts, logTimes, 0, 0, 0)                             // log编码次数(通常0~9)

	bsTmp := make([]byte, 0)
	bs1 := make([]byte, 0)
	bs2 := make([]byte, 0)
	bs3 := make([]byte, 0)
	for n := range blockSplatCount {
		b3x := cmn.EncodeSpxPositionUint24(cmn.EncodeLog(blockDatas[n].PositionX, int(logTimes)))
		b3y := cmn.EncodeSpxPositionUint24(cmn.EncodeLog(blockDatas[n].PositionY, int(logTimes)))
		b3z := cmn.EncodeSpxPositionUint24(cmn.EncodeLog(blockDatas[n].PositionZ, int(logTimes)))
		bs1 = append(bs1, b3x[0], b3y[0], b3z[0], 255)
		bs2 = append(bs2, b3x[1], b3y[1], b3z[1], 255)
		bs3 = append(bs3, b3x[2], b3y[2], b3z[2], 255)
	}
	bsTmp = append(bsTmp, bs1...)
	bsTmp = append(bsTmp, bs2...)
	bsTmp = append(bsTmp, bs3...)
	bsTmp, err := cmn.CompressWebp(bsTmp, oArg.webpQuality)
	cmn.ExitOnError(err)
	bts = append(bts, cmn.Uint32ToBytes(uint32(len(bsTmp)))...)
	bts = append(bts, bsTmp...)

	bsTmp = make([]byte, 0)
	for n := range blockSplatCount {
		bsTmp = append(bsTmp, cmn.EncodeSpxScale(blockDatas[n].ScaleX), cmn.EncodeSpxScale(blockDatas[n].ScaleY), cmn.EncodeSpxScale(blockDatas[n].ScaleZ), 255)
	}
	bsTmp, err = cmn.CompressWebp(bsTmp, oArg.webpQuality)
	cmn.ExitOnError(err)
	bts = append(bts, cmn.Uint32ToBytes(uint32(len(bsTmp)))...)
	bts = append(bts, bsTmp...)

	bsTmp = make([]byte, 0)
	for n := range blockSplatCount {
		bsTmp = append(bsTmp, blockDatas[n].ColorR, blockDatas[n].ColorG, blockDatas[n].ColorB, blockDatas[n].ColorA)
	}
	bsTmp, err = cmn.CompressWebp(bsTmp, oArg.webpQuality)
	cmn.ExitOnError(err)
	bts = append(bts, cmn.Uint32ToBytes(uint32(len(bsTmp)))...)
	bts = append(bts, bsTmp...)

	bsTmp = make([]byte, 0)
	for n := range blockSplatCount {
		bsTmp = append(bsTmp, blockDatas[n].RotationX, blockDatas[n].RotationY, blockDatas[n].RotationZ, 255)
	}
	bsTmp, err = cmn.CompressWebp(bsTmp, oArg.webpQuality)
	cmn.ExitOnError(err)
	bts = append(bts, cmn.Uint32ToBytes(uint32(len(bsTmp)))...)
	bts = append(bts, bsTmp...)

	blockByteLength := int32(len(bts))
	_, err = writer.Write(cmn.Int32ToBytes(blockByteLength))
	cmn.ExitOnError(err)
	_, err = writer.Write(bts)
	cmn.ExitOnError(err)
}
