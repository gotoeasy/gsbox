package gsplat

import (
	"bufio"
	"gsbox/cmn"
	"os"
)

func WriteSpz(spzFile string, rows []*SplatData, shDegree int) {
	file, err := os.Create(spzFile)
	cmn.ExitOnError(err)
	defer file.Close()

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
		if len(rows[0].SH1) > 0 {
			for i := range rows {
				for n := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH1[n]))
				}
			}
		} else if len(rows[0].SH2) > 0 {
			for i := range rows {
				for n := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH2[n]))
				}
			}
		} else {
			bts = append(bts, make([]byte, h.NumPoints*9)...)
		}
	} else if shDegree == 2 {
		if len(rows[0].SH1) > 0 {
			for i := range rows {
				for n := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH1[n]))
				}
				bts = append(bts, make([]byte, 15)...)
			}
		} else if len(rows[0].SH2) > 0 {
			for i := range rows {
				for n := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH2[n]))
				}
				for j := 9; j < 24; j++ {
					bts = append(bts, cmn.SpzEncodeSH23(rows[i].SH2[j]))
				}
			}
		} else {
			bts = append(bts, make([]byte, h.NumPoints*24)...)
		}
	} else if shDegree == 3 {
		if len(rows[0].SH3) > 0 {
			for i := range rows {
				for j := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH2[j]))
				}
				for j := 9; j < 24; j++ {
					bts = append(bts, cmn.SpzEncodeSH23(rows[i].SH2[j]))
				}
				for j := range 21 {
					bts = append(bts, cmn.SpzEncodeSH23(rows[i].SH3[j]))
				}
			}
		} else if len(rows[0].SH2) > 0 {
			for i := range rows {
				for j := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH2[j]))
				}
				for j := 9; j < 24; j++ {
					bts = append(bts, cmn.SpzEncodeSH23(rows[i].SH2[j]))
				}
				bts = append(bts, make([]byte, 21)...)
			}
		} else if len(rows[0].SH1) > 0 {
			for i := range rows {
				for n := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH1[n]))
				}
				bts = append(bts, make([]byte, 36)...)
			}
		} else {
			bts = append(bts, make([]byte, h.NumPoints*45)...)
		}
	}

	gzipDatas, err := cmn.GzipBytes(bts)
	cmn.ExitOnError(err)

	_, err = writer.Write(gzipDatas)
	cmn.ExitOnError(err)

	err = writer.Flush()
	cmn.ExitOnError(err)
}
