package gsplat

import (
	"errors"
	"gsbox/cmn"
	"log"
	"path/filepath"
	"sync"
)

func ReadSog(fileSogMeta string) ([]*SplatData, uint8) {
	isNetFile := cmn.IsNetFile(fileSogMeta)
	if isNetFile {
		if cmn.Endwiths(fileSogMeta, "/meta.json") {
			return readHttpSog(fileSogMeta)
		}

		tmpdir, err := cmn.CreateTempDir()
		cmn.ExitOnError(err)
		downloadSog := filepath.Join(tmpdir, cmn.FileName(fileSogMeta))
		log.Println("[Info]", "Download start,", fileSogMeta)
		cmn.HttpDownload(fileSogMeta, downloadSog, nil)
		log.Println("[Info]", "Download finish")
		fileSogMeta = downloadSog
	}

	dir := cmn.Dir(fileSogMeta)
	if cmn.Endwiths(fileSogMeta, ".sog", true) {
		tmpdir, err := cmn.CreateTempDir()
		cmn.ExitOnError(err)
		downloadTmpDir := dir
		dir = tmpdir

		defer func() {
			cmn.RemoveAllFile(tmpdir) // 清除解压的临时文件
			if isNetFile {
				cmn.RemoveAllFile(downloadTmpDir)
			}
		}()

		cmn.Unzip(fileSogMeta, tmpdir)
		fileSogMeta = filepath.Join(tmpdir, "meta.json")
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

func ReadSogInfo(fileSogMeta string) (version, count int, shDegree uint8, totalFileSize int64) {

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

func readHttpSog(urlMetaJson string) ([]*SplatData, uint8) {

	dir, err := cmn.CreateTempDir()
	cmn.ExitOnError(err)
	defer func() {
		cmn.RemoveAllFile(dir)
	}()

	fileMeta := filepath.Join(dir, "meta.json")
	cmn.HttpDownload(urlMetaJson, filepath.Join(dir, "meta.json"), nil)
	strMeta, err := cmn.ReadFileString(fileMeta)
	cmn.ExitOnError(err)

	meta, err := ParseSogMeta(strMeta)
	cmn.ExitOnError(err)

	var wg sync.WaitGroup
	log.Println("[Info]", "Download start")
	ary := cmn.Split(urlMetaJson, "/")
	log.Println("[Info]", meta.Means.Files[0])
	ary[len(ary)-1] = meta.Means.Files[0]
	wg.Add(1)
	go cmn.HttpDownload(cmn.Join(ary, "/"), filepath.Join(dir, meta.Means.Files[0]), &wg)
	log.Println("[Info]", meta.Means.Files[1])
	ary[len(ary)-1] = meta.Means.Files[1]
	wg.Add(1)
	go cmn.HttpDownload(cmn.Join(ary, "/"), filepath.Join(dir, meta.Means.Files[1]), &wg)
	log.Println("[Info]", meta.Scales.Files[0])
	ary[len(ary)-1] = meta.Scales.Files[0]
	wg.Add(1)
	go cmn.HttpDownload(cmn.Join(ary, "/"), filepath.Join(dir, meta.Scales.Files[0]), &wg)
	log.Println("[Info]", meta.Quats.Files[0])
	ary[len(ary)-1] = meta.Quats.Files[0]
	wg.Add(1)
	go cmn.HttpDownload(cmn.Join(ary, "/"), filepath.Join(dir, meta.Quats.Files[0]), &wg)
	log.Println("[Info]", meta.Sh0.Files[0])
	ary[len(ary)-1] = meta.Sh0.Files[0]
	wg.Add(1)
	go cmn.HttpDownload(cmn.Join(ary, "/"), filepath.Join(dir, meta.Sh0.Files[0]), &wg)
	if meta.ShN != nil {
		log.Println("[Info]", meta.ShN.Files[0])
		ary[len(ary)-1] = meta.ShN.Files[0]
		wg.Add(1)
		go cmn.HttpDownload(cmn.Join(ary, "/"), filepath.Join(dir, meta.ShN.Files[0]), &wg)
		log.Println("[Info]", meta.ShN.Files[1])
		ary[len(ary)-1] = meta.ShN.Files[1]
		wg.Add(1)
		go cmn.HttpDownload(cmn.Join(ary, "/"), filepath.Join(dir, meta.ShN.Files[1]), &wg)
	}
	wg.Wait() // 阻塞直到所有下载完成
	log.Println("[Info]", "Download finish")

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
