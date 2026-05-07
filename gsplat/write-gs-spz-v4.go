package gsplat

import (
	"bufio"
	"gsbox/cmn"
	"log"
	"os"
)

func writeSpzV4(spzFile string, rows []*SplatData) {
	file, err := os.Create(spzFile)
	cmn.ExitOnError(err)
	defer file.Close()

	outputShDegree := GetArgShDegree()
	writer := bufio.NewWriter(file)
	if !oArg.isCut {
		log.Println("[Info] output spz version:", 4)
		log.Println("[Info] output shDegree:", outputShDegree)
	}

	numStreams := uint8(5)
	if outputShDegree > 0 {
		numStreams = 6
	}

	h := &SpzHeader{
		Magic:          SPZ_MAGIC,
		Version:        4,
		NumPoints:      uint32(len(rows)),
		ShDegree:       uint8(outputShDegree),
		FractionalBits: 12,
		Flags:          0,
		NumStreams:     numStreams,
		TocByteOffset:  32, // 默认不支持扩充数据
		Reserved:       0,
	}

	if outputShDegree > 0 {
		log.Println("[Info] quality level:", oArg.Quality, "(range 1~9)")
		if oArg.Quality < 9 {
			// 小于9级使用聚类，第9级按通常处理
			_, _, paletteSize := ReWriteShByKmeans(rows)
			log.Println("[Info] sh palette size", paletteSize)
		}
	}

	bts := make([]byte, 0)
	bts = append(bts, h.ToBytes()...)

	positions := make([]byte, 0)
	alphas := make([]byte, 0)
	colors := make([]byte, 0)
	scales := make([]byte, 0)
	rotations := make([]byte, 0)
	shs := make([]byte, 0)

	for i := range rows {
		positions = append(positions, cmn.SpzEncodePosition(rows[i].PositionX)...)
		positions = append(positions, cmn.SpzEncodePosition(rows[i].PositionY)...)
		positions = append(positions, cmn.SpzEncodePosition(rows[i].PositionZ)...)
	}
	OnProgress(PhaseWrite, 15, 100)
	for i := range rows {
		alphas = append(alphas, rows[i].ColorA)
	}
	OnProgress(PhaseWrite, 30, 100)
	for i := range rows {
		colors = append(colors, cmn.SpzEncodeColor(rows[i].ColorR), cmn.SpzEncodeColor(rows[i].ColorG), cmn.SpzEncodeColor(rows[i].ColorB))
	}
	OnProgress(PhaseWrite, 45, 100)
	for i := range rows {
		scales = append(scales, cmn.SpzEncodeScale(rows[i].ScaleX), cmn.SpzEncodeScale(rows[i].ScaleY), cmn.SpzEncodeScale(rows[i].ScaleZ))
	}
	OnProgress(PhaseWrite, 60, 100)
	for i := range rows {
		rotations = append(rotations, cmn.SpzEncodeRotationsV3V4(rows[i].RotationW, rows[i].RotationX, rows[i].RotationY, rows[i].RotationZ)...)
	}
	OnProgress(PhaseWrite, 75, 100)

	switch outputShDegree {
	case 1:
		for i := range rows {
			if len(rows[i].SH45) > 0 {
				for j := range 9 {
					shs = append(shs, cmn.SpzEncodeSH1(rows[i].SH45[j]))
				}
			} else {
				for range 9 {
					shs = append(shs, 128)
				}
			}
		}
	case 2:
		for i := range rows {
			if len(rows[i].SH45) > 0 {
				for j := range 9 {
					shs = append(shs, cmn.SpzEncodeSH1(rows[i].SH45[j]))
				}
				for j := 9; j < 24; j++ {
					shs = append(shs, cmn.SpzEncodeSH23(rows[i].SH45[j]))
				}
			} else {
				for range 24 {
					shs = append(shs, 128)
				}
			}
		}
	case 3:
		for i := range rows {
			if len(rows[i].SH45) > 0 {
				for j := range 9 {
					shs = append(shs, cmn.SpzEncodeSH1(rows[i].SH45[j]))
				}
				for j := 9; j < 45; j++ {
					shs = append(shs, cmn.SpzEncodeSH23(rows[i].SH45[j]))
				}
			} else {
				for range 45 {
					shs = append(shs, 128)
				}
			}
		}
	}
	OnProgress(PhaseWrite, 90, 100)

	zstdPositions := cmn.CompressZstd(positions)
	zstdAlphas := cmn.CompressZstd(alphas)
	zstdColors := cmn.CompressZstd(colors)
	zstdScales := cmn.CompressZstd(scales)
	zstdRotations := cmn.CompressZstd(rotations)
	var zstdShs []byte

	bts = append(bts, cmn.Uint64ToBytes(uint64(len(zstdPositions)))...)
	bts = append(bts, cmn.Uint64ToBytes(uint64(len(positions)))...)
	bts = append(bts, cmn.Uint64ToBytes(uint64(len(zstdAlphas)))...)
	bts = append(bts, cmn.Uint64ToBytes(uint64(len(alphas)))...)
	bts = append(bts, cmn.Uint64ToBytes(uint64(len(zstdColors)))...)
	bts = append(bts, cmn.Uint64ToBytes(uint64(len(colors)))...)
	bts = append(bts, cmn.Uint64ToBytes(uint64(len(zstdScales)))...)
	bts = append(bts, cmn.Uint64ToBytes(uint64(len(scales)))...)
	bts = append(bts, cmn.Uint64ToBytes(uint64(len(zstdRotations)))...)
	bts = append(bts, cmn.Uint64ToBytes(uint64(len(rotations)))...)
	if outputShDegree > 0 {
		zstdShs = cmn.CompressZstd(shs)
		bts = append(bts, cmn.Uint64ToBytes(uint64(len(zstdShs)))...)
		bts = append(bts, cmn.Uint64ToBytes(uint64(len(shs)))...)
	}

	bts = append(bts, zstdPositions...)
	bts = append(bts, zstdAlphas...)
	bts = append(bts, zstdColors...)
	bts = append(bts, zstdScales...)
	bts = append(bts, zstdRotations...)
	if outputShDegree > 0 {
		bts = append(bts, zstdShs...)
	}

	_, err = writer.Write(bts)
	cmn.ExitOnError(err)

	err = writer.Flush()
	cmn.ExitOnError(err)
	OnProgress(PhaseWrite, 100, 100)
}
