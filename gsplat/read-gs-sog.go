package gsplat

import (
	"errors"
	"gsbox/cmn"
	"path/filepath"
)

func ReadSog(fileSogMeta string) ([]*SplatData, int) {

	dir := cmn.Dir(fileSogMeta)
	if cmn.Endwiths(fileSogMeta, ".sog") {
		tmpdir, err := cmn.CreateTempDir()
		cmn.ExitOnError(err)
		dir = tmpdir

		defer func() {
			cmn.RemoveAllFile(dir) // 清除解压的临时文件
		}()

		cmn.Unzip(fileSogMeta, dir)
		fileSogMeta = filepath.Join(dir, "meta.json")
	}

	strMeta, err := cmn.ReadFileString(fileSogMeta)
	cmn.ExitOnError(err)

	meta, err := ParseSogMeta(strMeta)
	cmn.ExitOnError(err)

	switch meta.Version {
	case 0:
		return ReadSogV1(meta, dir)
	case 2:
		return ReadSogV2(meta, dir)
	default:
		cmn.ExitOnError(errors.New("unsupported sog version"))
	}
	return nil, 0
}
