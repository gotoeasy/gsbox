package gsplat

import (
	"bufio"
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
