package gsplat

import (
	"bufio"
	"fmt"
	"gsbox/cmn"
	"math"
	"math/rand/v2"
	"os"
	"sort"
)

const BlockSize = 20480
const MinGzipBlockSize = 64

func WriteSpx(spxFile string, rows []*SplatData, comment string, shDegree int) {
	file, err := os.Create(spxFile)
	cmn.ExitOnError(err)
	defer file.Close()

	writer := bufio.NewWriter(file)
	if shDegree < 0 || shDegree > 3 {
		shDegree = 0 // 球谐系数范围[0,1,2,3]，其他值都按0处理
	}

	header := genSpx1Header(rows, comment, shDegree)
	_, err = writer.Write(header.ToSpx1Bytes())
	cmn.ExitOnError(err)

	var blockDatasList [][]*SplatData
	blockCnt := (int(header.SplatCount) + BlockSize - 1) / BlockSize
	for i := range blockCnt {
		blockDatas := make([]*SplatData, 0)
		max := min(i*BlockSize+BlockSize, int(header.SplatCount))
		for n := i * BlockSize; n < max; n++ {
			blockDatas = append(blockDatas, rows[n])
		}
		writeSpxBlockSplat20(writer, blockDatas, len(blockDatas))
		blockDatasList = append(blockDatasList, blockDatas)
	}

	if shDegree >= 1 {
		for i := range blockDatasList {
			writeSpxBlockSH1(writer, blockDatasList[i])
		}
	}
	if shDegree >= 2 {
		for i := range blockDatasList {
			writeSpxBlockSH2(writer, blockDatasList[i])
		}
	}
	if shDegree == 3 {
		for i := range blockDatasList {
			writeSpxBlockSH3(writer, blockDatasList[i])
		}
	}

	err = writer.Flush()
	cmn.ExitOnError(err)
}

func genSpx1Header(datas []*SplatData, comment string, shDegree int) *SpxHeader {

	header := new(SpxHeader)
	header.Fixed = "spx"
	header.Version = 1
	header.SplatCount = int32(len(datas))

	header.CreateDate = cmn.GetSystemDateYYYYMMDD() // 创建日期
	header.CreaterId = 1202056903                   // 0:官方默认识别号，（这里参考阿佩里常数1.202056903159594…以示区分，此常数由瑞士数学家罗杰·阿佩里在1978年证明其无理数性质而闻名）
	header.ExclusiveId = 0                          // 0:官方开放格式的识别号
	header.ShDegree = int32(shDegree)
	header.Reserve1 = rand.Float32() * 10
	header.Reserve2 = rand.Float32() * 20
	del, comment := cmn.RemoveNonASCII(comment)
	if del {
		fmt.Println("[WARN] The existing non-ASCII characters in the comment have been removed!")
	}
	header.Comment = comment // 注释
	if header.Comment == "" {
		header.Comment = "created by gsbox"
	}

	if len(datas) > 0 {
		minX := float64(datas[0].PositionX)
		minY := float64(datas[0].PositionY)
		minZ := float64(datas[0].PositionZ)
		maxX := float64(datas[0].PositionX)
		maxY := float64(datas[0].PositionY) // 模型是Y轴倒立
		maxZ := float64(datas[0].PositionZ)

		for i := 1; i < len(datas); i++ {
			minX = math.Min(minX, float64(datas[i].PositionX))
			minY = math.Max(minY, float64(datas[i].PositionY))
			minZ = math.Min(minZ, float64(datas[i].PositionZ))
			maxX = math.Max(maxX, float64(datas[i].PositionX))
			maxY = math.Min(maxY, float64(datas[i].PositionY))
			maxZ = math.Max(maxZ, float64(datas[i].PositionZ))
		}
		header.MinX = cmn.ToFloat32(minX)
		header.MaxX = cmn.ToFloat32(maxX)
		header.MinY = cmn.ToFloat32(minY)
		header.MaxY = cmn.ToFloat32(maxY)
		header.MinZ = cmn.ToFloat32(minZ)
		header.MaxZ = cmn.ToFloat32(maxZ)
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
		bts = append(bts, cmn.EncodeFloat32ToBytes3(blockDatas[n].PositionX)...)
	}
	for n := range blockSplatCount {
		bts = append(bts, cmn.EncodeFloat32ToBytes3(blockDatas[n].PositionY)...)
	}
	for n := range blockSplatCount {
		bts = append(bts, cmn.EncodeFloat32ToBytes3(blockDatas[n].PositionZ)...)
	}
	for n := range blockSplatCount {
		bts = append(bts, cmn.EncodeFloat32ToByte(blockDatas[n].ScaleX))
	}
	for n := range blockSplatCount {
		bts = append(bts, cmn.EncodeFloat32ToByte(blockDatas[n].ScaleY))
	}
	for n := range blockSplatCount {
		bts = append(bts, cmn.EncodeFloat32ToByte(blockDatas[n].ScaleZ))
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
	for n := range blockSplatCount {
		bts = append(bts, blockDatas[n].RotationW)
	}

	if blockSplatCount >= MinGzipBlockSize {
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

	if len(blockDatas[0].SH1) > 0 {
		// 有SH1的数据
		for n := range blockSplatCount {
			bts = append(bts, blockDatas[n].SH1...)
		}
	} else {
		// 没有SH1的数据
		bts = append(bts, make([]byte, blockSplatCount*9)...) // n*9个0
	}

	if blockSplatCount >= MinGzipBlockSize {
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

	if len(blockDatas[0].SH2) > 0 {
		// 有SH2的数据
		for n := range blockSplatCount {
			bts = append(bts, blockDatas[n].SH2...)
		}
	} else {
		// 没有SH2的数据
		bts = append(bts, make([]byte, blockSplatCount*15)...) // n*15个0
	}

	if blockSplatCount >= MinGzipBlockSize {
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

	if len(blockDatas[0].SH3) > 0 {
		// 有SH3的数据
		for n := range blockSplatCount {
			bts = append(bts, blockDatas[n].SH3...)
		}
	} else {
		// 没有SH3的数据
		bts = append(bts, make([]byte, blockSplatCount*21)...) // n*15个0
	}

	if blockSplatCount >= MinGzipBlockSize {
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
