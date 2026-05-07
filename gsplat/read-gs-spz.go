package gsplat

import (
	"encoding/binary"
	"errors"
	"gsbox/cmn"
	"io"
	"log"
	"os"
	"path/filepath"
)

func ReadSpz(spzFile string, readHeaderOnly ...bool) (*SpzHeader, []*SplatData) {
	isNetFile := cmn.IsNetFile(spzFile)
	if isNetFile {
		tmpdir, err := cmn.CreateTempDir()
		cmn.ExitOnError(err)
		downloadFile := filepath.Join(tmpdir, cmn.FileName(spzFile))
		log.Println("[Info]", "download start,", spzFile)
		err = cmn.HttpDownload(spzFile, downloadFile, nil)
		cmn.RemoveAllFileIfError(err, tmpdir)
		cmn.ExitOnError(err)
		log.Println("[Info]", "download finish")
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
	if !IsSpzV4(gzipDatas) {
		ungzipDatas, err := cmn.DecompressGzip(gzipDatas)
		cmn.ExitOnConditionError(err != nil, errors.New("[SPZ ERROR] UnGzip failed"))
		cmn.ExitOnConditionError(len(ungzipDatas) < HeaderSizeSpzV3, errors.New("[SPZ ERROR] Invalid spz header"))

		header := ParseSpzHeader(ungzipDatas[0:HeaderSizeSpzV3])
		if len(readHeaderOnly) > 0 && readHeaderOnly[0] {
			return header, nil
		}
		if header.Version != 2 && header.Version != 3 {
			cmn.ExitOnError(errors.New("[SPZ ERROR] deserializePackedGaussians: version not supported: " + cmn.Uint32ToString(header.Version)))
		}
		datas := readSpzDatasV2V3(ungzipDatas[HeaderSizeSpzV3:], header)
		OnProgress(PhaseRead, 100, 100)
		return header, datas
	}

	// v4
	header := ParseSpzHeader(gzipDatas[0:HeaderSizeSpzV4])
	if len(readHeaderOnly) > 0 && readHeaderOnly[0] {
		return header, nil
	}

	datas := readSpzDatasV4(gzipDatas, header)
	OnProgress(PhaseRead, 100, 100)
	return header, datas
}

func readSpzDatasV2V3(datas []byte, h *SpzHeader) []*SplatData {
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
		OnProgress(PhaseRead, i, int(h.NumPoints))
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
			data.RotationW, data.RotationX, data.RotationY, data.RotationZ = cmn.SpzDecodeRotationsV3V4(rotations[i*4 : i*4+4])
		}
		switch h.ShDegree {
		case 1:
			data.SH45 = InitZeroSH45()
			copy(data.SH45, shs[i*9:i*9+9])
		case 2:
			data.SH45 = InitZeroSH45()
			copy(data.SH45, shs[i*24:i*24+24])
		case 3:
			data.SH45 = InitZeroSH45()
			copy(data.SH45, shs[i*45:i*45+45])
		}

		splatDatas = append(splatDatas, data)
	}

	return splatDatas
}

// Header(32),Extensions(N*12),TOC(NumStreams*16),zstd datas
func readSpzDatasV4(datas []byte, h *SpzHeader) []*SplatData {

	datas = datas[h.TocByteOffset:]

	zstdSizes := make([]uint64, h.NumStreams)
	for i := range h.NumStreams {
		zstdSizes[i] = binary.LittleEndian.Uint64(datas[0:8])
		datas = datas[16:]
	}

	// Positions
	n := 0
	zstdBytes := datas[0:zstdSizes[n]]
	positions, err := cmn.DecompressZstd(zstdBytes)
	cmn.ExitOnError(err)
	datas = datas[zstdSizes[n]:]
	n++

	// Alphas
	zstdBytes = datas[0:zstdSizes[n]]
	alphas, err := cmn.DecompressZstd(zstdBytes)
	cmn.ExitOnError(err)
	datas = datas[zstdSizes[n]:]
	n++

	// Colors
	zstdBytes = datas[0:zstdSizes[n]]
	colors, err := cmn.DecompressZstd(zstdBytes)
	cmn.ExitOnError(err)
	datas = datas[zstdSizes[n]:]
	n++

	// Scales
	zstdBytes = datas[0:zstdSizes[n]]
	scales, err := cmn.DecompressZstd(zstdBytes)
	cmn.ExitOnError(err)
	datas = datas[zstdSizes[n]:]
	n++

	// Rotations
	zstdBytes = datas[0:zstdSizes[n]]
	rotations, err := cmn.DecompressZstd(zstdBytes)
	cmn.ExitOnError(err)
	datas = datas[zstdSizes[n]:]
	n++

	// SH
	var shs []byte
	if h.NumStreams > 5 {
		zstdBytes = datas[0:zstdSizes[n]]
		shs, err = cmn.DecompressZstd(zstdBytes)
		cmn.ExitOnError(err)
		datas = datas[zstdSizes[n]:]
		n++
	}

	var splatDatas []*SplatData
	for i := range int(h.NumPoints) {
		OnProgress(PhaseRead, i, int(h.NumPoints))
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
		data.RotationW, data.RotationX, data.RotationY, data.RotationZ = cmn.SpzDecodeRotationsV3V4(rotations[i*4 : i*4+4])

		switch h.ShDegree {
		case 1:
			data.SH45 = InitZeroSH45()
			copy(data.SH45, shs[i*9:i*9+9])
		case 2:
			data.SH45 = InitZeroSH45()
			copy(data.SH45, shs[i*24:i*24+24])
		case 3:
			data.SH45 = InitZeroSH45()
			copy(data.SH45, shs[i*45:i*45+45])
		}

		splatDatas = append(splatDatas, data)
	}

	return splatDatas
}
