package gsplat

import (
	"bufio"
	"gsbox/cmn"
	"log"
	"math"
	"os"
)

func WriteSpxV3(spxFile string, rows []*SplatData, comment string, shDegree uint8) {
	file, err := os.Create(spxFile)
	cmn.ExitOnError(err)
	defer file.Close()

	writer := bufio.NewWriter(file)

	inputBlockSize := Args.GetArgIgnorecase("-bs", "--block-size")
	blockSize := cmn.StringToInt(inputBlockSize, DefaultBlockSize)
	if cmn.EqualsIngoreCase(inputBlockSize, "max") {
		blockSize = MaxBlockSize // 支持 -bs max 写法
	} else if blockSize < MinCompressBlockSize || blockSize > MaxBlockSize {
		blockSize = DefaultBlockSize // 超出范围按默认看待
	}

	header := genSpxHeaderV3(rows, comment, shDegree)
	_, err = writer.Write(header.ToBytes())
	cmn.ExitOnError(err)

	bf := cmn.StringToInt(Args.GetArgIgnorecase("-bf", "--block-format"), BF_SPLAT220_WEBP)
	if bf != BF_SPLAT22 && bf != BF_SPLAT220_WEBP {
		bf = BF_SPLAT22 // 默认格式
	}
	log.Println("[Info] (Parameter) block format:", bf, BlockFormatDesc(bf))
	log.Println("[Info] (Parameter) block size:", blockSize)

	var compressType uint8 = CT_XZ // 默认xz
	ct := Args.GetArgIgnorecase("-ct", "--compress-type")
	if bf != BF_SPLAT220_WEBP {
		if cmn.EqualsIngoreCase(ct, "gzip") || ct == "0" {
			compressType = CT_GZIP
			log.Println("[Info] block compress type: gzip")
		} else {
			log.Println("[Info] block compress type: xz")
		}
	}

	for _, d := range rows {
		d.RotationW, d.RotationX, d.RotationY, d.RotationZ = cmn.NormalizeRotations(d.RotationW, d.RotationX, d.RotationY, d.RotationZ)
	}
	shCentroids, _ := ReWriteShByKmeans(rows)

	// 分块
	blockCnt := (int(header.SplatCount) + blockSize - 1) / blockSize
	var blockList [][]*SplatData
	for i := range blockCnt {
		blockDatas := make([]*SplatData, 0)
		max := min(i*blockSize+blockSize, int(header.SplatCount))
		for n := i * blockSize; n < max; n++ {
			blockDatas = append(blockDatas, rows[n])
		}
		blockList = append(blockList, blockDatas)
	}

	// 多块且最后块太小时并入前一块
	if len(blockList) > 1 && len(blockList[len(blockList)-1]) < MinWebpBlockSize {
		last := blockList[len(blockList)-1]
		blockList = blockList[:len(blockList)-1]
		blockList[len(blockList)-1] = append(blockList[len(blockList)-1], last...)
	}

	// 写文件
	palettesDone := shDegree == 0
	writeCnt := 0
	for _, blockDatas := range blockList {
		if bf == BF_SPLAT220_WEBP && (len(blockDatas) >= MinWebpBlockSize || Args.HasArgIgnorecase("-bf", "--block-format")) {
			// 默认提示WEBP编码且数据量够大，或强制参数要求WEBP编码
			writeSpxBF220_WEBP_V3(writer, blockDatas, shDegree)
		} else {
			writeSpxBF22_V3(writer, blockDatas, shDegree, compressType)
		}
		writeCnt += len(blockDatas)

		if !palettesDone && writeCnt >= DefaultBlockSize*4 {
			// 调色板插在较前处写入，避免中断下载后无法读取
			if bf == BF_SPLAT220_WEBP || compressType == CT_XZ {
				writePalettesWebp_V3(writer, shCentroids)
			} else {
				writePalettes_V3(writer, shCentroids, compressType)
			}
			palettesDone = true
		}
	}

	if !palettesDone {
		// 确保调色板存在时会写入
		if bf == BF_SPLAT220_WEBP || compressType == CT_XZ {
			writePalettesWebp_V3(writer, shCentroids)
		} else {
			writePalettes_V3(writer, shCentroids, compressType)
		}
	}

	err = writer.Flush()
	cmn.ExitOnError(err)
}

