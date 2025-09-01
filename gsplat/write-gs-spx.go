package gsplat

import "log"

const MinCompressBlockSize = 64

func WriteSpx(spxFile string, rows []*SplatData, comment string, shDegree int, flag1 uint8, flag2 uint8, flag3 uint8) {
	ver := Args.GetArgIgnorecase("-ov", "--output-version")
	if ver != "1" && ver != "2" && ver != "" {
		log.Println("[Warn] Ignore invalid output version:", ver)
	}

	if ver == "1" {
		log.Println("[Info] output spx version: 1")
		WriteSpxV1(spxFile, rows, comment, shDegree, flag1, flag2, flag3)
	} else {
		log.Println("[Info] output spx version: 2")
		WriteSpxV2(spxFile, rows, comment, shDegree, flag1, flag2, flag3)
	}
}
