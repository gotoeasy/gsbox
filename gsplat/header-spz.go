package gsplat

import (
	"errors"
	"fmt"
	"gsbox/cmn"
)

const HeaderSizeSpz = 16
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
	/** Reserved fields */
	FractionalBits uint8
	/** Reserved fields */
	Flags uint8
	/** 0 */
	Reserved uint8
}

func (h *SpzHeader) ToBytes() []byte {
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

func ParseSpzHeader(bts []byte) *SpzHeader {
	header := &SpzHeader{
		Magic:          cmn.BytesToUint32(bts[0:4]),
		Version:        cmn.BytesToUint32(bts[4:8]),
		NumPoints:      cmn.BytesToUint32(bts[8:12]),
		ShDegree:       bts[12],
		FractionalBits: bts[13],
		Flags:          bts[14],
		Reserved:       bts[15],
	}

	if header.Magic != SPZ_MAGIC {
		cmn.ExitOnError(errors.New("[SPZ ERROR] deserializePackedGaussians: header not found"))
	}
	if header.Version != 2 {
		cmn.ExitOnError(errors.New("[SPZ ERROR] deserializePackedGaussians: version not supported: " + cmn.Uint32ToString(header.Version)))
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
	return fmt.Sprintf("3DGS model format spz\nMagic          : %v\nVersion        : %v\nNumPoints      : %v\nShDegree       : %v\nFractionalBits : %v\nFlags          : %v\nReserved       : %v",
		h.Magic, h.Version, h.NumPoints, h.ShDegree, h.FractionalBits, h.Flags, h.Reserved)
}
