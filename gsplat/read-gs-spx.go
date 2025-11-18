package gsplat

import (
	"errors"
	"gsbox/cmn"
	"log"
	"path/filepath"
)

var inputSpxHeader *SpxHeader

func ReadSpx(spxFile string) (*SpxHeader, []*SplatData) {
	isNetFile := cmn.IsNetFile(spxFile)
	if isNetFile {
		tmpdir, err := cmn.CreateTempDir()
		cmn.ExitOnError(err)
		downloadFile := filepath.Join(tmpdir, cmn.FileName(spxFile))
		log.Println("[Info]", "Download start,", spxFile)
		err = cmn.HttpDownload(spxFile, downloadFile, nil)
		cmn.RemoveAllFileIfError(err, tmpdir)
		cmn.ExitOnError(err)
		log.Println("[Info]", "Download finish")
		spxFile = downloadFile
		defer cmn.RemoveAllFile(tmpdir)
	}

	header := ParseSpxHeader(spxFile)
	inputSpxHeader = header
	if !header.IsValid() && IsOpenCreaterId(header.CreaterId) && IsOpenExclusiveId(header.ExclusiveId) {
		log.Println("[Warn] hash check failed! CreaterId:" + cmn.Uint32ToString(header.CreaterId) + ", ExclusiveId:" + cmn.Uint32ToString(header.ExclusiveId))
	}

	switch header.Version {
	case 1:
		return ReadSpxV1(spxFile, header)
	case 2:
		return ReadSpxV2(spxFile, header)
	case 3:
		return ReadSpxV3(spxFile, header)
	default:
		cmn.ExitOnError(errors.New("unsupport spx version: " + cmn.IntToString(int(header.Version))))
		return nil, nil
	}
}
