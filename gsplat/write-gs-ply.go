package gsplat

import (
	"bufio"
	"gsbox/cmn"
	"math"
	"os"
)

func WritePly(plyFile string, datas []*SplatData, comment string, simplePly bool) {
	file, err := os.Create(plyFile)
	cmn.ExitOnError(err)
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(genPlyHeader(len(datas), comment, simplePly))
	cmn.ExitOnError(err)
	for i := 0; i < len(datas); i++ {
		_, err = writer.Write(genPlyDataBin(datas[i], simplePly))
		cmn.ExitOnError(err)
	}
	err = writer.Flush()
	cmn.ExitOnError(err)
}

func genPlyDataBin(splatData *SplatData, simplePly bool) []byte {

	bts := []byte{}
	bts = append(bts, cmn.Float32ToBytes(splatData.PositionX)...) // x
	bts = append(bts, cmn.Float32ToBytes(splatData.PositionY)...) // y
	bts = append(bts, cmn.Float32ToBytes(splatData.PositionZ)...) // z
	if !simplePly {
		bts = append(bts, make([]byte, 3*4)...) // nx, ny, nz
	}
	bts = append(bts, cmn.ToFloat32Bytes((float64(splatData.ColorR)/255-0.5)/SH_C0)...) // f_dc_0
	bts = append(bts, cmn.ToFloat32Bytes((float64(splatData.ColorG)/255-0.5)/SH_C0)...) // f_dc_1
	bts = append(bts, cmn.ToFloat32Bytes((float64(splatData.ColorB)/255-0.5)/SH_C0)...) // f_dc_2
	if !simplePly {
		bts = append(bts, make([]byte, 45*4)...) // f_rest_0 ~ f_rest_44
	}
	bts = append(bts, cmn.ToFloat32Bytes(-math.Log((1/(float64(splatData.ColorA)/255))-1))...) // opacity
	bts = append(bts, cmn.ToFloat32Bytes(math.Log(float64(splatData.ScaleX)))...)              // scale_0
	bts = append(bts, cmn.ToFloat32Bytes(math.Log(float64(splatData.ScaleY)))...)              // scale_1
	bts = append(bts, cmn.ToFloat32Bytes(math.Log(float64(splatData.ScaleZ)))...)              // scale_2
	bts = append(bts, cmn.ToFloat32Bytes((float64(splatData.RotationX)-128)/128)...)           // rot_0
	bts = append(bts, cmn.ToFloat32Bytes((float64(splatData.RotationY)-128)/128)...)           // rot_1
	bts = append(bts, cmn.ToFloat32Bytes((float64(splatData.RotationZ)-128)/128)...)           // rot_2
	bts = append(bts, cmn.ToFloat32Bytes((float64(splatData.RotationW)-128)/128)...)           // rot_3

	return bts
}

func genPlyHeader(count int, comment string, simplePly bool) string {
	lines := []string{}
	lines = append(lines, "ply")
	lines = append(lines, "format binary_little_endian 1.0")
	lines = append(lines, "element vertex "+cmn.IntToString(count))
	if comment != "" {
		comment = cmn.ReplaceAll(comment, "\r", "\\r")
		comment = cmn.ReplaceAll(comment, "\n", "\\n")
		lines = append(lines, "comment "+comment)
	}
	lines = append(lines, "property float x")
	lines = append(lines, "property float y")
	lines = append(lines, "property float z")
	if !simplePly {
		lines = append(lines, "property float nx")
		lines = append(lines, "property float ny")
		lines = append(lines, "property float nz")
	}
	lines = append(lines, "property float f_dc_0")
	lines = append(lines, "property float f_dc_1")
	lines = append(lines, "property float f_dc_2")
	if !simplePly {
		for i := 0; i <= 44; i++ {
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
