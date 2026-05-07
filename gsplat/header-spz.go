package gsplat

import (
	"errors"
	"fmt"
	"gsbox/cmn"
)

const HeaderSizeSpzV3 = 16
const HeaderSizeSpzV4 = 32
const SPZ_MAGIC = 0x5053474e // NGSP = Niantic gaussian splat

type SpzHeader struct {
	/** 1347635022 */
	Magic uint32
	/** 2 */
	Version uint32
	/** Number of Gaussian primitives, must be specified */
	NumPoints uint32
	/** 0,1,2,3 */
	ShDegree uint8
	/** 位置定点数的小数位数 */
	FractionalBits uint8
	/** bit0: 是否抗锯齿, bit1: 是否有扩展【V4】 */
	Flags uint8
	/** 属性流个数（通常 6）【V4】 */
	NumStreams uint8
	/** TableOfContent 相对文件起始位置的字节偏移量【V4】 */
	TocByteOffset uint32
	/** 0 */
	Reserved uint8
}

func (h *SpzHeader) ToBytes() []byte {
	if h.Version < 4 {
		bts := make([]byte, 0)
		bts = append(bts, cmn.Uint32ToBytes(h.Magic)...)
		bts = append(bts, cmn.Uint32ToBytes(h.Version)...)
		bts = append(bts, cmn.Uint32ToBytes(h.NumPoints)...)
		bts = append(bts, h.ShDegree)
		bts = append(bts, h.FractionalBits)
		bts = append(bts, h.Flags)
		bts = append(bts, h.Reserved)
		return bts
	}

	// v4
	bts := make([]byte, 0)
	bts = append(bts, cmn.Uint32ToBytes(h.Magic)...)
	bts = append(bts, cmn.Uint32ToBytes(h.Version)...)
	bts = append(bts, cmn.Uint32ToBytes(h.NumPoints)...)
	bts = append(bts, h.ShDegree)
	bts = append(bts, h.FractionalBits)
	bts = append(bts, h.Flags)
	bts = append(bts, h.NumStreams)
	bts = append(bts, cmn.Uint32ToBytes(h.TocByteOffset)...)
	for range 12 {
		bts = append(bts, 0)
	}
	return bts
}

func IsSpzV4(bts []byte) bool {
	magic := cmn.BytesToUint32(bts[0:4])
	if magic != SPZ_MAGIC {
		return false
	}
	version := cmn.BytesToUint32(bts[4:8])
	return version == 4
}

func ParseSpzHeader(bts []byte) *SpzHeader {

	magic := cmn.BytesToUint32(bts[0:4])
	if magic != SPZ_MAGIC {
		cmn.ExitOnError(errors.New("[SPZ ERROR] deserializePackedGaussians: header not found"))
	}

	header := &SpzHeader{}
	version := cmn.BytesToUint32(bts[4:8])
	if version == 2 || version == 3 {
		// v2,v3
		header = &SpzHeader{
			Magic:          magic,
			Version:        version,
			NumPoints:      cmn.BytesToUint32(bts[8:12]),
			ShDegree:       bts[12],
			FractionalBits: bts[13],
			Flags:          bts[14],
			Reserved:       bts[15],
		}
	} else {
		// v4
		header = &SpzHeader{
			Magic:          magic,
			Version:        version,
			NumPoints:      cmn.BytesToUint32(bts[8:12]),
			ShDegree:       bts[12],
			FractionalBits: bts[13],
			Flags:          bts[14],
			NumStreams:     bts[15],
			TocByteOffset:  cmn.BytesToUint32(bts[16:20]),
			Reserved:       bts[20],
		}
	}

	if header.ShDegree > 3 {
		cmn.ExitOnError(errors.New("[SPZ ERROR] deserializePackedGaussians: Unsupported SH degree: " + cmn.IntToString(int(header.ShDegree))))
	}
	if header.FractionalBits != 12 {
		// 仅支持这一种编码方式（坐标24位整数编码）
		cmn.ExitOnError(errors.New("[SPZ ERROR] deserializePackedGaussians: Unsupported FractionalBits: " + cmn.IntToString(int(header.FractionalBits))))
	}

	return header
}

func (h *SpzHeader) ToString() string {
	if h.Version < 4 {
		return fmt.Sprintf("3DGS model format spz\nMagic          : %v\nVersion        : %v\nNumPoints      : %v\nShDegree       : %v\nFractionalBits : %v\nFlags          : %v\nReserved       : %v",
			h.Magic, h.Version, h.NumPoints, h.ShDegree, h.FractionalBits, h.Flags, h.Reserved)
	} else {
		return fmt.Sprintf("3DGS model format spz\nMagic          : %v\nVersion        : %v\nNumPoints      : %v\nShDegree       : %v\nFractionalBits : %v\nFlags          : %v\nNumStreams     : %v\nTocByteOffset  : %v",
			h.Magic, h.Version, h.NumPoints, h.ShDegree, h.FractionalBits, h.Flags, h.NumStreams, h.TocByteOffset)
	}
}
