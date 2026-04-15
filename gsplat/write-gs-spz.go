package gsplat

import (
	"bufio"
	"gsbox/cmn"
	"log"
	"os"
)

func WriteSpz(spzFile string, rows []*SplatData) {
	file, err := os.Create(spzFile)
	cmn.ExitOnError(err)
	defer file.Close()

	outputShDegree := GetArgShDegree()
	writer := bufio.NewWriter(file)
	ver := cmn.StringToInt(Args.GetArgIgnorecase("-ov", "--output-version"), 2)
	if ver < 2 || ver > 3 {
		log.Println("[Warn] Ignore invalid output version:", ver)
		ver = 2
	}
	log.Println("[Info] output spz version:", ver)
	log.Println("[Info] output shDegree:", outputShDegree)

	h := &SpzHeader{
		Magic:          SPZ_MAGIC,
		Version:        uint32(ver),
		NumPoints:      uint32(len(rows)),
		ShDegree:       uint8(outputShDegree),
		FractionalBits: 12,
		Flags:          0,
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

	for i := range rows {
		bts = append(bts, cmn.SpzEncodePosition(rows[i].PositionX)...)
		bts = append(bts, cmn.SpzEncodePosition(rows[i].PositionY)...)
		bts = append(bts, cmn.SpzEncodePosition(rows[i].PositionZ)...)
	}
	OnProgress(PhaseWrite, 15, 100)
	for i := range rows {
		bts = append(bts, rows[i].ColorA)
	}
	OnProgress(PhaseWrite, 30, 100)
	for i := range rows {
		bts = append(bts, cmn.SpzEncodeColor(rows[i].ColorR), cmn.SpzEncodeColor(rows[i].ColorG), cmn.SpzEncodeColor(rows[i].ColorB))
	}
	OnProgress(PhaseWrite, 45, 100)
	for i := range rows {
		bts = append(bts, cmn.SpzEncodeScale(rows[i].ScaleX), cmn.SpzEncodeScale(rows[i].ScaleY), cmn.SpzEncodeScale(rows[i].ScaleZ))
	}
	OnProgress(PhaseWrite, 60, 100)
	for i := range rows {
		if h.Version >= 3 {
			bts = append(bts, cmn.SpzEncodeRotationsV3(rows[i].RotationW, rows[i].RotationX, rows[i].RotationY, rows[i].RotationZ)...)
		} else {
			bts = append(bts, cmn.SpzEncodeRotations(rows[i].RotationW, rows[i].RotationX, rows[i].RotationY, rows[i].RotationZ)...)
		}
	}
	OnProgress(PhaseWrite, 75, 100)

	switch outputShDegree {
	case 1:
		for i := range rows {
			if len(rows[i].SH45) > 0 {
				for j := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH45[j]))
				}
			} else {
				for range 9 {
					bts = append(bts, 128)
				}
			}
		}
	case 2:
		for i := range rows {
			if len(rows[i].SH45) > 0 {
				for j := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH45[j]))
				}
				for j := 9; j < 24; j++ {
					bts = append(bts, cmn.SpzEncodeSH23(rows[i].SH45[j]))
				}
			} else {
				for range 24 {
					bts = append(bts, 128)
				}
			}
		}
	case 3:
		for i := range rows {
			if len(rows[i].SH45) > 0 {
				for j := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH45[j]))
				}
				for j := 9; j < 45; j++ {
					bts = append(bts, cmn.SpzEncodeSH23(rows[i].SH45[j]))
				}
			} else {
				for range 45 {
					bts = append(bts, 128)
				}
			}
		}
	}
	OnProgress(PhaseWrite, 90, 100)

	gzipDatas, err := cmn.CompressGzip(bts)
	cmn.ExitOnError(err)

	_, err = writer.Write(gzipDatas)
	cmn.ExitOnError(err)

	err = writer.Flush()
	cmn.ExitOnError(err)
	OnProgress(PhaseWrite, 100, 100)
}
