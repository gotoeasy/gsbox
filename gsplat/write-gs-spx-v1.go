package gsplat

import (
	"bufio"
	"gsbox/cmn"
	"log"
	"math"
	"os"
	"sort"
)

func WriteSpxV1(spxFile string, rows []*SplatData, comment string, shDegree int, flag1 uint8, flag2 uint8, flag3 uint8) {
	file, err := os.Create(spxFile)
	cmn.ExitOnError(err)
	defer file.Close()

	log.Println("[Info] output shDegree:", shDegree)
	writer := bufio.NewWriter(file)

	blockSize := cmn.StringToInt(Args.GetArgIgnorecase("-bs", "--block-size"), 20480)
	if blockSize <= 0 {
		blockSize = len(rows) // 所有数据放到一个块
	} else if blockSize < MinCompressBlockSize {
		blockSize = MinCompressBlockSize // 最小
	} else if blockSize > 512000 {
		blockSize = 512000 // 最大512000
	}

	header := genSpxHeader(rows, comment, shDegree, flag1, flag2, flag3)
	_, err = writer.Write(header.ToBytes())
	cmn.ExitOnError(err)

	var blockDatasList [][]*SplatData
	blockCnt := (int(header.SplatCount) + blockSize - 1) / blockSize
	for i := range blockCnt {
		blockDatas := make([]*SplatData, 0)
		max := min(i*blockSize+blockSize, int(header.SplatCount))
		for n := i * blockSize; n < max; n++ {
			blockDatas = append(blockDatas, rows[n])
		}
		writeSpxBlockSplat20(writer, blockDatas, len(blockDatas))
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

func genSpxHeader(datas []*SplatData, comment string, shDegree int, flag1 uint8, flag2 uint8, flag3 uint8) *SpxHeader {

	header := new(SpxHeader)
	header.Fixed = "spx"
	header.Version = 1
	header.SplatCount = int32(len(datas))

	header.CreateDate = cmn.GetSystemDateYYYYMMDD() // 创建日期
	header.CreaterId = ID1202056903                 // 0:官方默认识别号，（这里参考阿佩里常数1.202056903159594…以示区分，此常数由瑞士数学家罗杰·阿佩里在1978年证明其无理数性质而闻名）
	header.ExclusiveId = 0                          // 0:官方开放格式的识别号
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

func writeSpxBlockSplat20(writer *bufio.Writer, blockDatas []*SplatData, blockSplatCount int) {
	sort.Slice(blockDatas, func(i, j int) bool {
		return blockDatas[i].PositionY < blockDatas[j].PositionY // 坐标分别占3字节，按其中任一排序以更利于压缩
	})

	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(blockSplatCount))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(20)...)                      // 开放的块数据格式 20:splat20重排

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
		bts, err := cmn.GzipBytes(bts)
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

func writeSpxBlockSH1(writer *bufio.Writer, blockDatas []*SplatData) {
	blockSplatCount := len(blockDatas)
	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(blockSplatCount))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(1)...)                       // 开放的块数据格式 1:sh1

	for n := range blockSplatCount {
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

	if blockSplatCount >= MinCompressBlockSize {
		bts, err := cmn.GzipBytes(bts)
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

func writeSpxBlockSH2(writer *bufio.Writer, blockDatas []*SplatData) {
	blockSplatCount := len(blockDatas)
	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(blockSplatCount))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(2)...)                       // 开放的块数据格式 2:sh2

	for n := range blockSplatCount {
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

	if blockSplatCount >= MinCompressBlockSize {
		bts, err := cmn.GzipBytes(bts)
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

func writeSpxBlockSH3(writer *bufio.Writer, blockDatas []*SplatData) {
	blockSplatCount := len(blockDatas)
	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(uint32(blockSplatCount))...) // 块中的高斯点个数
	bts = append(bts, cmn.Uint32ToBytes(3)...)                       // 开放的块数据格式 3:sh3

	for n := range blockSplatCount {
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

	if blockSplatCount >= MinCompressBlockSize {
		bts, err := cmn.GzipBytes(bts)
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
