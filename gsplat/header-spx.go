package gsplat

import (
	"errors"
	"fmt"
	"gsbox/cmn"
	"os"
	"strings"
)

/** spx文件头长度 */
const HeaderSizeSpx = 128

/** Spx format file header */
type SpxHeader struct {
	/** Fixed string, fixed to spx */
	Fixed string
	/** Spx version number, fixed to 1 */
	Version uint8
	/** Number of Gaussian primitives, must be specified */
	SplatCount int32
	/** Model bounding box vertices */
	MinX float32
	/** Model bounding box vertices */
	MaxX float32
	/** Model bounding box vertices */
	MinY float32
	/** Model bounding box vertices */
	MaxY float32
	/** Model bounding box vertices */
	MinZ float32
	/** Model bounding box vertices */
	MaxZ float32
	/** Min Center height */
	MinTopY float32
	/** Max Center height */
	MaxTopY float32
	/** Creation date (YYYYMMDD) */
	CreateDate uint32
	/** Creator identifier, 0 represents official tools */
	CreaterId uint32
	/** Exclusive format identifier, 0 represents official open format, can be customized, must be specified */
	ExclusiveId uint32
	/** spherical harmonics degree(1/2/3, others mean 0) */
	ShDegree uint8
	/** gaussian splat data type (default 0) */
	Flag1 uint8 // 废弃
	Flag2 uint8 // 废弃
	Flag3 uint8 // 废弃
	/** Reserved fields */
	Reserve1 uint32
	/** Reserved fields */
	Reserve2 uint32
	/** Comments (only supports ASCII characters) */
	Comment string
	/** Hash */
	Hash uint32

	/** v2 */
	Flags uint8

	/** v3 */
	Lod uint8
	/** v3 */
	Reserve3 uint8
	/** v3 */
	Palettes []uint8

	checkHash bool
}

func (h *SpxHeader) IsValid() bool {
	return h.checkHash
}

/** 是否倒立的模型(v2) */
func (h *SpxHeader) IsInverted() bool {
	return (h.Flags & 0b10000000) > 0
}

/** 是否优化的大场景模型(v2) */
func (h *SpxHeader) IsLargeScene() bool {
	return (h.Flags & 0b1) > 0
}

func (h *SpxHeader) ToBytes() []byte {
	bts := make([]byte, 0)
	bts = append(bts, h.Fixed...)
	bts = append(bts, h.Version)
	bts = append(bts, cmn.Int32ToBytes(h.SplatCount)...)
	bts = append(bts, cmn.Float32ToBytes(h.MinX)...)
	bts = append(bts, cmn.Float32ToBytes(h.MaxX)...)
	bts = append(bts, cmn.Float32ToBytes(h.MinY)...)
	bts = append(bts, cmn.Float32ToBytes(h.MaxY)...)
	bts = append(bts, cmn.Float32ToBytes(h.MinZ)...)
	bts = append(bts, cmn.Float32ToBytes(h.MaxZ)...)
	bts = append(bts, cmn.Float32ToBytes(h.MinTopY)...)
	bts = append(bts, cmn.Float32ToBytes(h.MaxTopY)...)
	bts = append(bts, cmn.Uint32ToBytes(h.CreateDate)...)
	bts = append(bts, cmn.Uint32ToBytes(h.CreaterId)...)
	bts = append(bts, cmn.Uint32ToBytes(h.ExclusiveId)...)
	bts = append(bts, h.ShDegree)
	if h.Version == 1 {
		// v1
		bts = append(bts, h.Flag1)
		bts = append(bts, h.Flag2)
		bts = append(bts, h.Flag3)
	} else {
		// v2
		bts = append(bts, h.Flags)
		bts = append(bts, h.Lod, h.Reserve3)
	}
	bts = append(bts, cmn.Uint32ToBytes(h.Reserve1)...)
	bts = append(bts, cmn.Uint32ToBytes(h.Reserve2)...)
	bts = append(bts, cmn.StringToBytes(cmn.Left(cmn.Trim(h.Comment)+strings.Repeat(" ", 60), 60))...) // 右边补足空格取60个accsii字符
	bts = append(bts, cmn.Uint32ToBytes(GetSpxOutputHeaderHash(bts[0:124]))...)
	return bts
}

func ParseSpxHeader(spxFile string) *SpxHeader {
	file, err := os.Open(spxFile)
	cmn.ExitOnError(err)
	defer file.Close()

	bs := make([]byte, HeaderSizeSpx)
	_, err = file.Read(bs)
	if err != nil {
		cmn.ExitOnError(err)
	}

	spxVer := int(bs[3])
	if bs[0] == 's' && bs[1] == 'p' && bs[2] == 'x' && (spxVer > 0 && spxVer <= NewestSpxVersion) {
		header := readSpxHeader(bs)
		if !CanParseExclusiveId(header.ExclusiveId) && !Args.HasCmd("info") {
			cmn.ExitOnError(errors.New("unknown exclusive id: " + cmn.Uint32ToString(header.ExclusiveId))) // 内含不识别的专属格式时，退出
		}
		return header
	}

	cmn.ExitOnError(errors.New("unknown format: " + spxFile))
	return nil
}

