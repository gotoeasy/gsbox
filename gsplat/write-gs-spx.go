package gsplat

import (
	"errors"
	"gsbox/cmn"
	"log"
)

const MaxBlockSize = int(1600 * 10000) // 不做特别限制，注意过大对低端设备不友好，需自行注意
const MinBlockSize = int(4096)         // 64*64，过小性能反差，参数边界值应限制
const DefaultBlockSize = int(67600)    // 260*260，应不断观察调整为合适值
const MinCompressBlockSize = int(64)   // 再小就别压了
const MinWebpBlockSize = int(4096)     // 64*64，再小就别webp了

func WriteSpx(spxFile string, rows []*SplatData) {

	ver := OutputSpxVersion()
	shDegree := GetArgShDegree()
	comment := Args.GetArgIgnorecase("-c", "--comment")

	log.Println("[Info] output spx version:", ver)
	if ver < NewestSpxVersion {
		log.Println("[Warn] IT IS HIGHLY RECOMMENDED TO MIGRATE TO THE SPX VERSION", NewestSpxVersion)
	}
	log.Println("[Info] output shDegree:", shDegree)

	switch ver {
	case 1:
		WriteSpxV1(spxFile, rows, comment, shDegree)
	case 2:
		WriteSpxV2(spxFile, rows, comment, shDegree)
	case 3:
		WriteSpxV3(spxFile, rows, comment, shDegree)
	default:
		cmn.ExitOnError(errors.New("unknow output version: " + cmn.IntToString(ver)))
	}

	OnProgress(PhaseWrite, 100, 100)
}
