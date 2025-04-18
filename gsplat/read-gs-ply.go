package gsplat

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"gsbox/cmn"
	"math"
	"os"
	"sort"
)

const SH_C0 float64 = 0.28209479177387814

func ReadPlyHeader(plyFile string) (*PlyHeader, error) {
	file, err := os.Open(plyFile)
	cmn.ExitOnError(err)
	defer file.Close()
	return getPlyHeader(file, 2048)
}

func ReadPly(plyFile string, shDegree int, plyTypes ...string) []*SplatData {
	file, err := os.Open(plyFile)
	cmn.ExitOnError(err)
	defer file.Close()

	header, err := getPlyHeader(file, 2048)
	cmn.ExitOnError(err)

	if len(plyTypes) > 0 && cmn.EqualsIngoreCase(plyTypes[0], "ply-3dgs") {
		if !header.IsOfficialPly() {
			cmn.ExitOnError(errors.New("unsupported ply file: " + plyFile))
		}
	}

	datas := make([]*SplatData, header.VertexCount)
	for i := 0; i < header.VertexCount; i++ {
		dataBytes := make([]byte, header.RowLength)
		_, err := file.ReadAt(dataBytes, int64(header.HeaderLength+i*header.RowLength))
		cmn.ExitOnError(err)

		data := &SplatData{}
		data.PositionX = cmn.ClipFloat32(readValue(header, "x", dataBytes))
		data.PositionY = cmn.ClipFloat32(readValue(header, "y", dataBytes))
		data.PositionZ = cmn.ClipFloat32(readValue(header, "z", dataBytes))
		data.ScaleX = cmn.ClipFloat32(math.Exp(readValue(header, "scale_0", dataBytes)))
		data.ScaleY = cmn.ClipFloat32(math.Exp(readValue(header, "scale_1", dataBytes)))
		data.ScaleZ = cmn.ClipFloat32(math.Exp(readValue(header, "scale_2", dataBytes)))
		data.ColorR = cmn.ClipUint8((0.5 + SH_C0*readValue(header, "f_dc_0", dataBytes)) * 255.0)
		data.ColorG = cmn.ClipUint8((0.5 + SH_C0*readValue(header, "f_dc_1", dataBytes)) * 255.0)
		data.ColorB = cmn.ClipUint8((0.5 + SH_C0*readValue(header, "f_dc_2", dataBytes)) * 255.0)
		data.ColorA = cmn.ClipUint8((1.0 / (1.0 + math.Exp(-readValue(header, "opacity", dataBytes)))) * 255.0)

		r0, r1, r2, r3 := readValue(header, "rot_0", dataBytes), readValue(header, "rot_1", dataBytes), readValue(header, "rot_2", dataBytes), readValue(header, "rot_3", dataBytes)
		qlen := math.Sqrt(r0*r0 + r1*r1 + r2*r2 + r3*r3)
		data.RotationX = cmn.ClipUint8((r0/qlen)*128.0 + 128.0)
		data.RotationY = cmn.ClipUint8((r1/qlen)*128.0 + 128.0)
		data.RotationZ = cmn.ClipUint8((r2/qlen)*128.0 + 128.0)
		data.RotationW = cmn.ClipUint8((r3/qlen)*128.0 + 128.0)

		datas[i] = data

		if shDegree == 1 {
			for n := range 9 {
				data.SH1 = append(data.SH1, cmn.SpzEncodeSH1(readValue(header, "f_rest_"+cmn.IntToString(n), dataBytes)))
			}
		} else if shDegree == 2 {
			for n := range 24 {
				data.SH2 = append(data.SH2, cmn.SpzEncodeSH23(readValue(header, "f_rest_"+cmn.IntToString(n), dataBytes)))
			}
		} else if shDegree == 3 {
			for n := range 24 {
				data.SH2 = append(data.SH2, cmn.SpzEncodeSH23(readValue(header, "f_rest_"+cmn.IntToString(n), dataBytes)))
			}
			for n := 24; n < 45; n++ {
				data.SH3 = append(data.SH3, cmn.SpzEncodeSH23(readValue(header, "f_rest_"+cmn.IntToString(n), dataBytes)))
			}
		}

	}

	return datas
}

func readValue(header *PlyHeader, property string, splatDataBytes []byte) float64 {
	offset, typename := header.Property(property)
	if typename == "float" {
		var v float32
		cmn.ExitOnError(binary.Read(bytes.NewReader(splatDataBytes[offset:offset+4]), binary.LittleEndian, &v))
		return float64(v)
	} else if typename == "double" {
		var v float64
		cmn.ExitOnError(binary.Read(bytes.NewReader(splatDataBytes[offset:offset+8]), binary.LittleEndian, &v))
		return v
	} else if typename == "int" {
		var v int32
		cmn.ExitOnError(binary.Read(bytes.NewReader(splatDataBytes[offset:offset+4]), binary.LittleEndian, &v))
		return float64(v)
	} else if typename == "uint" {
		var v uint32
		cmn.ExitOnError(binary.Read(bytes.NewReader(splatDataBytes[offset:offset+4]), binary.LittleEndian, &v))
		return float64(v)
	} else if typename == "short" {
		var v int16
		cmn.ExitOnError(binary.Read(bytes.NewReader(splatDataBytes[offset:offset+2]), binary.LittleEndian, &v))
		return float64(v)
	} else if typename == "ushort" {
		var v uint16
		cmn.ExitOnError(binary.Read(bytes.NewReader(splatDataBytes[offset:offset+2]), binary.LittleEndian, &v))
		return float64(v)
	} else if typename == "uchar" {
		var v uint8
		cmn.ExitOnError(binary.Read(bytes.NewReader(splatDataBytes[offset:offset+1]), binary.LittleEndian, &v))
		return float64(v)
	}

	if cmn.Startwiths(property, "f_rest_") {
		return 0 // 球谐系数读取不到时，默认为0
	}

	fmt.Println("Unsupported property:", "property", typename, property)
	os.Exit(1)
	return 0
}

// 排序
func Sort(rows []*SplatData) {
	sort.Slice(rows, func(i, j int) bool {
		return math.Exp(float64(rows[i].ScaleX+rows[i].ScaleY+rows[i].ScaleZ))*float64(rows[i].ColorA) <
			math.Exp(float64(rows[j].ScaleX+rows[j].ScaleY+rows[j].ScaleZ))*float64(rows[j].ColorA)
	})
}
