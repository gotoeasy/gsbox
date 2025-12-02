package gsplat

import (
	"bufio"
	"fmt"
	"gsbox/cmn"
	"os"
)

func WriteSplat(splatFile string, rows []*SplatData) {
	file, err := os.Create(splatFile)
	cmn.ExitOnError(err)
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i := range rows {
		OnProgress(PhaseWrite, i, len(rows))
		bts := make([]byte, 0)
		bts = append(bts, cmn.Float32ToBytes(rows[i].PositionX)...)
		bts = append(bts, cmn.Float32ToBytes(rows[i].PositionY)...)
		bts = append(bts, cmn.Float32ToBytes(rows[i].PositionZ)...)
		bts = append(bts, cmn.Float32ToBytes(cmn.EncodeSplatScale(rows[i].ScaleX))...)
		bts = append(bts, cmn.Float32ToBytes(cmn.EncodeSplatScale(rows[i].ScaleY))...)
		bts = append(bts, cmn.Float32ToBytes(cmn.EncodeSplatScale(rows[i].ScaleZ))...)
		bts = append(bts, rows[i].ColorR, rows[i].ColorG, rows[i].ColorB, rows[i].ColorA)
		bts = append(bts, rows[i].RotationW, rows[i].RotationX, rows[i].RotationY, rows[i].RotationZ)
		_, err = writer.Write(bts)
		cmn.ExitOnError(err)
	}
	err = writer.Flush()
	cmn.ExitOnError(err)
	OnProgress(PhaseWrite, 100, 100)
}

func PrintSplat(splatFile string, rows []*SplatData) {
	file, err := os.Create(splatFile)
	cmn.ExitOnError(err)
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i := range rows {
		if len(rows[i].SH45) == 0 {
			rows[i].SH45 = InitZeroSH45()
		}
		shs := rows[i].SH45

		shDegree := GetArgShDegree()
		switch shDegree {
		case 1:
			line := fmt.Sprintf("%s,%s,%s, %s,%s,%s, %v,%v,%v,%v, %v,%v,%v,%v, %v,%v,%v,%v,%v,%v,%v,%v,%v\r\n",
				cmn.FormatFloat32(rows[i].PositionX), cmn.FormatFloat32(rows[i].PositionY), cmn.FormatFloat32(rows[i].PositionZ),
				cmn.FormatFloat32(cmn.EncodeSplatScale(rows[i].ScaleX)), cmn.FormatFloat32(cmn.EncodeSplatScale(rows[i].ScaleY)), cmn.FormatFloat32(cmn.EncodeSplatScale(rows[i].ScaleZ)),
				rows[i].ColorR, rows[i].ColorG, rows[i].ColorB, rows[i].ColorA,
				rows[i].RotationW, rows[i].RotationX, rows[i].RotationY, rows[i].RotationZ,
				shs[0], shs[1], shs[2], shs[3], shs[4], shs[5], shs[6], shs[7], shs[8],
			)
			_, err = writer.WriteString(line)
			cmn.ExitOnError(err)
		case 2:
			line := fmt.Sprintf("%s,%s,%s, %s,%s,%s, %v,%v,%v,%v, %v,%v,%v,%v, %v,%v,%v,%v,%v,%v,%v,%v,%v, %v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\r\n",
				cmn.FormatFloat32(rows[i].PositionX), cmn.FormatFloat32(rows[i].PositionY), cmn.FormatFloat32(rows[i].PositionZ),
				cmn.FormatFloat32(cmn.EncodeSplatScale(rows[i].ScaleX)), cmn.FormatFloat32(cmn.EncodeSplatScale(rows[i].ScaleY)), cmn.FormatFloat32(cmn.EncodeSplatScale(rows[i].ScaleZ)),
				rows[i].ColorR, rows[i].ColorG, rows[i].ColorB, rows[i].ColorA,
				rows[i].RotationW, rows[i].RotationX, rows[i].RotationY, rows[i].RotationZ,
				shs[0], shs[1], shs[2], shs[3], shs[4], shs[5], shs[6], shs[7], shs[8],
				shs[9], shs[10], shs[11], shs[12], shs[13], shs[14], shs[15], shs[16], shs[17], shs[18], shs[19], shs[20], shs[21], shs[22], shs[23],
			)
			_, err = writer.WriteString(line)
			cmn.ExitOnError(err)
		case 3:
			line := fmt.Sprintf("%s,%s,%s, %s,%s,%s, %v,%v,%v,%v, %v,%v,%v,%v, %v,%v,%v,%v,%v,%v,%v,%v,%v, %v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v, %v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\r\n",
				cmn.FormatFloat32(rows[i].PositionX), cmn.FormatFloat32(rows[i].PositionY), cmn.FormatFloat32(rows[i].PositionZ),
				cmn.FormatFloat32(cmn.EncodeSplatScale(rows[i].ScaleX)), cmn.FormatFloat32(cmn.EncodeSplatScale(rows[i].ScaleY)), cmn.FormatFloat32(cmn.EncodeSplatScale(rows[i].ScaleZ)),
				rows[i].ColorR, rows[i].ColorG, rows[i].ColorB, rows[i].ColorA,
				rows[i].RotationW, rows[i].RotationX, rows[i].RotationY, rows[i].RotationZ,
				shs[0], shs[1], shs[2], shs[3], shs[4], shs[5], shs[6], shs[7], shs[8],
				shs[9], shs[10], shs[11], shs[12], shs[13], shs[14], shs[15], shs[16], shs[17], shs[18], shs[19], shs[20], shs[21], shs[22], shs[23],
				shs[24], shs[25], shs[26], shs[27], shs[28], shs[29], shs[30], shs[31], shs[32], shs[33], shs[34], shs[35], shs[36], shs[37], shs[38], shs[39], shs[40], shs[41], shs[42], shs[43], shs[44],
			)
			_, err = writer.WriteString(line)
			cmn.ExitOnError(err)
		default:
			line := fmt.Sprintf("%s,%s,%s, %s,%s,%s, %v,%v,%v,%v, %v,%v,%v,%v\r\n",
				cmn.FormatFloat32(rows[i].PositionX), cmn.FormatFloat32(rows[i].PositionY), cmn.FormatFloat32(rows[i].PositionZ),
				cmn.FormatFloat32(cmn.EncodeSplatScale(rows[i].ScaleX)), cmn.FormatFloat32(cmn.EncodeSplatScale(rows[i].ScaleY)), cmn.FormatFloat32(cmn.EncodeSplatScale(rows[i].ScaleZ)),
				rows[i].ColorR, rows[i].ColorG, rows[i].ColorB, rows[i].ColorA,
				rows[i].RotationW, rows[i].RotationX, rows[i].RotationY, rows[i].RotationZ,
			)
			_, err = writer.WriteString(line)
			cmn.ExitOnError(err)
		}
	}
	err = writer.Flush()
	cmn.ExitOnError(err)
}
