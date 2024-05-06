package gsplat

import (
	"errors"
	"fmt"
	"gsbox/cmn"
	"os"
	"strings"
)

type PlyHeader struct {
	Format       string
	Comment      string
	VertexCount  int
	HeaderLength int
	RowLength    int
	text         string
	mapOffset    map[string]int
	mapType      map[string]string
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
	vertexCount := 0
	format := ""
	comment := ""
	offset := 0
	for i := 0; i < len(lines); i++ {
		if cmn.Startwiths(lines[i], "property ") {
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

	fmt.Println("Unsupported property type:", name)
	os.Exit(1)
	return 0 // unknown
}

func (p *PlyHeader) Property(property string) (int, string) {
	return p.mapOffset[property], p.mapType[property]
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
