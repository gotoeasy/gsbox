package gsplat

import (
	"gsbox/cmn"
	"log"
	"path/filepath"
	"sort"
)

func WriteSogLodMeta(output string, splatTiles *SplatTiles, lodMeta *LodMeta) {
	defer OnProgress(PhaseWrite, 100, 100)

	outputShDegree := GetArgShDegree()

	log.Println("[Info] output sog version: 2")
	log.Println("[Info] output shDegree:", outputShDegree)
	log.Println("[Info] quality level:", oArg.Quality, "(range 1~9)")
	if splatTiles.PaletteSize > 0 {
		log.Println("[Info] sh palette size", splatTiles.PaletteSize)
	}

	// 取出文件，排序，设定调色板所在文件
	var files []*SplatFile
	for _, file := range splatTiles.Files {
		files = append(files, file)
	}
	sort.Slice(files, func(i, j int) bool {
		if files[i].Lod == files[j].Lod {
			return files[i].Seq < files[j].Seq
		}
		return files[i].Lod < files[j].Lod
	})

	// 写文件
	MaxFileCnt := len(files) + 2
	prcCnt := 0
	OnProgress(PhaseWrite, prcCnt, MaxFileCnt)

	dir := cmn.Dir(output)
	cmn.MkdirAll(cmn.Dir(dir))
	i := 0
	total := cmn.IntToString(len(files))
	isSog := Args.GetArgIgnorecase("-of", "--output-format") != "meta.json"
	for _, splatFile := range files {
		OnProgress(PhaseWrite, prcCnt, MaxFileCnt)

		if isSog {
			WriteLodMetaDataSog(filepath.Join(dir, splatFile.Url), splatFile, splatTiles, outputShDegree)
		} else {
			WriteLodMetaDataMetaJson(dir, splatFile, splatTiles, outputShDegree)
		}
		splatFile.Datas = nil
		i++
		prcCnt++
		log.Println("[Info]", "write", splatFile.Url, "("+cmn.IntToString(i)+"/"+total+")")
	}
	inputEnv := Args.GetArgIgnorecase("-e", "--env", "--environment")
	environment := "environment.sog"
	if !isSog {
		environment = "environment/environment.sog"
	}
	if saveSogLodEnvModel(dir, inputEnv, environment, splatTiles) {
		lodMeta.Environment = environment
		log.Println("[Info]", "write", environment)
	}
	OnProgress(PhaseWrite, prcCnt+1, MaxFileCnt)
	cmn.WriteFileString(output, lodMeta.ToJson())

	cmn.PrintLibwebpInfo(true)
}

func WriteLodMetaDataSog(sogFile string, splatFile *SplatFile, splatTiles *SplatTiles, outputShDegree uint8) {
	tmpdir, err := cmn.CreateTempDir()
	cmn.ExitOnError(err)
	dir := tmpdir
	defer func() {
		cmn.RemoveAllFile(tmpdir)
	}()

	files, mm := writeMeans(dir, splatFile.Datas, false)
	files = append(files, writeScales(dir, splatFile.Datas, false)...)
	files = append(files, writeQuats(dir, splatFile.Datas, false)...)
	files = append(files, writeSh0(dir, splatFile.Datas, false)...)
	files = append(files, writeMeta(dir, mm, splatTiles.PaletteSize, len(splatFile.Datas), false)...)
	if outputShDegree > 0 {
		widths := []int{0, 96, 512, 960}
		bytsCentroids, err := cmn.CompressWebpByWidthHeight(splatTiles.ShCentroids, widths[outputShDegree], 1024, oArg.webpQuality)
		cmn.ExitOnError(err)
		bytsLabels, err := cmn.CompressWebp(GetShnLablesByPaletteIdx(splatFile.Datas), oArg.webpQuality)
		cmn.ExitOnError(err)
		files = append(files, writeShN(dir, bytsCentroids, bytsLabels, false)...)
	}

	cmn.ZipSogFiles(sogFile, files, false)
}

func WriteLodMetaDataMetaJson(dir string, splatFile *SplatFile, splatTiles *SplatTiles, outputShDegree uint8) {
	distDir := cmn.Dir(filepath.Join(dir, splatFile.Url))
	files, mm := writeMeans(distDir, splatFile.Datas, false)
	files = append(files, writeScales(distDir, splatFile.Datas, false)...)
	files = append(files, writeQuats(distDir, splatFile.Datas, false)...)
	files = append(files, writeSh0(distDir, splatFile.Datas, false)...)
	files = append(files, writeMeta(distDir, mm, splatTiles.PaletteSize, len(splatFile.Datas), false)...)
	if outputShDegree > 0 {
		widths := []int{0, 96, 512, 960}
		bytsCentroids, err := cmn.CompressWebpByWidthHeight(splatTiles.ShCentroids, widths[outputShDegree], 1024, oArg.webpQuality)
		cmn.ExitOnError(err)
		bytsLabels, err := cmn.CompressWebp(GetShnLablesByPaletteIdx(splatFile.Datas), oArg.webpQuality)
		cmn.ExitOnError(err)
		files = append(files, writeShN(distDir, bytsCentroids, bytsLabels, false)...)
	}
}

func saveSogLodEnvModel(dir string, inputEnvFile string, outputFile string, splatTiles *SplatTiles) bool {
	if inputEnvFile == "" && len(splatTiles.EnvironmentDatas) == 0 {
		return false
	}

	if inputEnvFile != "" && !cmn.IsExistFile(inputEnvFile) {
		log.Println("[Warn] environment file not found:", inputEnvFile)
		return false // 暂不支持网络文件
	}

	var datas []*SplatData
	if inputEnvFile != "" {
		if cmn.Endwiths(inputEnvFile, ".ply", true) {
			_, datas = ReadPly(inputEnvFile)
		} else if cmn.Endwiths(inputEnvFile, ".splat", true) {
			datas = ReadSplat(inputEnvFile)
		} else if cmn.Endwiths(inputEnvFile, ".spx", true) {
			_, datas = ReadSpx(inputEnvFile)
		} else if cmn.Endwiths(inputEnvFile, ".spz", true) {
			_, datas = ReadSpz(inputEnvFile)
		} else if cmn.Endwiths(inputEnvFile, ".ksplat", true) {
			_, _, datas = ReadKsplat(inputEnvFile)
		} else if cmn.Endwiths(inputEnvFile, ".sog", true) || cmn.FileName(inputEnvFile) == "meta.json" {
			datas, _ = ReadSog(inputEnvFile)
		}
	} else {
		datas = splatTiles.EnvironmentDatas
	}

	sogFile := filepath.Join(dir, outputFile)
	warpSplatFile := &SplatFile{Datas: datas}
	writeLodMetaEnvironmentSog(sogFile, warpSplatFile)

	return true
}

// 输出单个sog文件，无球谐系数
func writeLodMetaEnvironmentSog(sogFile string, splatFile *SplatFile) {
	tmpdir, err := cmn.CreateTempDir()
	cmn.ExitOnError(err)
	dir := tmpdir
	defer func() {
		cmn.RemoveAllFile(tmpdir)
	}()

	files, mm := writeMeans(dir, splatFile.Datas, false)
	files = append(files, writeScales(dir, splatFile.Datas, false)...)
	files = append(files, writeQuats(dir, splatFile.Datas, false)...)
	files = append(files, writeSh0(dir, splatFile.Datas, false)...)
	files = append(files, writeMeta(dir, mm, 0, len(splatFile.Datas), false)...)

	cmn.ZipSogFiles(sogFile, files, false)
}