func genSpxHeaderV3(datas []*SplatData, comment string, shDegree uint8) *SpxHeader {
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
	header.Version = 3
	header.SplatCount = int32(len(datas))

	header.CreateDate = cmn.GetSystemDateYYYYMMDD() // 创建日期
	header.CreaterId = GetOutputCreaterId()         // 0:官方默认识别号，（这里参考阿佩里常数1.202056903159594…以示区分，此常数由瑞士数学家罗杰·阿佩里在1978年证明其无理数性质而闻名）
	header.ExclusiveId = GetOutputExclusiveId()     // 0:官方开放格式的识别号
	header.ShDegree = uint8(shDegree)               // 0,1,2,3
	header.Flag = f1 | f8                           // v2+
	header.MaxFlagValue = 0
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

func writeSpxBF22_V3(writer *bufio.Writer, blockDatas []*SplatData, shDegree uint8, compressType uint8) {

	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(len(blockDatas)))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(BF_SPLAT22)...)              // 开放的块数据格式 22

	var bs0 []byte
	var bs1 []byte
	var bs2 []byte
	for _, d := range blockDatas {
		b3 := cmn.EncodeSpxPositionUint24(cmn.EncodeLog(d.PositionX, 1))
		bs0 = append(bs0, b3[0])
		bs1 = append(bs1, b3[1])
		bs2 = append(bs2, b3[2])
	}
	for _, d := range blockDatas {
		b3 := cmn.EncodeSpxPositionUint24(cmn.EncodeLog(d.PositionY, 1))
		bs0 = append(bs0, b3[0])
		bs1 = append(bs1, b3[1])
		bs2 = append(bs2, b3[2])
	}
	for _, d := range blockDatas {
		b3 := cmn.EncodeSpxPositionUint24(cmn.EncodeLog(d.PositionZ, 1))
		bs0 = append(bs0, b3[0])
		bs1 = append(bs1, b3[1])
		bs2 = append(bs2, b3[2])
	}
	bts = append(bts, bs0...)
	bts = append(bts, bs1...)
	bts = append(bts, bs2...)

	for _, d := range blockDatas {
		bts = append(bts, cmn.EncodeSpxScale(d.ScaleX))
	}
	for _, d := range blockDatas {
		bts = append(bts, cmn.EncodeSpxScale(d.ScaleY))
	}
	for _, d := range blockDatas {
		bts = append(bts, cmn.EncodeSpxScale(d.ScaleZ))
	}
	for _, d := range blockDatas {
		bts = append(bts, d.ColorR)
	}
	for _, d := range blockDatas {
		bts = append(bts, d.ColorG)
	}
	for _, d := range blockDatas {
		bts = append(bts, d.ColorB)
	}
	for _, d := range blockDatas {
		bts = append(bts, d.ColorA)
	}

	for _, d := range blockDatas {
		bts = append(bts, d.RotationW)
	}
	for _, d := range blockDatas {
		bts = append(bts, d.RotationX)
	}
	for _, d := range blockDatas {
		bts = append(bts, d.RotationY)
	}
	for _, d := range blockDatas {
		bts = append(bts, d.RotationZ)
	}

	if shDegree > 0 {
		for _, d := range blockDatas {
			bts = append(bts, byte(d.ShPaletteIdx&0xFF))
		}
		for _, d := range blockDatas {
			bts = append(bts, byte(d.ShPaletteIdx>>8))
		}
	}

	if len(blockDatas) >= MinCompressBlockSize {
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

func writeSpxBF220_WEBP_V3(writer *bufio.Writer, blockDatas []*SplatData, shDegree uint8) {
	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(len(blockDatas)))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(BF_SPLAT220_WEBP)...)        // 开放的块数据格式 220

	bsTmp := make([]byte, 0)
	bs1 := make([]byte, 0)
	bs2 := make([]byte, 0)
	bs3 := make([]byte, 0)
	for _, d := range blockDatas {
		b3x := cmn.EncodeSpxPositionUint24(cmn.EncodeLog(d.PositionX, 1))
		b3y := cmn.EncodeSpxPositionUint24(cmn.EncodeLog(d.PositionY, 1))
		b3z := cmn.EncodeSpxPositionUint24(cmn.EncodeLog(d.PositionZ, 1))
		bs1 = append(bs1, b3x[0], b3y[0], b3z[0], 255)
		bs2 = append(bs2, b3x[1], b3y[1], b3z[1], 255)
		bs3 = append(bs3, b3x[2], b3y[2], b3z[2], 255)
	}
	bsTmp = append(bsTmp, bs1...)
	bsTmp = append(bsTmp, bs2...)
	bsTmp = append(bsTmp, bs3...)
	bsTmp, err := cmn.CompressWebp(bsTmp)
	cmn.ExitOnError(err)
	bts = append(bts, cmn.Uint32ToBytes(uint32(len(bsTmp)))...)
	bts = append(bts, bsTmp...)

	bsTmp = make([]byte, 0)
	for _, d := range blockDatas {
		bsTmp = append(bsTmp, cmn.EncodeSpxScale(d.ScaleX), cmn.EncodeSpxScale(d.ScaleY), cmn.EncodeSpxScale(d.ScaleZ), 255)
	}
	bsTmp, err = cmn.CompressWebp(bsTmp)
	cmn.ExitOnError(err)
	bts = append(bts, cmn.Uint32ToBytes(uint32(len(bsTmp)))...)
	bts = append(bts, bsTmp...)

	bsTmp = make([]byte, 0)
	for _, d := range blockDatas {
		bsTmp = append(bsTmp, d.ColorR, d.ColorG, d.ColorB, d.ColorA)
	}
	bsTmp, err = cmn.CompressWebp(bsTmp)
	cmn.ExitOnError(err)
	bts = append(bts, cmn.Uint32ToBytes(uint32(len(bsTmp)))...)
	bts = append(bts, bsTmp...)

	bsTmp = make([]byte, 0)
	for _, d := range blockDatas {
		r0, r1, r2, ri := cmn.SogEncodeRotations(d.RotationW, d.RotationX, d.RotationY, d.RotationZ)
		bsTmp = append(bsTmp, r0, r1, r2, ri)
	}
	bsTmp, err = cmn.CompressWebp(bsTmp)
	cmn.ExitOnError(err)
	bts = append(bts, cmn.Uint32ToBytes(uint32(len(bsTmp)))...)
	bts = append(bts, bsTmp...)

	if shDegree > 0 {
		bsTmp = make([]byte, 0)
		for _, d := range blockDatas {
			bsTmp = append(bsTmp, byte(d.ShPaletteIdx&0xFF), byte(d.ShPaletteIdx>>8), 0, 255)
		}
		bsTmp, err = cmn.CompressWebp(bsTmp)
		cmn.ExitOnError(err)
		bts = append(bts, cmn.Uint32ToBytes(uint32(len(bsTmp)))...)
		bts = append(bts, bsTmp...)
	}

	blockByteLength := int32(len(bts))
	_, err = writer.Write(cmn.Int32ToBytes(blockByteLength))
	cmn.ExitOnError(err)
	_, err = writer.Write(bts)
	cmn.ExitOnError(err)
}

