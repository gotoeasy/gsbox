package gsplat

import (
	"encoding/json"
	"errors"
	"gsbox/cmn"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

func ReadLodMeta(input string) (dataShDegree uint8, lodLevels uint16, lodDatas []*SplatData, envDatas []*SplatData) {
	cmn.ExitOnConditionError(!cmn.Endwiths(input, "lod-meta.json"), errors.New("invalid lod-meta.json: "+input))

	isNetFile := cmn.IsNetFile(input)
	tmpdir := ""
	inputTmp := input
	if isNetFile {
		tmpdir = downloadLod(input)
		inputTmp = filepath.Join(tmpdir, cmn.FileName(input))
	}
	defer func() {
		if isNetFile {
			cmn.RemoveAllFile(cmn.Dir(tmpdir))
		}
	}()

	jsonStr, err := cmn.ReadFileString(inputTmp)
	cmn.ExitOnError(err)

	var lodMeta LodMeta
	err = json.Unmarshal([]byte(jsonStr), &lodMeta)
	cmn.ExitOnError(err)

	lodMetaDir := cmn.Dir(inputTmp)

	// 汇总文件
	mapFiles := make(map[string]*SplatFile)
	for i, filename := range lodMeta.Filenames {
		splatFile := &SplatFile{FileKey: cmn.IntToString(i), Url: filepath.Join(lodMetaDir, filename)}
		mapFiles[splatFile.FileKey] = splatFile
	}

	// 设定文件LOD级别
	traveTree(lodMeta.Tree, func(node *LodNode) bool {
		if node.Lods != nil {
			for k, v := range *node.Lods {
				fileKey := cmn.IntToString(v.File)
				mapFiles[fileKey].Lod = uint16(cmn.StringToInt(k))
			}
		}
		return true
	})

	// 计算LOD最大级别
	lodLevels = 0
	for _, splatFile := range mapFiles {
		lodLevels = max(lodLevels, splatFile.Lod+1)
	}

	// 数据部分只支持 *.sog 或 meta.json
	lodDatas = make([]*SplatData, 0)
	for _, splatFile := range mapFiles {
		datas, header := ReadSog(splatFile.Url)
		dataShDegree = max(dataShDegree, header.ShDegree)
		for _, d := range datas {
			d.Lod = splatFile.Lod
		}
		lodDatas = append(lodDatas, datas...)
	}

	// 环境天空盒，默认 *.sog，保险起见支持多种格式
	if lodMeta.Environment != "" {
		fileEnv := filepath.Join(lodMetaDir, lodMeta.Environment)
		if cmn.Endwiths(fileEnv, ".sog", true) || cmn.FileName(fileEnv) == "meta.json" {
			envDatas, _ = ReadSog(fileEnv)
		} else if cmn.Endwiths(fileEnv, ".spx", true) {
			_, envDatas = ReadSpx(fileEnv)
		} else if cmn.Endwiths(fileEnv, ".spz", true) {
			_, envDatas = ReadSpz(fileEnv)
		} else if cmn.Endwiths(fileEnv, ".ply", true) {
			_, envDatas = ReadPly(fileEnv)
		} else if cmn.Endwiths(fileEnv, ".splat", true) {
			envDatas = ReadSplat(fileEnv)
		}
	}

	return
}

func traveTree(node *LodNode, fnCallBack func(nd *LodNode) bool) {
	if fnCallBack(node) && node.Children != nil {
		for _, cld := range *node.Children {
			traveTree(cld, fnCallBack)
		}
	}
}

func downloadLod(input string) string {
	tmpdir, err := cmn.CreateTempDir()
	cmn.ExitOnError(err)
	jsonFileName := cmn.FileName(input)
	localJsonFile := filepath.Join(tmpdir, jsonFileName) // lod json
	log.Println("[Info]", "download", jsonFileName)
	err = cmn.HttpDownload(input, localJsonFile, nil)
	cmn.RemoveAllFileIfError(err, tmpdir)
	cmn.ExitOnError(err)

	jsonStr, err := cmn.ReadFileString(localJsonFile)
	cmn.ExitOnError(err)

	var lodMeta LodMeta
	err = json.Unmarshal([]byte(jsonStr), &lodMeta)
	cmn.ExitOnError(err)

	ary := strings.Split(input, "/")
	ary = ary[:len(ary)-1]
	baseUrl := cmn.Join(ary, "/")

	// 汇总文件
	var files []string
	for _, filename := range lodMeta.Filenames {
		urlFile := baseUrl + "/" + filename
		localFile := filepath.Join(tmpdir, filename)
		localDir := cmn.Dir(localFile)
		log.Println("[Info]", "download", filename)
		if cmn.Endwiths(filename, "meta.json") {
			downloadSogOfMetaJson(urlFile, localDir)
		} else {
			err := cmn.HttpDownload(urlFile, localFile, nil)
			cmn.ExitOnError(err)
		}
	}
	if lodMeta.Environment != "" {
		files = append(files, filepath.Join(baseUrl, baseUrl+"/"+lodMeta.Environment))
		urlFile := baseUrl + "/" + lodMeta.Environment
		localFile := filepath.Join(tmpdir, lodMeta.Environment)
		localDir := cmn.Dir(localFile)
		log.Println("[Info]", "download", lodMeta.Environment)
		if cmn.Endwiths(lodMeta.Environment, "meta.json") {
			downloadSogOfMetaJson(urlFile, localDir)
		} else {
			err := cmn.HttpDownload(urlFile, localFile, nil)
			cmn.ExitOnError(err)
		}
	}

	return tmpdir
}

func downloadSogOfMetaJson(urlMetaJson string, saveDir string) {

	fileMeta := filepath.Join(saveDir, "meta.json")
	err := cmn.HttpDownload(urlMetaJson, filepath.Join(saveDir, "meta.json"), nil)
	cmn.ExitOnError(err)
	strMeta, err := cmn.ReadFileString(fileMeta)
	cmn.ExitOnError(err)

	meta, err := ParseSogMeta(strMeta)
	cmn.ExitOnError(err)

	var wg sync.WaitGroup
	ary := cmn.Split(urlMetaJson, "/")
	ary[len(ary)-1] = meta.Means.Files[0]
	wg.Add(1)
	go cmn.HttpDownload(cmn.Join(ary, "/"), filepath.Join(saveDir, meta.Means.Files[0]), &wg)
	ary[len(ary)-1] = meta.Means.Files[1]
	wg.Add(1)
	go cmn.HttpDownload(cmn.Join(ary, "/"), filepath.Join(saveDir, meta.Means.Files[1]), &wg)
	ary[len(ary)-1] = meta.Scales.Files[0]
	wg.Add(1)
	go cmn.HttpDownload(cmn.Join(ary, "/"), filepath.Join(saveDir, meta.Scales.Files[0]), &wg)
	ary[len(ary)-1] = meta.Quats.Files[0]
	wg.Add(1)
	go cmn.HttpDownload(cmn.Join(ary, "/"), filepath.Join(saveDir, meta.Quats.Files[0]), &wg)
	ary[len(ary)-1] = meta.Sh0.Files[0]
	wg.Add(1)
	go cmn.HttpDownload(cmn.Join(ary, "/"), filepath.Join(saveDir, meta.Sh0.Files[0]), &wg)
	if meta.ShN != nil {
		ary[len(ary)-1] = meta.ShN.Files[0]
		wg.Add(1)
		go cmn.HttpDownload(cmn.Join(ary, "/"), filepath.Join(saveDir, meta.ShN.Files[0]), &wg)
		ary[len(ary)-1] = meta.ShN.Files[1]
		wg.Add(1)
		go cmn.HttpDownload(cmn.Join(ary, "/"), filepath.Join(saveDir, meta.ShN.Files[1]), &wg)
	}
	wg.Wait() // 阻塞直到所有下载完成
}
