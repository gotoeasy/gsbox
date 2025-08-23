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
}

func PrintSplat(splatFile string, rows []*SplatData) {
	file, err := os.Create(splatFile)
	cmn.ExitOnError(err)
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i := range rows {
		line := fmt.Sprintf("%s,%s,%s, %s,%s,%s, %v,%v,%v,%v, %v,%v,%v,%v\r\n",
			cmn.FormatFloat32(rows[i].PositionX), cmn.FormatFloat32(rows[i].PositionY), cmn.FormatFloat32(rows[i].PositionZ),
			cmn.FormatFloat32(cmn.EncodeSplatScale(rows[i].ScaleX)), cmn.FormatFloat32(cmn.EncodeSplatScale(rows[i].ScaleY)), cmn.FormatFloat32(cmn.EncodeSplatScale(rows[i].ScaleZ)),
			rows[i].ColorR, rows[i].ColorG, rows[i].ColorB, rows[i].ColorA,
			rows[i].RotationW, rows[i].RotationX, rows[i].RotationY, rows[i].RotationZ,
		)
		_, err = writer.WriteString(line)
		cmn.ExitOnError(err)
	}
	err = writer.Flush()
	cmn.ExitOnError(err)
}
