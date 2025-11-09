package gsplat

import (
	"errors"
	"gsbox/cmn"
	"log"
)

const MaxBlockSize = 1048576
const MinCompressBlockSize = 64
const MinWebpBlockSize = 4096
const DefaultBlockSize = 65536

func WriteSpx(spxFile string, rows []*SplatData) {

	ver := OutputSpxVersion()
	shDegree := GetArgShDegree()
	comment := Args.GetArgIgnorecase("-c", "--comment")

	log.Println("[Info] output spx version:", ver)
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
}
