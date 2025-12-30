package gsplat

import (
	"fmt"
	"gsbox/cmn"
	"log"
	"time"
)

const CreaterIdOpen uint32 = 1202056903 // 创建者ID，开放版
const ExclusiveIdOpen uint32 = 0        // 专属ID，开放版
const CreaterIdPrd uint32 = 0           // 创建者ID，官方

const printProgressLog = false
const PhaseJoin = -1
const PhaseRead = 0
const PhaseProc = 1
const PhaseWrite = 2
const PhaseKmean = 3

var Phases = []int{100, 100, 100, 100} // 0:读, 1:数据处理, 2:写, 3:聚类
var Progress = []int{0, 0, 0, 0}

func OnProgress(phase int, current int, total int) {
	if printProgressLog {
		if phase > 0 {
			Phases[phase] = total
			Progress[phase] = current
		} else if phase == 0 {
			// 单文件读
			if !oArg.isJoin {
				Phases[phase] = total
				Progress[phase] = current
			}
		} else {
			// 合并读
			Phases[PhaseRead] = total
			Progress[PhaseRead] = current
		}
	}
}

func init() {
	if !printProgressLog {
		return
	}

	go (func() {
		lastProgress := ""
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			var phasesLen int
			var totalProgress float64
			for i, v := range Progress {
				if Phases[i] > 0 {
					totalProgress += float64(v) / float64(Phases[i])
					phasesLen++
				}
			}
			totalPercent := totalProgress / float64(phasesLen)
			theProgress := fmt.Sprintf("%.1f%%", totalPercent*100)
			if lastProgress != theProgress {
				lastProgress = theProgress
				log.Println("[Info] progress:", lastProgress)
			}

			if totalPercent >= 1.0 {
				ticker.Stop()
				break
			}
		}
	})()
}

// 数据处理
func ProcessDatas(datas []*SplatData) []*SplatData {
	datas = FilterDatas(datas)
	OnProgress(PhaseProc, 30, 100)
	datas = TransformDatas(datas)
	OnProgress(PhaseProc, 80, 100)
	Sort(datas)
	OnProgress(PhaseProc, 100, 100)
	return datas
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
	return false
}

func GetSpxOutputHeaderHash(bts []byte) uint32 {
	return cmn.HashBytes(bts)
}

func CheckSpxHeaderHash(bts []byte, hash uint32) bool {
	return cmn.HashBytes(bts) == hash
}

func CheckHeaderHash(header *SpxHeader) bool {
	return header.IsValid()
}

func GetLodFlag() uint8 {
	return 0
}

func WriteSpxV1(spxFile string, rows []*SplatData, comment string, shDegree uint8) {
	WriteSpxOpenV1(spxFile, rows, comment, shDegree)
}

func ReadSpxV1(spxFile string, header *SpxHeader) (*SpxHeader, []*SplatData) {
	return ReadSpxOpenV1(spxFile, header)
}

func WriteSpxV2(spxFile string, rows []*SplatData, comment string, shDegree uint8) {
	WriteSpxOpenV2(spxFile, rows, comment, shDegree)
}

func ReadSpxV2(spxFile string, header *SpxHeader) (*SpxHeader, []*SplatData) {
	return ReadSpxOpenV2(spxFile, header)
}

func WriteSpxV3(spxFile string, rows []*SplatData, comment string, shDegree uint8) {
	WriteSpxOpenV3(spxFile, rows, comment, shDegree)
}

func ReadSpxV3(spxFile string, header *SpxHeader) (*SpxHeader, []*SplatData) {
	return ReadSpxOpenV3(spxFile, header)
}

func DefaultSpxComment() string {
	return "created by gsbox " + cmn.VER + " https://github.com/gotoeasy/gsbox"
}

func BlockFormatDesc(bf int) string {
	rs := ""
	switch bf {
	case BF_SPLAT22:
		rs = "(block format, splat per 20|22 bytes)"
	case BF_SPLAT23:
		rs = "(block format, splat per 20|22 bytes)"
	case BF_SPLAT220_WEBP:
		rs = "(block format, splat per 20|22 bytes, webp encoding)"
	case BF_SPLAT230_WEBP:
		rs = "(block format, splat per 20|22 bytes, webp encoding)"
	case BF_SPLAT19:
		rs = "(block format, splat per 19 bytes)"
	case BF_SPLAT10019:
		rs = "(block format, splat per 19 bytes)"
	case BF_SPLAT20:
		rs = "(block format, splat per 20 bytes)"
	case BF_SPLAT190_WEBP:
		rs = "(block format, splat per 19 bytes, webp encoding)"
	case BF_SPLAT10190_WEBP:
		rs = "(block format, splat per 19 bytes, webp encoding)"
	case BF_SH_PALETTES:
		rs = "(compress only)"
	case BF_SH_PALETTES_WEBP:
		rs = "(webp encoding)"
	}

	return rs
}
