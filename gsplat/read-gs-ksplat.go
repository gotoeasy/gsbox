package gsplat

import (
	"errors"
	"fmt"
	"gsbox/cmn"
	"math"
	"os"
)

/** ksplat 文件的主头信息。主头信息长度为 4096 字节，但只有少数字段被使用。未使用的空间可能是为未来的扩展预留的 */
type KsplatHeader struct {
	/** 主版本号 */
	MajorVersion uint8
	/** 次版本号 */
	MinorVersion uint8
	/** 分段数量 */
	SectionCount int
	/** 点（splat）的总数 */
	SplatCount int
	/** 压缩模式，代码支持 0、1 和 2 三种模式 */
	CompressionMode uint16
	/** 球谐函数的最小值，默认为 -1.5 */
	MinHarmonicsValue float32
	/** 球谐函数的最大值，默认为 1.5 */
	MaxHarmonicsValue float32

	ShDegree int
}

func (h *KsplatHeader) ToString() string {
	return fmt.Sprintf("KsplatVersion     : %v.%v\nSectionCount      : %v\nSplatCount        : %v\nCompressionMode   : %v\nMinHarmonicsValue : %s\nMaxHarmonicsValue : %s\nShDegree          : %v\n",
		h.MajorVersion, h.MinorVersion, h.SectionCount, h.SplatCount, h.CompressionMode, cmn.FormatFloat32(h.MinHarmonicsValue), cmn.FormatFloat32(h.MaxHarmonicsValue), h.ShDegree)
}

/** 分段头信息。每个分段头信息长度为 1024 字节，有未使用的空间可能是为未来的扩展预留的 */
type SectionHeader struct {
	/** 实际点（splat）数量 */
	SectionSplatCount uint32
	/** 能容纳的最大点数量 */
	SectionSplatCapacity uint32
	/** 每个桶的容量 */
	BucketCapacity uint32
	/** 桶的数量 */
	BucketCount uint32
	/** 空间块的大小 */
	BlockSize float32
	/** 桶存储的大小 */
	BucketSize uint16
	/** 量化范围，如果未指定则使用主头中的值 */
	QuantizationRange uint32
	/** 满桶的数量 */
	FullBucketCount uint32
	/** 未满桶的数量 */
	PartialBucketCount uint32
	/** 球谐函数的级别 */
	ShDegree uint16
}

func (h *SectionHeader) ToString() string {
	return fmt.Sprintf("SectionSplatCount    : %v\nSectionSplatCapacity : %v\nBucketCapacity       : %v\nBucketCount          : %v\nBlockSize            : %s\nBucketSize           : %v\nQuantizationRange    : %v\nFullBucketCount      : %v\nPartialBucketCount   : %v\nShDegree             : %v\n",
		h.SectionSplatCount, h.SectionSplatCapacity, h.BucketCapacity, h.BucketCount, cmn.FormatFloat32(h.BlockSize), h.BucketSize, h.QuantizationRange, h.FullBucketCount, h.PartialBucketCount, h.ShDegree)
}

type CompressionConfig struct {
	centerBytes        int
	scaleBytes         int
	rotationBytes      int
	colorBytes         int
	harmonicsBytes     int
	scaleStartByte     int
	rotationStartByte  int
	colorStartByte     int
	harmonicsStartByte int
	scaleQuantRange    uint32
}

