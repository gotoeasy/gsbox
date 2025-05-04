package gsplat

import (
	"errors"
	"gsbox/cmn"
	"log"
	"os"
	"strings"
)

type PlyHeader struct {
	Declare      string
	Format       string
	Comment      string
	VertexCount  int
	HeaderLength int
	RowLength    int
	text         string
	mapOffset    map[string]int
	mapType      map[string]string
}

func ReadPlyHeaderString(plyFile string, readLen int) (string, error) {
	file, err := os.Open(plyFile)
	cmn.ExitOnError(err)
	defer file.Close()

	bs := make([]byte, readLen)
	_, err = file.Read(bs)
	if err != nil {
		return "", err
	}

	str := cmn.BytesToString(bs)
	if !cmn.Contains(str, "end_header\n") {
		if readLen > 1024*10 {
			return "", errors.New("ply header not found")
		}
		return ReadPlyHeaderString(plyFile, readLen+1024)
	}

	header := cmn.Split(str, "end_header\n")[0] + "end_header\n"
	return header, nil
}

func getPlyHeader(file *os.File, readLen int) (*PlyHeader, error) {
	bs := make([]byte, readLen)
	_, err := file.Read(bs)
	if err != nil {
		return nil, err
	}

	str := cmn.BytesToString(bs)
	if !cmn.Contains(str, "end_header\n") {
		if readLen > 1024*10 {
			return nil, errors.New("ply header not found")
		}
		return getPlyHeader(file, readLen+1024)
	}

	mapOffset := make(map[string]int)
	mapType := make(map[string]string)

	header := cmn.Split(str, "end_header\n")[0] + "end_header\n"
	lines := cmn.Split(strings.TrimRight(header, "\n"), "\n")
	vertexCount := -1
	declare := ""
	format := ""
	comment := ""
	offset := 0
	for i := 0; i < len(lines); i++ {
		if i == 0 {
			declare = lines[i]
		} else if cmn.Startwiths(lines[i], "property ") {
			ary := cmn.Split(lines[i], " ")
			mapOffset[ary[2]] = offset
			mapType[ary[2]] = ary[1]
			offset += getTypeSize(ary[1])
		} else if cmn.Startwiths(lines[i], "format ") {
			format = cmn.ReplaceAll(lines[i], "format ", "")
		} else if cmn.Startwiths(lines[i], "element vertex ") {
			vertexCount = cmn.StringToInt(cmn.ReplaceAll(lines[i], "element vertex ", ""))
		} else if cmn.Startwiths(lines[i], "comment ") {
			comment = cmn.ReplaceAll(lines[i], "comment ", "")
		}
	}

	plyHeader := &PlyHeader{
		text:         header,
		Declare:      declare,
		Format:       format,
		Comment:      comment,
		VertexCount:  vertexCount,
		HeaderLength: len(header),
		RowLength:    offset,
		mapType:      mapType,
		mapOffset:    mapOffset,
	}
	return plyHeader, nil
}

func getTypeSize(name string) int {
	if name == "float" {
		return 4
	} else if name == "double" {
		return 8
	} else if name == "int" {
		return 4
	} else if name == "uint" {
		return 4
	} else if name == "short" {
		return 2
	} else if name == "ushort" {
		return 2
	} else if name == "uchar" {
		return 1
	}

	log.Println("[Error] unsupported property type:", name)
	os.Exit(1)
	return 0 // unknown
}

func (p *PlyHeader) Property(property string) (int, string) {
	return p.mapOffset[property], p.mapType[property]
}

func (p *PlyHeader) MaxShDegree() int {
	if p.mapType["f_rest_44"] != "" {
		return 3
	} else if p.mapType["f_rest_23"] != "" {
		return 2
	} else if p.mapType["f_rest_8"] != "" {
		return 1
	}
	return 0
}

func (p *PlyHeader) IsPly() bool {
	return cmn.Startwiths(p.text, "ply\n")
}

func (p *PlyHeader) GetComment() string {
	return p.Comment
}

func (p *PlyHeader) GetFormat() string {
	return p.Format
}

func (p *PlyHeader) ToString() string {
	return p.text
}

// 是否3dgs官方格式ply
func (p *PlyHeader) IsOfficialPly() bool {
	return p.Declare == "ply" &&
		p.Format == "binary_little_endian 1.0" &&
		p.VertexCount > 0 &&
		p.mapType["x"] != "" && p.mapType["y"] != "" && p.mapType["z"] != "" &&
		p.mapType["f_dc_0"] != "" && p.mapType["f_dc_1"] != "" && p.mapType["f_dc_2"] != "" &&
		p.mapType["opacity"] != "" && p.mapType["scale_0"] != "" && p.mapType["scale_1"] != "" && p.mapType["scale_2"] != "" &&
		p.mapType["rot_0"] != "" && p.mapType["rot_1"] != "" && p.mapType["rot_2"] != "" && p.mapType["rot_3"] != ""
}
