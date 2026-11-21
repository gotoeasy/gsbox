package gsplat

import (
	"errors"
	"gsbox/cmn"
	"log"
	"path/filepath"
	"sync"
)

// 支持输出zip压缩单个sog文件，或压缩前的多个文件
func WriteSogV3(sogOrJsonFile string, rows []*SplatData) (fileSize int64) {
	ver := Args.GetArgIgnorecase("-ov", "--output-version")
	cmn.ExitOnConditionError(ver != "" && ver != "2" && ver != "3", errors.New("support sog version 2|3 only"))

	outputShDegree := GetArgShDegree()
	log.Println("[Info] output sog version: 3")
	log.Println("[Info] output shDegree:", outputShDegree)
	log.Println("[Info] quality level:", oArg.Quality, "(range 1~9)")
	if hasRgbPointCloudData {
		log.Println("[Info] NOTICE: Output contains color point cloud data")
	}

	dir := cmn.Dir(sogOrJsonFile)
	isSog := !cmn.Endwiths(sogOrJsonFile, "meta.json", true)
	if isSog {
		tmpdir, err := cmn.CreateTempDir()
		cmn.ExitOnError(err)
		dir = tmpdir
		defer func() {
			cmn.RemoveAllFile(tmpdir)
		}()
	}

	var wg sync.WaitGroup
	mm := &V3MinMax{}
	var paletteSize int

	wg.Add(5)
	go asyncWriteMeans(&wg, dir, rows, mm, true)
	go asyncWriteScales(&wg, dir, rows, true)
	go asyncWriteQuats(&wg, dir, rows, true)
	go asyncWriteSh0(&wg, dir, rows, true)
	go asyncWriteShN(&wg, dir, rows, &paletteSize, true)
	wg.Wait()

	fileMeansL := filepath.Join(dir, "means_l.webp")
	fileMeansU := filepath.Join(dir, "means_u.webp")
	fileScales := filepath.Join(dir, "scales.webp")
	fileQuats := filepath.Join(dir, "quats.webp")
	fileSh0 := filepath.Join(dir, "sh0.webp")
	files := []string{fileMeansL, fileMeansU, fileScales, fileQuats, fileSh0}
	if outputShDegree > 0 {
		fileCentroids := filepath.Join(dir, "shN_centroids.webp")
		fileLabels := filepath.Join(dir, "shN_labels.webp")
		files = append(files, fileCentroids, fileLabels)
	}
	fileMeta := filepath.Join(dir, "meta.json")
	files = append(files, fileMeta)

	writeMetaV3(dir, mm, paletteSize, len(rows))
	cmn.PrintLibwebpInfo(true)

	if isSog {
		cmn.ZipSogFiles(sogOrJsonFile, files)
	} else {
		fileSize = 0
		for _, f := range files {
			fileSize += cmn.GetFileSize(f)
		}
	}
	OnProgress(PhaseWrite, 100, 100)
	return
}

func writeQuatsV3(dir string, rows []*SplatData, printLogs ...bool) []string {
	printLog := len(printLogs) == 0 || printLogs[0]
	fileQuats := filepath.Join(dir, "quats.webp")
	if printLog {
		log.Println("[Info] generate quats.webp")
	}
	rgbaQuats := getQuatsRgbaV3(rows)
	bytsQuats, err := cmn.CompressWebp(rgbaQuats, oArg.webpQuality)
	cmn.ExitOnError(err)
	err = cmn.WriteFileBytes(fileQuats, bytsQuats)
	cmn.ExitOnError(err)
	return []string{fileQuats}
}

func getQuatsRgbaV3(rows []*SplatData) []uint8 {
	dataCnt := len(rows)
	rgba := make([]uint8, dataCnt*4)
	for i := range dataCnt {
		data := rows[i]
		idx := i * 4
		rgba[idx], rgba[idx+1], rgba[idx+2], rgba[idx+3] = cmn.EncodeRotationsExpMapSogV3(data.RotationW, data.RotationX, data.RotationY, data.RotationZ)
	}
	return rgba
}