func ReadKsplat(ksplatFile string, readHeadOnly bool) (*SectionHeader, *KsplatHeader, []*SplatData) {

	file, err := os.Open(ksplatFile)
	cmn.ExitOnError(err)
	defer file.Close()

	// 文件头读取
	HeadSize := 4096
	SectionSize := 1024
	bts := make([]byte, HeadSize)
	size, err := file.Read(bts)
	cmn.ExitOnError(err)
	if size != HeadSize {
		cmn.ExitOnError(errors.New("[KSPLAT ERROR] invalid ksplat format")) // 还不够文件头大小
	}

	mainHeader := &KsplatHeader{
		MajorVersion:      bts[0],
		MinorVersion:      bts[1],
		SectionCount:      int(cmn.BytesToInt32(bts[4:8])),
		SplatCount:        int(cmn.BytesToInt32(bts[16:20])),
		CompressionMode:   cmn.BytesToUint16(bts[20:22]),
		MinHarmonicsValue: cmn.BytesToFloat32(bts[36:40]),
		MaxHarmonicsValue: cmn.BytesToFloat32(bts[40:44]),
	}
	if mainHeader.MajorVersion != 0 || mainHeader.MinorVersion != 1 {
		cmn.ExitOnError(errors.New("[KSPLAT ERROR] unsupported version"))
	}
	if mainHeader.CompressionMode > 2 {
		cmn.ExitOnError(errors.New("[KSPLAT ERROR] invalid compression mode"))
	}
	if mainHeader.SplatCount == 0 {
		cmn.ExitOnError(errors.New("[KSPLAT ERROR] data empty"))
	}
	if mainHeader.MinHarmonicsValue == 0 {
		mainHeader.MinHarmonicsValue = -1.5
	}
	if mainHeader.MaxHarmonicsValue == 0 {
		mainHeader.MaxHarmonicsValue = 1.5
	}

	// 按压缩模式对应
	cc := &CompressionConfig{}
	switch mainHeader.CompressionMode {
	case 0:
		cc.centerBytes = 12
		cc.scaleBytes = 12
		cc.rotationBytes = 16
		cc.colorBytes = 4
		cc.harmonicsBytes = 4
		cc.scaleStartByte = 12
		cc.rotationStartByte = 24
		cc.colorStartByte = 40
		cc.harmonicsStartByte = 44
		cc.scaleQuantRange = 1
	case 1:
		cc.centerBytes = 6
		cc.scaleBytes = 6
		cc.rotationBytes = 8
		cc.colorBytes = 4
		cc.harmonicsBytes = 2
		cc.scaleStartByte = 6
		cc.rotationStartByte = 12
		cc.colorStartByte = 20
		cc.harmonicsStartByte = 24
		cc.scaleQuantRange = 32767
	case 2:
		cc.centerBytes = 6
		cc.scaleBytes = 6
		cc.rotationBytes = 8
		cc.colorBytes = 4
		cc.harmonicsBytes = 1
		cc.scaleStartByte = 6
		cc.rotationStartByte = 12
		cc.colorStartByte = 20
		cc.harmonicsStartByte = 24
		cc.scaleQuantRange = 32767
	}

	// 分段头读取
	secHeaders := make([]*SectionHeader, mainHeader.SectionCount)
	for i := 0; i < mainHeader.SectionCount; i++ {
		bs := make([]byte, SectionSize)
		_, err := file.ReadAt(bs, int64(HeadSize+i*SectionSize))
		cmn.ExitOnError(err)

		h := &SectionHeader{
			SectionSplatCount:    cmn.BytesToUint32(bs[0:4]),
			SectionSplatCapacity: cmn.BytesToUint32(bs[4:8]),
			BucketCapacity:       cmn.BytesToUint32(bs[8:12]),
			BucketCount:          cmn.BytesToUint32(bs[12:16]),
			BlockSize:            cmn.BytesToFloat32(bs[16:20]),
			BucketSize:           cmn.BytesToUint16(bs[20:22]),
			QuantizationRange:    cmn.BytesToUint32(bs[24:28]),
			FullBucketCount:      cmn.BytesToUint32(bs[32:36]),
			PartialBucketCount:   cmn.BytesToUint32(bs[36:40]),
			ShDegree:             cmn.BytesToUint16(bs[40:42]),
		}
		if h.QuantizationRange == 0 {
			h.QuantizationRange = cc.scaleQuantRange
		}
		secHeaders[i] = h

		// 球谐系数通过遍历分段头取最大值求得
		if mainHeader.ShDegree < int(h.ShDegree) {
			mainHeader.ShDegree = int(h.ShDegree)
		}
	}

	// 仅需读取概要信息时，直接返回头部信息跳过内容读取
	if readHeadOnly {
		return secHeaders[0], mainHeader, nil
	}

	shDims := []int{0, 9, 24, 15}
	shComponents := shDims[mainHeader.ShDegree]
	offset := HeadSize + mainHeader.SectionCount*SectionSize
	datas := make([]*SplatData, mainHeader.SplatCount)
	n := 0
	for i := 0; i < mainHeader.SectionCount; i++ {
		secHead := secHeaders[i]

		bytesPerSplat := cc.centerBytes + cc.scaleBytes + cc.rotationBytes + cc.colorBytes + cc.harmonicsBytes*shComponents
		positionScaleFactor := float64(secHead.BlockSize) / 2.0 / float64(secHead.QuantizationRange)

		// 部分桶元数据
		partialBucketMetaSize := int(secHead.PartialBucketCount * 4) // 各未满桶中的点数量
		partialBucketSizes := make([]byte, partialBucketMetaSize)
		_, err := file.ReadAt(partialBucketSizes, int64(offset))
		cmn.ExitOnError(err)
		offset += partialBucketMetaSize
		// 桶中心数据
		bucketCentersSize := int(secHead.BucketCount * 3 * 4) // 每个桶相应有xyz坐标基数
		bucketCenters := make([]byte, bucketCentersSize)
		_, err = file.ReadAt(bucketCenters, int64(offset))
		cmn.ExitOnError(err)
		offset += bucketCentersSize
		// 点数据
		sectionDataSize := bytesPerSplat * int(secHead.SectionSplatCapacity) // 按分段数据最大容量读取
		splatData := make([]byte, sectionDataSize)
		_, err = file.ReadAt(splatData, int64(offset))
		cmn.ExitOnError(err)
		offset += sectionDataSize

		fullBucketSplats := int(secHead.FullBucketCount * secHead.BucketCapacity)
		currentPartialBucket := int(secHead.FullBucketCount)
		currentPartialBase := fullBucketSplats
		sectionSplatCount := int(secHead.SectionSplatCount)
		for j := range sectionSplatCount {
			var bucketIdx int
			if secHead.BucketCapacity > 0 {
				if j < fullBucketSplats {
					bucketIdx = int(math.Floor(float64(j) / float64(secHead.BucketCapacity)))
				} else {
					partialIdx := currentPartialBucket - int(secHead.FullBucketCount)                              // 未满桶范围的下标
					currentBucketSize := int(cmn.BytesToUint32(partialBucketSizes[partialIdx*4 : partialIdx*4+4])) // 未满桶中的数量
					if j >= currentPartialBase+currentBucketSize {
						// 当前未满桶已读取完，计算准备下一个未满桶
						currentPartialBucket++                  // 下一个未满桶
						currentPartialBase += currentBucketSize // 总数累加
					}
					bucketIdx = currentPartialBucket
				}
			}

			data := &SplatData{}

			// Decode position
			splatByteOffset := j * bytesPerSplat
			if mainHeader.CompressionMode == 0 {
				data.PositionX = cmn.BytesToFloat32(splatData[splatByteOffset : splatByteOffset+4])
				data.PositionY = cmn.BytesToFloat32(splatData[splatByteOffset+4 : splatByteOffset+8])
				data.PositionZ = cmn.BytesToFloat32(splatData[splatByteOffset+8 : splatByteOffset+12])
			} else {
				data.PositionX = cmn.ClipFloat32((float64(cmn.BytesToUint16(splatData[splatByteOffset:splatByteOffset+2]))-float64(secHead.QuantizationRange))*positionScaleFactor + float64(cmn.BytesToFloat32(bucketCenters[bucketIdx*3*4:bucketIdx*3*4+4])))
				data.PositionY = cmn.ClipFloat32((float64(cmn.BytesToUint16(splatData[splatByteOffset+2:splatByteOffset+4]))-float64(secHead.QuantizationRange))*positionScaleFactor + float64(cmn.BytesToFloat32(bucketCenters[bucketIdx*3*4+4:bucketIdx*3*4+8])))
				data.PositionZ = cmn.ClipFloat32((float64(cmn.BytesToUint16(splatData[splatByteOffset+4:splatByteOffset+6]))-float64(secHead.QuantizationRange))*positionScaleFactor + float64(cmn.BytesToFloat32(bucketCenters[bucketIdx*3*4+8:bucketIdx*3*4+12])))
			}

			// Decode scales
			if mainHeader.CompressionMode == 0 {
				data.ScaleX = cmn.DecodeSplatScale(cmn.BytesToFloat32(splatData[splatByteOffset+cc.scaleStartByte : splatByteOffset+cc.scaleStartByte+4]))
				data.ScaleY = cmn.DecodeSplatScale(cmn.BytesToFloat32(splatData[splatByteOffset+cc.scaleStartByte+4 : splatByteOffset+cc.scaleStartByte+8]))
				data.ScaleZ = cmn.DecodeSplatScale(cmn.BytesToFloat32(splatData[splatByteOffset+cc.scaleStartByte+8 : splatByteOffset+cc.scaleStartByte+12]))
			} else {
				data.ScaleX = cmn.DecodeSplatScale(cmn.DecodeFloat16(cmn.BytesToUint16(splatData[splatByteOffset+cc.scaleStartByte : splatByteOffset+cc.scaleStartByte+2])))
				data.ScaleY = cmn.DecodeSplatScale(cmn.DecodeFloat16(cmn.BytesToUint16(splatData[splatByteOffset+cc.scaleStartByte+2 : splatByteOffset+cc.scaleStartByte+4])))
				data.ScaleZ = cmn.DecodeSplatScale(cmn.DecodeFloat16(cmn.BytesToUint16(splatData[splatByteOffset+cc.scaleStartByte+4 : splatByteOffset+cc.scaleStartByte+6])))
			}

			// Decode rotation quaternion
			var rot0 float32
			var rot1 float32
			var rot2 float32
			var rot3 float32
			if mainHeader.CompressionMode == 0 {
				rot0 = cmn.BytesToFloat32(splatData[splatByteOffset+cc.rotationStartByte : splatByteOffset+cc.rotationStartByte+4])
				rot1 = cmn.BytesToFloat32(splatData[splatByteOffset+cc.rotationStartByte+4 : splatByteOffset+cc.rotationStartByte+8])
				rot2 = cmn.BytesToFloat32(splatData[splatByteOffset+cc.rotationStartByte+8 : splatByteOffset+cc.rotationStartByte+12])
				rot3 = cmn.BytesToFloat32(splatData[splatByteOffset+cc.rotationStartByte+12 : splatByteOffset+cc.rotationStartByte+16])
			} else {
				rot0 = cmn.DecodeFloat16(cmn.BytesToUint16(splatData[splatByteOffset+cc.rotationStartByte : splatByteOffset+cc.rotationStartByte+2]))
				rot1 = cmn.DecodeFloat16(cmn.BytesToUint16(splatData[splatByteOffset+cc.rotationStartByte+2 : splatByteOffset+cc.rotationStartByte+4]))
				rot2 = cmn.DecodeFloat16(cmn.BytesToUint16(splatData[splatByteOffset+cc.rotationStartByte+4 : splatByteOffset+cc.rotationStartByte+6]))
				rot3 = cmn.DecodeFloat16(cmn.BytesToUint16(splatData[splatByteOffset+cc.rotationStartByte+6 : splatByteOffset+cc.rotationStartByte+8]))
			}
			data.RotationW = cmn.EncodeSplatRotation(float64(rot0))
			data.RotationX = cmn.EncodeSplatRotation(float64(rot1))
			data.RotationY = cmn.EncodeSplatRotation(float64(rot2))
			data.RotationZ = cmn.EncodeSplatRotation(float64(rot3))

			// Decode color and opacity
			data.ColorR = uint8(splatData[splatByteOffset+cc.colorStartByte])
			data.ColorG = uint8(splatData[splatByteOffset+cc.colorStartByte+1])
			data.ColorB = uint8(splatData[splatByteOffset+cc.colorStartByte+2])
			data.ColorA = uint8(splatData[splatByteOffset+cc.colorStartByte+3])

			// Decode Harmonic
			if mainHeader.ShDegree > 0 {
				shIndexs := []int{0, 3, 6, 1, 4, 7, 2, 5, 8,
					9, 14, 19, 10, 15, 20, 11, 16, 21, 12, 17, 22, 13, 18, 23,
					24, 31, 38, 25, 32, 39, 26, 33, 40, 27, 34, 41, 28, 35, 42, 29, 36, 43, 30, 37, 44}
				shDims := []int{0, 3, 8, 15}

				shCnt := shDims[mainHeader.ShDegree] * 3
				shs := make([]byte, 45)
				cnt := 0
				for k := range shCnt {
					switch mainHeader.CompressionMode {
					case 0:
						shOffset := splatByteOffset + cc.harmonicsStartByte + shIndexs[k]*4
						shs[cnt] = cmn.EncodeSplatSH(float64(cmn.BytesToFloat32(splatData[shOffset : shOffset+4])))
					case 1:
						shOffset := splatByteOffset + cc.harmonicsStartByte + shIndexs[k]*2
						shs[cnt] = cmn.EncodeSplatSH(float64(cmn.DecodeFloat16(cmn.BytesToUint16(splatData[shOffset : shOffset+2]))))
					case 2:
						shOffset := splatByteOffset + cc.harmonicsStartByte + shIndexs[k]
						shs[cnt] = cmn.EncodeSplatSH(float64(mainHeader.MinHarmonicsValue + (float32(uint8(splatData[shOffset]))/255.0)*(mainHeader.MaxHarmonicsValue-mainHeader.MinHarmonicsValue)))
					}
					cnt++
				}

				for ; cnt < 45; cnt++ {
					shs[cnt] = cmn.EncodeSplatSH(0)
				}
				switch mainHeader.ShDegree {
				case 3:
					data.SH2 = shs[:24]
					data.SH3 = shs[24:]
				case 2:
					data.SH2 = shs[:24]
				case 1:
					data.SH1 = shs[:9]
				}
			}

			// 存放数据
			datas[n] = data
			n++
		}
	}

	return secHeaders[0], mainHeader, datas
}
