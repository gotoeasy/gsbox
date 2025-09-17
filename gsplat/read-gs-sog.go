package gsplat

import (
	"errors"
	"gsbox/cmn"
	"path/filepath"
)

func ReadSog(fileSogMeta string) ([]*SplatData, int) {

	dir := cmn.Dir(fileSogMeta)
	if cmn.Endwiths(fileSogMeta, ".sog", true) {
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

func ReadSogInfo(fileSogMeta string) (version, count, shDegree int, totalFileSize int64) {

	dir := cmn.Dir(fileSogMeta)
	isSog := cmn.Endwiths(fileSogMeta, ".sog", true)
	if isSog {
		tmpdir, err := cmn.CreateTempDir()
		cmn.ExitOnError(err)
		dir = tmpdir

		defer func() {
			cmn.RemoveAllFile(dir) // 清除解压的临时文件
		}()

		totalFileSize = cmn.GetFileSize(fileSogMeta)

		cmn.Unzip(fileSogMeta, dir)
		fileSogMeta = filepath.Join(dir, "meta.json")
	}

	strMeta, err := cmn.ReadFileString(fileSogMeta)
	cmn.ExitOnError(err)

	meta, err := ParseSogMeta(strMeta)
	cmn.ExitOnError(err)

	if meta.Version == 0 {
		version = 1
		count = meta.Means.Shape[0]

		shDegree = 0
		if meta.ShN != nil {
			switch meta.ShN.Shape[1] {
			case 45, 15:
				shDegree = 3
			case 24, 8:
				shDegree = 2
			case 9, 3:
				shDegree = 1
			}
		}
	} else {
		version = meta.Version
		count = meta.Count
		if meta.ShN == nil {
			shDegree = 0
		} else {
			shDegree = 3
		}
	}

	if !isSog {
		totalFileSize = cmn.GetFileSize(fileSogMeta)
		totalFileSize += cmn.GetFileSize(filepath.Join(dir, meta.Means.Files[0]))
		totalFileSize += cmn.GetFileSize(filepath.Join(dir, meta.Means.Files[1]))
		totalFileSize += cmn.GetFileSize(filepath.Join(dir, meta.Scales.Files[0]))
		totalFileSize += cmn.GetFileSize(filepath.Join(dir, meta.Quats.Files[0]))
		totalFileSize += cmn.GetFileSize(filepath.Join(dir, meta.Sh0.Files[0]))
		if meta.ShN != nil {
			totalFileSize += cmn.GetFileSize(filepath.Join(dir, meta.ShN.Files[0]))
			totalFileSize += cmn.GetFileSize(filepath.Join(dir, meta.ShN.Files[1]))
		}
	}

	return
}
