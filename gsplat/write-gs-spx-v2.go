package gsplat

import (
	"bufio"
	"gsbox/cmn"
	"log"
	"math"
	"os"
)

func WriteSpxV2(spxFile string, rows []*SplatData, comment string, shDegree int) {
	file, err := os.Create(spxFile)
	cmn.ExitOnError(err)
	defer file.Close()

	log.Println("[Info] output shDegree:", shDegree)
	writer := bufio.NewWriter(file)

	blockSize := cmn.StringToInt(Args.GetArgIgnorecase("-bs", "--block-size"), DefaultBlockSize)
	if blockSize < MinCompressBlockSize || blockSize > MaxBlockSize {
		blockSize = MaxBlockSize // 默认及超出范围都按最大看待
	}

	header := genSpxHeaderV2(rows, comment, shDegree)
	_, err = writer.Write(header.ToBytes())
	cmn.ExitOnError(err)

	bf := cmn.StringToInt(Args.GetArgIgnorecase("-bf", "--block-format"), 19)
	if bf != 20 && bf != 16 {
		bf = 19 // 默认splat19格式
	}
	log.Println("[Info] (Parameter) data block format:", bf)
	log.Println("[Info] (Parameter) block size:", blockSize)

	// var compressType uint8 = 0 // 默认gzip
	// 试验下来 zstd 比 gzip 的压缩率更差，暂不支持
	// ct := Args.GetArgIgnorecase("-ct", "--compress-type")
	// if cmn.EqualsIngoreCase(ct, "zstd") {
	// 	compressType = 1
	// 	log.Println("[Info] block compress type: zstd")
	// } else {
	// 	log.Println("[Info] block compress type: gzip")
	// }

	var blockDatasList [][]*SplatData
	blockCnt := (int(header.SplatCount) + blockSize - 1) / blockSize
	for i := range blockCnt {
		blockDatas := make([]*SplatData, 0)
		max := min(i*blockSize+blockSize, int(header.SplatCount))
		for n := i * blockSize; n < max; n++ {
			blockDatas = append(blockDatas, rows[n])
		}
		switch bf {
		case 20:
			// splat20 格式，优势不够突出，spx v1版本使用
			writeSpxBlockSplat20(writer, blockDatas, len(blockDatas))
		case 16:
			// splat16 格式，压缩率较好，有些情况误差较大，但若肉眼可接受可能换取更高的压缩率，供手动选用
			writeSpxBlockSplat16(writer, blockDatas, len(blockDatas))
		default:
			// splat19 格式，压缩率等综合表现比较优秀，默认使用
			writeSpxBlockSplat19(writer, blockDatas, len(blockDatas))
		}
		blockDatasList = append(blockDatasList, blockDatas)
	}

	switch shDegree {
	case 1:
		for i := range blockDatasList {
			writeSpxBlockSH1(writer, blockDatasList[i])
		}
	case 2:
		for i := range blockDatasList {
			writeSpxBlockSH2(writer, blockDatasList[i])
		}
	case 3:
		for i := range blockDatasList {
			writeSpxBlockSH2(writer, blockDatasList[i])
			writeSpxBlockSH3(writer, blockDatasList[i])
		}
	}

	err = writer.Flush()
	cmn.ExitOnError(err)
}

func genSpxHeaderV2(datas []*SplatData, comment string, shDegree int) *SpxHeader {
	var f1 uint8 = 0 // 是否Y轴倒立模型
	var f8 uint8 = 0 // 是否大场景

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
	header.CreaterId = ID1202056903                 // 0:官方默认识别号，（这里参考阿佩里常数1.202056903159594…以示区分，此常数由瑞士数学家罗杰·阿佩里在1978年证明其无理数性质而闻名）
	header.ExclusiveId = 0                          // 0:官方开放格式的识别号
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
		header.Comment = "created by gsbox " + cmn.VER + " https://github.com/gotoeasy/gsbox"
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

func writeSpxBlockSplat19(writer *bufio.Writer, blockDatas []*SplatData, blockSplatCount int) {
	SortBlockDatas4Compress(blockDatas)
	for n := range blockSplatCount {
		blockDatas[n].RotationW, blockDatas[n].RotationX, blockDatas[n].RotationY, blockDatas[n].RotationZ = cmn.NormalizeRotations(blockDatas[n].RotationW, blockDatas[n].RotationX, blockDatas[n].RotationY, blockDatas[n].RotationZ)
	}

	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(blockSplatCount))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(19)...)                      // 开放的块数据格式 19

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
		bts, err := cmn.CompressGzip(bts)
		cmn.ExitOnError(err)
		blockByteLength := -int32(len(bts))
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

func writeSpxBlockSplat16(writer *bufio.Writer, blockDatas []*SplatData, blockSplatCount int) {
	mm := ComputeXyzMinMax(blockDatas)
	SortBlockDatas4Compress(blockDatas)
	for n := range blockSplatCount {
		blockDatas[n].RotationW, blockDatas[n].RotationX, blockDatas[n].RotationY, blockDatas[n].RotationZ = cmn.NormalizeRotations(blockDatas[n].RotationW, blockDatas[n].RotationX, blockDatas[n].RotationY, blockDatas[n].RotationZ)
	}

	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(blockSplatCount))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(16)...)                      // 开放的块数据格式 16
	bts = append(bts, cmn.Float32ToBytes(mm.MinX)...)
	bts = append(bts, cmn.Float32ToBytes(mm.MaxX)...)
	bts = append(bts, cmn.Float32ToBytes(mm.MinY)...)
	bts = append(bts, cmn.Float32ToBytes(mm.MaxY)...)
	bts = append(bts, cmn.Float32ToBytes(mm.MinZ)...)
	bts = append(bts, cmn.Float32ToBytes(mm.MaxZ)...)

	var bs0 []byte
	var bs1 []byte
	for n := range blockSplatCount {
		b2 := cmn.EncodeSpxPositionUint16(blockDatas[n].PositionX, mm.MinX, mm.MaxX)
		bs0 = append(bs0, b2[0])
		bs1 = append(bs1, b2[1])
	}
	for n := range blockSplatCount {
		b2 := cmn.EncodeSpxPositionUint16(blockDatas[n].PositionY, mm.MinY, mm.MaxY)
		bs0 = append(bs0, b2[0])
		bs1 = append(bs1, b2[1])
	}
	for n := range blockSplatCount {
		b2 := cmn.EncodeSpxPositionUint16(blockDatas[n].PositionZ, mm.MinZ, mm.MaxZ)
		bs0 = append(bs0, b2[0])
		bs1 = append(bs1, b2[1])
	}
	bts = append(bts, bs0...)
	bts = append(bts, bs1...)

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
		bts, err := cmn.CompressGzip(bts)
		cmn.ExitOnError(err)
		blockByteLength := -int32(len(bts))
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
