package gsplat

import (
	"bufio"
	"gsbox/cmn"
	"log"
	"os"
)

func WriteSpz(spzFile string, rows []*SplatData, shDegree int) {
	file, err := os.Create(spzFile)
	cmn.ExitOnError(err)
	defer file.Close()

	log.Println("[Info] output shDegree:", shDegree)
	writer := bufio.NewWriter(file)

	h := &SpzHeader{
		Magic:          SPZ_MAGIC,
		Version:        2,
		NumPoints:      uint32(len(rows)),
		ShDegree:       uint8(shDegree),
		FractionalBits: 12,
		Flags:          0,
		Reserved:       0,
	}

	bts := make([]byte, 0)
	bts = append(bts, h.ToBytes()...)

	for i := range rows {
		bts = append(bts, cmn.SpzEncodePosition(rows[i].PositionX)...)
		bts = append(bts, cmn.SpzEncodePosition(rows[i].PositionY)...)
		bts = append(bts, cmn.SpzEncodePosition(rows[i].PositionZ)...)
	}
	for i := range rows {
		bts = append(bts, rows[i].ColorA)
	}
	for i := range rows {
		bts = append(bts, cmn.SpzEncodeColor(rows[i].ColorR), cmn.SpzEncodeColor(rows[i].ColorG), cmn.SpzEncodeColor(rows[i].ColorB))
	}
	for i := range rows {
		bts = append(bts, cmn.SpzEncodeScale(rows[i].ScaleX), cmn.SpzEncodeScale(rows[i].ScaleY), cmn.SpzEncodeScale(rows[i].ScaleZ))
	}
	for i := range rows {
		bts = append(bts, cmn.SpzEncodeRotations(rows[i].RotationW, rows[i].RotationX, rows[i].RotationY, rows[i].RotationZ)...)
	}

	if shDegree == 1 {
		for i := range rows {
			if len(rows[i].SH1) > 0 {
				for n := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH1[n]))
				}
			} else if len(rows[i].SH2) > 0 {
				for n := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH2[n]))
				}
			} else {
				for range 9 {
					bts = append(bts, cmn.EncodeSplatSH(0.0))
				}
			}
		}
	} else if shDegree == 2 {
		for i := range rows {
			if len(rows[i].SH1) > 0 {
				for n := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH1[n]))
				}
				for range 15 {
					bts = append(bts, cmn.EncodeSplatSH(0.0))
				}
			} else if len(rows[i].SH2) > 0 {
				for n := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH2[n]))
				}
				for j := 9; j < 24; j++ {
					bts = append(bts, cmn.SpzEncodeSH23(rows[i].SH2[j]))
				}
			} else {
				for range 24 {
					bts = append(bts, cmn.EncodeSplatSH(0.0))
				}
			}
		}
	} else if shDegree == 3 {
		for i := range rows {
			if len(rows[i].SH3) > 0 {
				for j := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH2[j]))
				}
				for j := 9; j < 24; j++ {
					bts = append(bts, cmn.SpzEncodeSH23(rows[i].SH2[j]))
				}
				for j := range 21 {
					bts = append(bts, cmn.SpzEncodeSH23(rows[i].SH3[j]))
				}
			} else if len(rows[i].SH2) > 0 {
				for j := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH2[j]))
				}
				for j := 9; j < 24; j++ {
					bts = append(bts, cmn.SpzEncodeSH23(rows[i].SH2[j]))
				}
				for range 21 {
					bts = append(bts, cmn.EncodeSplatSH(0.0))
				}
			} else if len(rows[i].SH1) > 0 {
				for n := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH1[n]))
				}
				for range 36 {
					bts = append(bts, cmn.EncodeSplatSH(0.0))
				}
			} else {
				for range 45 {
					bts = append(bts, cmn.EncodeSplatSH(0.0))
				}
			}
		}
	}

	gzipDatas, err := cmn.GzipBytes(bts)
	cmn.ExitOnError(err)

	_, err = writer.Write(gzipDatas)
	cmn.ExitOnError(err)

	err = writer.Flush()
	cmn.ExitOnError(err)
}