func writeMetaV3(dir string, mm *V3MinMax, paletteSize int, count int, printLogs ...bool) []string {
	printLog := len(printLogs) == 0 || printLogs[0]
	m := new(Meta)
	m.Version = 3
	m.Asset = &Asset{Generator: "gsbox " + cmn.VER}
	m.Count = count
	_, comment := cmn.RemoveNonASCII(Args.GetArgIgnorecase("-c", "--comment"))
	m.Comment = comment
	var means Means
	means.Mins = []float32{mm.MinX, mm.MinY, mm.MinZ}
	means.Maxs = []float32{mm.MaxX, mm.MaxY, mm.MaxZ}
	means.Files = []string{"means_l.webp", "means_u.webp"}
	m.Means = means
	var scales Scales
	scales.Codebook = make([]float32, 256)
	for i := range 256 {
		scales.Codebook[i] = cmn.DecodeSpxScale(uint8(i))
	}
	scales.Files = []string{"scales.webp"}
	m.Scales = scales
	var quats Quats
	quats.Files = []string{"quats.webp"}
	m.Quats = quats
	var sh0 Sh0
	sh0.Codebook = make([]float32, 256)
	for i := range 256 {
		sh0.Codebook[i] = cmn.DecodeSplatColor(uint8(i))
	}
	sh0.Files = []string{"sh0.webp"}
	m.Sh0 = sh0
	if paletteSize > 0 && GetArgShDegree() > 0 {
		var shn ShN
		shn.Count = paletteSize
		shn.Bands = GetArgShDegree()
		shn.Codebook = make([]float32, 256)
		for i := range 256 {
			shn.Codebook[i] = cmn.DecodeSplatSH(uint8(i))
		}
		shn.Files = []string{"shN_centroids.webp", "shN_labels.webp"}
		m.ShN = &shn
	}

	bytsJson, err := m.ToJSON()
	cmn.ExitOnError(err)
	fileMeta := filepath.Join(dir, "meta.json")
	if printLog {
		log.Println("[Info] generate meta.json")
	}
	err = cmn.WriteFileBytes(fileMeta, bytsJson)
	cmn.ExitOnError(err)
	return []string{fileMeta}
}

func asyncWriteMeans(wg *sync.WaitGroup, dir string, rows []*SplatData, mm *V3MinMax, printLog bool) {
	defer wg.Done()
	_, MM := writeMeans(dir, rows, printLog)
	mm.CenterX = MM.CenterX
	mm.CenterY = MM.CenterY
	mm.CenterZ = MM.CenterZ
	mm.LenX = MM.LenX
	mm.LenY = MM.LenY
	mm.LenZ = MM.LenZ
	mm.MaxX = MM.MaxX
	mm.MaxY = MM.MaxY
	mm.MaxZ = MM.MaxZ
	mm.MinX = MM.MinX
	mm.MinY = MM.MinY
	mm.MinZ = MM.MinZ
}

func asyncWriteScales(wg *sync.WaitGroup, dir string, rows []*SplatData, printLog bool) {
	defer wg.Done()
	writeScales(dir, rows, printLog)
}

func asyncWriteQuats(wg *sync.WaitGroup, dir string, rows []*SplatData, printLog bool) {
	defer wg.Done()
	writeQuatsV3(dir, rows, printLog)
}

func asyncWriteSh0(wg *sync.WaitGroup, dir string, rows []*SplatData, printLog bool) {
	defer wg.Done()
	writeSh0(dir, rows, printLog)
}

func asyncWriteShN(wg *sync.WaitGroup, dir string, rows []*SplatData, paletteSize *int, printLog bool) {
	defer wg.Done()
	outputShDegree := GetArgShDegree()
	if outputShDegree > 0 {
		var shN_centroids []uint8
		var shN_labels []uint8
		shN_centroids, shN_labels, *paletteSize = ReWriteShByKmeans(rows)
		widths := []int{0, 96, 512, 960}
		bytsCentroids, err := cmn.CompressWebpByWidthHeight(shN_centroids, widths[outputShDegree], 1024, oArg.webpQuality)
		cmn.ExitOnError(err)
		bytsLabels, err := cmn.CompressWebp(shN_labels, oArg.webpQuality)
		cmn.ExitOnError(err)

		writeShN(dir, bytsCentroids, bytsLabels, printLog)
		log.Println("[Info] sh palette size", *paletteSize)
	}
}
