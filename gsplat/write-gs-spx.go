package gsplat

import (
	"errors"
	"gsbox/cmn"
	"log"
)

const MaxBlockSize = int(1000000)    // 1000*1000，过大会增加前端wasm解析的内存压力，会对低端设备不友好，应不断观察调整为合适值
const MinBlockSize = int(4096)       // 64*64，过小性能反差，参数边界值应限制
const DefaultBlockSize = int(67600)  // 260*260，应不断观察调整为合适值
const MinCompressBlockSize = int(64) // 再小就别压了
const MinWebpBlockSize = int(4096)   // 64*64，再小就别webp了

func WriteSpx(spxFile string, rows []*SplatData) {

	ver := OutputSpxVersion()
	shDegree := GetArgShDegree()
	comment := Args.GetArgIgnorecase("-c", "--comment")

	log.Println("[Info] output spx version:", ver)
	if ver < NewestSpxVersion {
		log.Println("[Warn] it is highly recommended to migrate to the spx version", NewestSpxVersion)
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
