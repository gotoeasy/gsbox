package gsplat

import (
	"errors"
	"gsbox/cmn"
	"io"
	"log"
	"os"
	"path/filepath"
)

func ReadSpz(spzFile string, readHeadOnly bool) (*SpzHeader, []*SplatData) {
	isNetFile := cmn.IsNetFile(spzFile)
	if isNetFile {
		tmpdir, err := cmn.CreateTempDir()
		cmn.ExitOnError(err)
		downloadFile := filepath.Join(tmpdir, cmn.FileName(spzFile))
		log.Println("[Info]", "Download start,", spzFile)
		cmn.HttpDownload(spzFile, downloadFile, nil)
		log.Println("[Info]", "Download finish")
		spzFile = downloadFile
	}

	file, err := os.Open(spzFile)
	cmn.ExitOnConditionError(err != nil, errors.New("[SPZ ERROR] File open failed"))
	defer func() {
		file.Close()
		if isNetFile {
			cmn.RemoveAllFile(cmn.Dir(spzFile))
		}
	}()

	gzipDatas, err := io.ReadAll(file)
	cmn.ExitOnConditionError(err != nil, errors.New("[SPZ ERROR] File read failed"))

	ungzipDatas, err := cmn.DecompressGzip(gzipDatas)
	cmn.ExitOnConditionError(err != nil, errors.New("[SPZ ERROR] UnGzip failed"))
	cmn.ExitOnConditionError(len(ungzipDatas) < HeaderSizeSpz, errors.New("[SPZ ERROR] Invalid spz header"))

	header := ParseSpzHeader(ungzipDatas[0:HeaderSizeSpz])
	if readHeadOnly {
		return header, nil
	}

	datas := readSpzDatas(ungzipDatas[HeaderSizeSpz:], header)
	return header, datas
}

func readSpzDatas(datas []byte, h *SpzHeader) []*SplatData {
	positionSize := int(h.NumPoints) * 9
	alphaSize := int(h.NumPoints)
	colorSize := int(h.NumPoints) * 3
	scaleSize := int(h.NumPoints) * 3
	rotationSize := int(h.NumPoints) * 3
	size := 19
	if h.Version >= 3 {
		rotationSize = int(h.NumPoints) * 4
		size = 20
	}

	offsetpositions := 0
	offsetAlphas := offsetpositions + positionSize
	offsetColors := offsetAlphas + alphaSize
	offsetScales := offsetColors + colorSize
	offsetRotations := offsetScales + scaleSize
	offsetShs := offsetRotations + rotationSize

	shDim := 0
	switch h.ShDegree {
	case 1:
		shDim = int(h.NumPoints) * 9
	case 2:
		shDim = int(h.NumPoints) * 24
	case 3:
		shDim = int(h.NumPoints) * 45
	}
	cmn.ExitOnConditionError(len(datas) != int(h.NumPoints)*size+shDim, errors.New("[SPZ ERROR] Invalid spz data"))

	positions := datas[0:positionSize]
	scales := datas[offsetScales : offsetScales+scaleSize]
	rotations := datas[offsetRotations : offsetRotations+rotationSize]
	alphas := datas[offsetAlphas : offsetAlphas+alphaSize]
	colors := datas[offsetColors : offsetColors+colorSize]
	shs := datas[offsetShs:]

	var splatDatas []*SplatData
	for i := range int(h.NumPoints) {
		data := &SplatData{}
		data.PositionX = cmn.SpzDecodePosition(positions[i*9:i*9+3], h.FractionalBits)
		data.PositionY = cmn.SpzDecodePosition(positions[i*9+3:i*9+6], h.FractionalBits)
		data.PositionZ = cmn.SpzDecodePosition(positions[i*9+6:i*9+9], h.FractionalBits)
		data.ScaleX = cmn.SpzDecodeScale(scales[i*3])
		data.ScaleY = cmn.SpzDecodeScale(scales[i*3+1])
		data.ScaleZ = cmn.SpzDecodeScale(scales[i*3+2])
		data.ColorR = cmn.SpzDecodeColor(colors[i*3])
		data.ColorG = cmn.SpzDecodeColor(colors[i*3+1])
		data.ColorB = cmn.SpzDecodeColor(colors[i*3+2])
		data.ColorA = alphas[i]
		if h.Version == 2 {
			data.RotationW, data.RotationX, data.RotationY, data.RotationZ = cmn.SpzDecodeRotations(rotations[i*3], rotations[i*3+1], rotations[i*3+2])
		} else {
			data.RotationW, data.RotationX, data.RotationY, data.RotationZ = cmn.SpzDecodeRotationsV3(rotations[i*4 : i*4+4])
		}
		switch h.ShDegree {
		case 1:
			data.SH1 = shs[i*9 : i*9+9]
		case 2:
			data.SH2 = shs[i*24 : i*24+24]
		case 3:
			data.SH2 = shs[i*45 : i*45+24]
			data.SH3 = shs[i*45+24 : i*45+45]
		}

		splatDatas = append(splatDatas, data)
	}

	return splatDatas
}
