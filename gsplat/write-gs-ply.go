package gsplat

import (
	"bufio"
	"gsbox/cmn"
	"log"
	"os"
)

func WritePly(plyFile string, datas []*SplatData) {
	file, err := os.Create(plyFile)
	cmn.ExitOnError(err)
	defer file.Close()

	comment := Args.GetArgIgnorecase("-c", "--comment")
	shDegree := GetArgShDegree()
	log.Println("[Info] output shDegree:", shDegree)
	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(genPlyHeader(len(datas), comment, shDegree))
	cmn.ExitOnError(err)
	for i := 0; i < len(datas); i++ {
		_, err = writer.Write(genPlyDataBin(datas[i], shDegree))
		cmn.ExitOnError(err)
	}
	err = writer.Flush()
	cmn.ExitOnError(err)
}

func genPlyDataBin(splatData *SplatData, shDegree uint8) []byte {

	bts := []byte{}
	bts = append(bts, cmn.Float32ToBytes(splatData.PositionX)...)                    // x
	bts = append(bts, cmn.Float32ToBytes(splatData.PositionY)...)                    // y
	bts = append(bts, cmn.Float32ToBytes(splatData.PositionZ)...)                    // z
	bts = append(bts, make([]byte, 3*4)...)                                          // nx, ny, nz
	bts = append(bts, cmn.Float32ToBytes(cmn.DecodeSplatColor(splatData.ColorR))...) // f_dc_0
	bts = append(bts, cmn.Float32ToBytes(cmn.DecodeSplatColor(splatData.ColorG))...) // f_dc_1
	bts = append(bts, cmn.Float32ToBytes(cmn.DecodeSplatColor(splatData.ColorB))...) // f_dc_2

	if shDegree > 0 {
		shDim := 0
		var shs []byte
		switch shDegree {
		case 1:
			shDim = 3
			if len(splatData.SH1) > 0 {
				shs = append(shs, splatData.SH1...)
			} else if len(splatData.SH2) > 0 {
				shs = append(shs, splatData.SH2[:9]...)
			}
		case 2:
			shDim = 8
			if len(splatData.SH1) > 0 {
				shs = append(shs, splatData.SH1...)
			} else if len(splatData.SH2) > 0 {
				shs = append(shs, splatData.SH2...)
			}
		case 3:
			shDim = 15
			if len(splatData.SH3) > 0 {
				shs = append(shs, splatData.SH2...)
				shs = append(shs, splatData.SH3...)
			} else if len(splatData.SH2) > 0 {
				shs = append(shs, splatData.SH2...)
			} else if len(splatData.SH1) > 0 {
				shs = append(shs, splatData.SH1...)
			}
		}

		for i := len(shs); i < 45; i++ {
			shs = append(shs, cmn.EncodeSplatSH(0.0))
		}

		for c := range 3 {
			for i := range shDim {
				bts = append(bts, cmn.Float32ToBytes(cmn.DecodeSplatSH(shs[c+i*3]))...) // f_rest_0 ... f_rest_n
			}
		}
	}

	bts = append(bts, cmn.Float32ToBytes(cmn.DecodeSplatOpacity(splatData.ColorA))...)     // opacity
	bts = append(bts, cmn.Float32ToBytes(splatData.ScaleX)...)                             // scale_0
	bts = append(bts, cmn.Float32ToBytes(splatData.ScaleY)...)                             // scale_1
	bts = append(bts, cmn.Float32ToBytes(splatData.ScaleZ)...)                             // scale_2
	bts = append(bts, cmn.Float32ToBytes(cmn.DecodeSplatRotation(splatData.RotationW))...) // rot_0
	bts = append(bts, cmn.Float32ToBytes(cmn.DecodeSplatRotation(splatData.RotationX))...) // rot_1
	bts = append(bts, cmn.Float32ToBytes(cmn.DecodeSplatRotation(splatData.RotationY))...) // rot_2
	bts = append(bts, cmn.Float32ToBytes(cmn.DecodeSplatRotation(splatData.RotationZ))...) // rot_3

	return bts
}

func genPlyHeader(count int, comment string, shDegree uint8) string {
	lines := []string{}
	lines = append(lines, "ply")
	lines = append(lines, "format binary_little_endian 1.0")
	lines = append(lines, "element vertex "+cmn.IntToString(count))
	if comment != "" {
		_, comment = cmn.RemoveNonASCII(comment)
		lines = append(lines, "comment "+comment)
	}
	lines = append(lines, "property float x")
	lines = append(lines, "property float y")
	lines = append(lines, "property float z")
	lines = append(lines, "property float nx")
	lines = append(lines, "property float ny")
	lines = append(lines, "property float nz")
	lines = append(lines, "property float f_dc_0")
	lines = append(lines, "property float f_dc_1")
	lines = append(lines, "property float f_dc_2")
	switch shDegree {
	case 1:
		for i := range 9 {
			lines = append(lines, "property float f_rest_"+cmn.IntToString(i))
		}
	case 2:
		for i := range 24 {
			lines = append(lines, "property float f_rest_"+cmn.IntToString(i))
		}
	case 3:
		for i := range 45 {
			lines = append(lines, "property float f_rest_"+cmn.IntToString(i))
		}
	}
	lines = append(lines, "property float opacity")
	lines = append(lines, "property float scale_0")
	lines = append(lines, "property float scale_1")
	lines = append(lines, "property float scale_2")
	lines = append(lines, "property float rot_0")
	lines = append(lines, "property float rot_1")
	lines = append(lines, "property float rot_2")
	lines = append(lines, "property float rot_3")
	lines = append(lines, "end_header")
	return cmn.Join(lines, "\n") + "\n"
}
