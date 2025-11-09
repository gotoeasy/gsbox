package gsplat

import "gsbox/cmn"

const CreaterIdOpen uint32 = 1202056903  // 创建者ID，开放版
const ExclusiveIdOpen uint32 = 0         // 专属ID，开放版
const CreaterIdPrd uint32 = 0            // 创建者ID，官方
const ExclusiveIdPrd uint32 = 3141592653 // 专属ID，专业版

// 数据处理
func ProcessDatas(datas []*SplatData) []*SplatData {
	datas = FilterDatas(datas)
	datas = TransformDatas(datas)
	Sort(datas)
	return datas
}

func GetOutputCreaterId() uint32 {
	return CreaterIdOpen
}

func GetOutputExclusiveId() uint32 {
	return ExclusiveIdOpen
}

func IsOpenCreaterId(id uint32) bool {
	return id == CreaterIdOpen
}

func IsOpenExclusiveId(id uint32) bool {
	return id == ExclusiveIdOpen
}

func CanParseExclusiveId(id uint32) bool {
	return id == ExclusiveIdOpen
}

func IsPrdExclusiveId(id uint32) bool {
	return id == ExclusiveIdPrd
}

func GetSpxOutputHeaderHash(bts []byte) uint32 {
	return cmn.HashBytes(bts)
}

func CheckSpxHeaderHash(bts []byte, hash uint32) bool {
	return cmn.HashBytes(bts) == hash
}

func GetLodFlag() uint8 {
	return 0
}

func WriteSpxV2(spxFile string, rows []*SplatData, comment string, shDegree uint8) {
	WriteSpxOpenV2(spxFile, rows, comment, shDegree)
}

func ReadSpxV2(spxFile string, header *SpxHeader) (*SpxHeader, []*SplatData) {
	return ReadSpxOpenV2(spxFile, header)
}

func DefaultSpxComment() string {
	return "created by gsbox " + cmn.VER + " https://github.com/gotoeasy/gsbox"
}

func BlockFormatDesc(bf int) string {
	rs := ""
	switch bf {
	case BF_SPLAT19:
		rs = "(splat per 19 bytes)"
	case BF_SPLAT10019:
		rs = "(splat per 19 bytes)"
	case BF_SPLAT20:
		rs = "(splat per 20 bytes)"
	case BF_SPLAT190_WEBP:
		rs = "(splat per 19 bytes, webp encoding)"
	case BF_SPLAT10190_WEBP:
		rs = "(splat per 19 bytes, webp encoding)"
	case BF_SH_PALETTES:
		rs = "(compress only)"
	case BF_SH_PALETTES_WEBP:
		rs = "(webp encoding)"
	}

	return rs
}
