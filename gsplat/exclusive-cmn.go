package gsplat

import "gsbox/cmn"

const CreaterIdOpen uint32 = 1202056903  // 创建者ID，开放版
const ExclusiveIdOpen uint32 = 0         // 专属ID，开放版
const CreaterIdPrd uint32 = 0            // 创建者ID，官方
const ExclusiveIdPrd uint32 = 3141592653 // 专属ID，专业版

// 数据处理
func ProcessDatas(datas []*SplatData) []*SplatData {
	TransformDatas(datas)
	rsDatas := FilterDatas(datas)
	Sort(rsDatas)
	return rsDatas
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

func IsPrdExclusiveId(id uint32) bool {
	return id == ExclusiveIdPrd
}

func GetSpxOutputHeaderHash(bts []byte) uint32 {
	return cmn.HashBytes(bts)
}

func CheckSpxHeaderHash(bts []byte, hash uint32) bool {
	return cmn.HashBytes(bts) == hash
}
