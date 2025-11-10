package gsplat

import (
	"encoding/json"
	"errors"
	"gsbox/cmn"
	"log"
	"path/filepath"
)

// 支持输出zip压缩单个sog文件，或压缩前的多个文件
func WriteSog(sogOrJsonFile string, rows []*SplatData) (fileSize int64) {
	ver := Args.GetArgIgnorecase("-ov", "--output-version")
	cmn.ExitOnConditionError(ver != "" && ver != "2", errors.New("support sog version 2 only"))

	shDegree := GetArgShDegree()
	log.Println("[Info] output sog version: 2")
	log.Println("[Info] output shDegree:", shDegree)
	GetArgKmeansIterations(true)
	GetArgKmeansNearestNodes(true)

	dir := cmn.Dir(sogOrJsonFile)
	isSog := cmn.Endwiths(sogOrJsonFile, ".sog", true)
	if isSog {
		tmpdir, err := cmn.CreateTempDir()
		cmn.ExitOnError(err)
		dir = tmpdir
		defer func() {
			cmn.RemoveAllFile(tmpdir)
		}()
	}

	files, mm := writeMeans(dir, rows)
	files = append(files, writeScales(dir, rows)...)
	files = append(files, writeQuats(dir, rows)...)
	files = append(files, writeSh0(dir, rows)...)
	if shDegree > 0 {
		shN_centroids, shN_labels := ReWriteShByKmeans(rows)
		bytsCentroids, err := cmn.CompressWebpByWidthHeight(shN_centroids, 960, 1024)
		cmn.ExitOnError(err)
		bytsLabels, err := cmn.CompressWebp(shN_labels)
		cmn.ExitOnError(err)

		files = append(files, writeShN(dir, bytsCentroids, bytsLabels)...)
	}
	files = append(files, writeMeta(dir, mm, len(rows))...)

	if !isSog {
		fileSize = 0
		for _, f := range files {
			fileSize += cmn.GetFileSize(f)
		}
	}

	if isSog {
		cmn.ZipSogFiles(sogOrJsonFile, files)
	}
	return
}

