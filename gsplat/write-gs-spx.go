package gsplat

import "log"

const MaxBlockSize = 1048576
const MinCompressBlockSize = 64
const MinWebpBlockSize = 4096
const DefaultBlockSize = 65536

func WriteSpx(spxFile string, rows []*SplatData) {
	ver := Args.GetArgIgnorecase("-ov", "--output-version")
	if ver != "1" && ver != "2" && ver != "" {
		log.Println("[Warn] Ignore invalid output version:", ver)
	}

	comment := Args.GetArgIgnorecase("-c", "--comment")
	shDegree := GetArgShDegree()

	if ver == "1" {
		log.Println("[Info] output spx version: 1")
		WriteSpxOpenV1(spxFile, rows, comment, shDegree)
	} else {
		log.Println("[Info] output spx version: 2")
		WriteSpxV2(spxFile, rows, comment, shDegree)
	}
}