func readSpxHeader(bts []byte) *SpxHeader {
	header := &SpxHeader{
		Fixed:       "spx",
		Version:     bts[3],
		SplatCount:  cmn.BytesToInt32(bts[4:8]),
		MinX:        cmn.BytesToFloat32(bts[8:12]),
		MaxX:        cmn.BytesToFloat32(bts[12:16]),
		MinY:        cmn.BytesToFloat32(bts[16:20]),
		MaxY:        cmn.BytesToFloat32(bts[20:24]),
		MinZ:        cmn.BytesToFloat32(bts[24:28]),
		MaxZ:        cmn.BytesToFloat32(bts[28:32]),
		MinTopY:     cmn.BytesToFloat32(bts[32:36]),
		MaxTopY:     cmn.BytesToFloat32(bts[36:40]),
		CreateDate:  cmn.BytesToUint32(bts[40:44]),
		CreaterId:   cmn.BytesToUint32(bts[44:48]),
		ExclusiveId: cmn.BytesToUint32(bts[48:52]),
		ShDegree:    bts[52],
		Flag1:       bts[53], // v1
		Flag2:       bts[54], // v1
		Flag3:       bts[55], // v1
		Flags:       bts[53], // v2
		Lod:         bts[54], // v3
		Reserve3:    bts[55], // v3
		Reserve1:    cmn.BytesToUint32(bts[56:60]),
		Reserve2:    cmn.BytesToUint32(bts[60:64]),
		Hash:        cmn.BytesToUint32(bts[124:128]),
		checkHash:   CheckSpxHeaderHash(bts[0:124], cmn.BytesToUint32(bts[124:128])),
	}
	if header.ShDegree != 1 && header.ShDegree != 2 && header.ShDegree != 3 {
		header.ShDegree = 0
	}

	comment := ""
	for i := 64; i < 124; i++ {
		comment += cmn.BytesToString(bts[i : i+1])
	}
	header.Comment = cmn.Trim(comment)

	return header
}

func (h *SpxHeader) ToStringSpx() string {
	switch h.Version {
	case 1:
		// v1
		return fmt.Sprintf("3DGS model format spx\nSpx version  : 1\nSplatCount   : %v\nMinX, MaxX   : %v, %v\nMinY, MaxY   : %v, %v\nMinZ, MaxZ   : %v, %v\nMinTopY      : %v\nMaxTopY      : %v\nCreateDate   : %v\nCreaterId    : %v\nExclusiveId  : %v\nShDegree     : %v\nFlag1        : %v\nFlag2        : %v\nFlag3        : %v\nComment      : %v\nHash         : %v (%v)",
			h.SplatCount, h.MinX, h.MaxX, h.MinY, h.MaxY, h.MinZ, h.MaxZ, h.MinTopY, h.MaxTopY, h.CreateDate, h.CreaterId, h.ExclusiveId, h.ShDegree, h.Flag1, h.Flag2, h.Flag3, h.Comment, h.Hash, h.checkHash)
	case 2:
		// v2
		return fmt.Sprintf("3DGS model format spx\nSpx version  : 2\nSplatCount   : %v\nMinX, MaxX   : %v, %v\nMinY, MaxY   : %v, %v\nMinZ, MaxZ   : %v, %v\nMinTopY      : %v\nMaxTopY      : %v\nCreateDate   : %v\nCreaterId    : %v\nExclusiveId  : %v\nShDegree     : %v\nIsInverted   : %v\nIsLargeScene : %v\nComment      : %v\nHash         : %v (%v)",
			h.SplatCount, h.MinX, h.MaxX, h.MinY, h.MaxY, h.MinZ, h.MaxZ, h.MinTopY, h.MaxTopY, h.CreateDate, h.CreaterId, h.ExclusiveId, h.ShDegree, h.IsInverted(), h.IsLargeScene(), h.Comment, h.Hash, h.checkHash)
	case 3:
		// v3
		return fmt.Sprintf("3DGS model format spx\nSpx version  : %v\nSplatCount   : %v\nMinX, MaxX   : %v, %v\nMinY, MaxY   : %v, %v\nMinZ, MaxZ   : %v, %v\nMinTopY      : %v\nMaxTopY      : %v\nCreateDate   : %v\nCreaterId    : %v\nExclusiveId  : %v\nShDegree     : %v\nIsInverted   : %v\nIsLargeScene : %v\nLOD          : %v\nComment      : %v\nHash         : %v (%v)",
			h.Version, h.SplatCount, h.MinX, h.MaxX, h.MinY, h.MaxY, h.MinZ, h.MaxZ, h.MinTopY, h.MaxTopY, h.CreateDate, h.CreaterId, h.ExclusiveId, h.ShDegree, h.IsInverted(), h.IsLargeScene(), h.Lod, h.Comment, h.Hash, h.checkHash)
	default:
		// v3
		return fmt.Sprintf("3DGS model format spx\nSpx version  : %v\nSplatCount   : %v\nMinX, MaxX   : %v, %v\nMinY, MaxY   : %v, %v\nMinZ, MaxZ   : %v, %v\nMinTopY      : %v\nMaxTopY      : %v\nCreateDate   : %v\nCreaterId    : %v\nExclusiveId  : %v\nShDegree     : %v\nIsInverted   : %v\nIsLargeScene : %v\nLOD          : %v\nComment      : %v\nHash         : %v (%v)",
			h.Version, h.SplatCount, h.MinX, h.MaxX, h.MinY, h.MaxY, h.MinZ, h.MaxZ, h.MinTopY, h.MaxTopY, h.CreateDate, h.CreaterId, h.ExclusiveId, h.ShDegree, h.IsInverted(), h.IsLargeScene(), h.Lod, h.Comment, h.Hash, h.checkHash)
	}

}