func writePalettes_V3(writer *bufio.Writer, shCentroids []byte, compressType uint8) {
	log.Println("[Info] palettes block format:", BF_SH_PALETTES, BlockFormatDesc(BF_SH_PALETTES))

	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(0)...)              // 占位
	bts = append(bts, cmn.Uint32ToBytes(BF_SH_PALETTES)...) // 球谐系数调色板
	bts = append(bts, shCentroids...)                       // 数据

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
}

func writePalettesWebp_V3(writer *bufio.Writer, shCentroids []byte) {
	log.Println("[Info] palettes block format:", BF_SH_PALETTES_WEBP, BlockFormatDesc(BF_SH_PALETTES_WEBP))

	webpBytes, err := cmn.CompressWebpByWidthHeight(shCentroids, 960, 1024)
	cmn.ExitOnError(err)

	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(0)...)                   // 占位
	bts = append(bts, cmn.Uint32ToBytes(BF_SH_PALETTES_WEBP)...) // 球谐系数调色板, WEBP压缩
	bts = append(bts, webpBytes...)                              // 数据

	blockByteLength := int32(len(bts))
	_, err = writer.Write(cmn.Int32ToBytes(blockByteLength))
	cmn.ExitOnError(err)
	_, err = writer.Write(bts)
	cmn.ExitOnError(err)
}
