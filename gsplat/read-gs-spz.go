package gsplat

import (
	"errors"
	"gsbox/cmn"
	"io"
	"os"
)

func ReadSpz(spzFile string, readHeadOnly bool) (*SpzHeader, []*SplatData) {

	file, err := os.Open(spzFile)
	cmn.ExitOnConditionError(err != nil, errors.New("[SPZ ERROR] File open failed"))
	defer file.Close()

	gzipDatas, err := io.ReadAll(file)
	cmn.ExitOnConditionError(err != nil, errors.New("[SPZ ERROR] File read failed"))

	ungzipDatas, err := cmn.UnGzipBytes(gzipDatas)
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

	offsetpositions := 0
	offsetAlphas := offsetpositions + positionSize
	offsetColors := offsetAlphas + alphaSize
	offsetScales := offsetColors + colorSize
	offsetRotations := offsetScales + scaleSize
	offsetShs := offsetRotations + rotationSize

	shDim := 0
	if h.ShDegree == 1 {
		shDim = int(h.NumPoints) * 9
	} else if h.ShDegree == 2 {
		shDim = int(h.NumPoints) * 24
	} else if h.ShDegree == 3 {
		shDim = int(h.NumPoints) * 45
	}
	cmn.ExitOnConditionError(len(datas) != int(h.NumPoints)*19+shDim, errors.New("[SPZ ERROR] Invalid spz data"))

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
		data.RotationW, data.RotationX, data.RotationY, data.RotationZ = cmn.SpzDecodeRotations(rotations[i*3], rotations[i*3+1], rotations[i*3+2])
		if h.ShDegree == 1 {
			data.SH1 = shs[i*9 : i*9+9]
		} else if h.ShDegree == 2 {
			data.SH2 = shs[i*24 : i*24+24]
		} else if h.ShDegree == 3 {
			data.SH2 = shs[i*45 : i*45+24]
			data.SH3 = shs[i*45+24 : i*45+45]
		}

		splatDatas = append(splatDatas, data)
	}

	return splatDatas
}