func writeMeta(dir string, mm *V3MinMax, count int) []string {
	m := new(Meta)
	m.Version = 2
	m.Count = count
	_, comment := cmn.RemoveNonASCII(Args.GetArgIgnorecase("-c", "--comment"))
	if comment == "" {
		comment = DefaultSpxComment()
	}
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
	shDegree := GetArgShDegree()
	if shDegree > 0 {
		var shn ShN
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
	log.Println("[Info] generate meta.json")
	err = cmn.WriteFileBytes(fileMeta, bytsJson)
	cmn.ExitOnError(err)
	return []string{fileMeta}
}

func writeMeans(dir string, rows []*SplatData) ([]string, *V3MinMax) {
	dataCnt := len(rows)
	mm := ComputeXyzLogMinMax(rows)
	rgbaMeansL := make([]uint8, dataCnt*4)
	rgbaMeansU := make([]uint8, dataCnt*4)
	for i := range dataCnt {
		data := rows[i]
		x := cmn.ClipUint16(65535.0*(cmn.SogEncodeLog(data.PositionX)-mm.MinX)/mm.LenX + 0.5)
		y := cmn.ClipUint16(65535.0*(cmn.SogEncodeLog(data.PositionY)-mm.MinY)/mm.LenY + 0.5)
		z := cmn.ClipUint16(65535.0*(cmn.SogEncodeLog(data.PositionZ)-mm.MinZ)/mm.LenZ + 0.5)
		rgbaMeansL[i*4+0] = uint8(x & 0xFF)
		rgbaMeansL[i*4+1] = uint8(y & 0xFF)
		rgbaMeansL[i*4+2] = uint8(z & 0xFF)
		rgbaMeansL[i*4+3] = 255
		rgbaMeansU[i*4+0] = uint8(x >> 8)
		rgbaMeansU[i*4+1] = uint8(y >> 8)
		rgbaMeansU[i*4+2] = uint8(z >> 8)
		rgbaMeansU[i*4+3] = 255
	}

	fileMeansL := filepath.Join(dir, "means_l.webp")
	log.Println("[Info] generate means_l.webp")
	bytsMeansL, err := cmn.CompressWebp(rgbaMeansL)
	cmn.ExitOnError(err)
	err = cmn.WriteFileBytes(fileMeansL, bytsMeansL)
	cmn.ExitOnError(err)

	fileMeansU := filepath.Join(dir, "means_u.webp")
	log.Println("[Info] generate means_u.webp")
	bytsMeansU, err := cmn.CompressWebp(rgbaMeansU)
	cmn.ExitOnError(err)
	err = cmn.WriteFileBytes(fileMeansU, bytsMeansU)
	cmn.ExitOnError(err)

	return []string{fileMeansL, fileMeansU}, mm
}

func writeScales(dir string, rows []*SplatData) []string {
	fileScales := filepath.Join(dir, "scales.webp")
	log.Println("[Info] generate scales.webp")
	rgbaScales := getScalesRgba(rows)
	bytsScales, err := cmn.CompressWebp(rgbaScales)
	cmn.ExitOnError(err)
	err = cmn.WriteFileBytes(fileScales, bytsScales)
	cmn.ExitOnError(err)
	return []string{fileScales}
}

func writeQuats(dir string, rows []*SplatData) []string {
	fileQuats := filepath.Join(dir, "quats.webp")
	log.Println("[Info] generate quats.webp")
	rgbaQuats := getQuatsRgba(rows)
	bytsQuats, err := cmn.CompressWebp(rgbaQuats)
	cmn.ExitOnError(err)
	err = cmn.WriteFileBytes(fileQuats, bytsQuats)
	cmn.ExitOnError(err)
	return []string{fileQuats}
}

func writeSh0(dir string, rows []*SplatData) []string {
	fileSh0 := filepath.Join(dir, "sh0.webp")
	log.Println("[Info] generate sh0.webp")
	rgbaSh0 := getSh0Rgba(rows)
	bytsSh0, err := cmn.CompressWebp(rgbaSh0)
	cmn.ExitOnError(err)
	err = cmn.WriteFileBytes(fileSh0, bytsSh0)
	cmn.ExitOnError(err)
	return []string{fileSh0}
}

func writeShN(dir string, bytsCentroids []uint8, bytsLabels []uint8) []string {
	fileCentroids := filepath.Join(dir, "shN_centroids.webp")
	log.Println("[Info] generate shN_centroids.webp")
	err := cmn.WriteFileBytes(fileCentroids, bytsCentroids)
	cmn.ExitOnError(err)

	fileLabels := filepath.Join(dir, "shN_labels.webp")
	log.Println("[Info] generate shN_labels.webp")
	err = cmn.WriteFileBytes(fileLabels, bytsLabels)
	cmn.ExitOnError(err)
	return []string{fileCentroids, fileLabels}
}

func getSh0Rgba(rows []*SplatData) []uint8 {
	dataCnt := len(rows)
	rgba := make([]uint8, dataCnt*4)
	for i := range dataCnt {
		data := rows[i]
		rgba[i*4+0] = data.ColorR
		rgba[i*4+1] = data.ColorG
		rgba[i*4+2] = data.ColorB
		rgba[i*4+3] = data.ColorA
	}
	return rgba
}

func getScalesRgba(rows []*SplatData) []uint8 {
	dataCnt := len(rows)
	rgba := make([]uint8, dataCnt*4)
	for i := range dataCnt {
		data := rows[i]
		rgba[i*4+0] = cmn.EncodeSpxScale(data.ScaleX)
		rgba[i*4+1] = cmn.EncodeSpxScale(data.ScaleY)
		rgba[i*4+2] = cmn.EncodeSpxScale(data.ScaleZ)
		rgba[i*4+3] = 255
	}
	return rgba
}

func getQuatsRgba(rows []*SplatData) []uint8 {
	dataCnt := len(rows)
	rgba := make([]uint8, dataCnt*4)
	for i := range dataCnt {
		data := rows[i]
		idx := i * 4
		rgba[idx], rgba[idx+1], rgba[idx+2], rgba[idx+3] = cmn.SogEncodeRotations(data.RotationW, data.RotationX, data.RotationY, data.RotationZ)
	}
	return rgba
}

// 当前仅针对 sog version 2
type Meta struct {
	Version int    `json:"version"`
	Count   int    `json:"count"`
	Comment string `json:"comment,omitempty"`
	Means   Means  `json:"means"`
	Scales  Scales `json:"scales"`
	Quats   Quats  `json:"quats"`
	Sh0     Sh0    `json:"sh0"`
	ShN     *ShN   `json:"shN,omitempty"`
}

func (m *Meta) ToJSON() ([]byte, error) {
	// return json.MarshalIndent(m, "", "  ") // 生成带缩进的格式
	return json.Marshal(m)
}

type Means struct {
	Mins  []float32 `json:"mins"`
	Maxs  []float32 `json:"maxs"`
	Files []string  `json:"files"`
}

type Scales struct {
	Codebook []float32 `json:"codebook"`
	Files    []string  `json:"files"`
}

type Quats struct {
	Files []string `json:"files"`
}

type Sh0 struct {
	Codebook []float32 `json:"codebook"`
	Files    []string  `json:"files"`
}

type ShN struct {
	Count    int       `json:"count"`
	Bands    int       `json:"bands"`
	Codebook []float32 `json:"codebook"`
	Files    []string  `json:"files"`
}
