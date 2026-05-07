package gsplat

import (
	"gsbox/cmn"
	"log"
)

func WriteSpz(spzFile string, rows []*SplatData) {
	ver := cmn.StringToInt(Args.GetArgIgnorecase("-ov", "--output-version"), 2)
	if ver < 2 || ver > 4 {
		log.Println("[Warn] Ignore invalid output version:", ver)
		ver = 2
	}

	if ver < 4 {
		writeSpzV2V3(spzFile, rows, ver)
	} else {
		writeSpzV4(spzFile, rows)
	}
}
