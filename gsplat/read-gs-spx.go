package gsplat

import (
	"gsbox/cmn"
	"log"
)

func ReadSpx(spxFile string) (*SpxHeader, []*SplatData) {

	header := ParseSpxHeader(spxFile)
	if !header.IsValid() && header.CreaterId == ID1202056903 && header.ExclusiveId == 0 {
		log.Println("[Warn] hash check failed! CreaterId:" + cmn.Uint32ToString(header.CreaterId) + ", ExclusiveId:" + cmn.Uint32ToString(header.ExclusiveId))
	}

	if header.Version == 1 {
		return ReadSpxV1(spxFile, header)
	}
	return ReadSpxV2(spxFile, header)
}
