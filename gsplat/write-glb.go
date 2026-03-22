package gsplat

import (
	"bufio"
	"gsbox/cmn"
	"log"
	"os"
)

const KHR_gaussian_splatting = "KHR_gaussian_splatting"                                     // 可被降级渲染，当前为候选阶段
const KHR_gaussian_splatting_compression_spz_2 = "KHR_gaussian_splatting_compression_spz_2" // 内嵌spz，当前为提案阶段

func WriteGlb(glbFile string, rows []*SplatData) int64 {
	fmt := Args.GetArgIgnorecase("-of", "--output-format")
	if fmt == "" {
		fmt = KHR_gaussian_splatting
	} else if cmn.Contains(fmt, "spz") {
		fmt = KHR_gaussian_splatting_compression_spz_2
	} else {
		fmt = KHR_gaussian_splatting
	}

	log.Println("[Info] (parameter) of:", fmt, "(glTF extension)")

	if fmt == KHR_gaussian_splatting_compression_spz_2 {
		return writeGlb_genGlbJson_KHR_gaussian_splatting_compression_spz_2(glbFile, rows)
	}

	return writeGlb_KHR_gaussian_splatting(glbFile, rows)
}

func writeGlb_KHR_gaussian_splatting(glbFile string, rows []*SplatData) int64 {
	file, err := os.Create(glbFile)
	cmn.ExitOnError(err)
	defer file.Close()
	writer := bufio.NewWriter(file)

	strJson := genJson_KHR_gaussian_splatting(len(rows), ComputeXyzMinMax(rows))
	outputShDegree := GetArgShDegree()
	log.Println("[Info] output shDegree:", outputShDegree)

	binLength := len(rows) * (12 + 4 + 16 + 12 + 4 + 12)
	if outputShDegree > 0 {
		binLength += len(rows) * (12 * 3)
	}
	if outputShDegree > 1 {
		binLength += len(rows) * (12 * 5)
	}
	if outputShDegree > 2 {
		binLength += len(rows) * (12 * 7)
	}

	jsonLength := len(strJson)
	jsonLength4 := (jsonLength + 3) &^ 3                 // 4字节对齐
	binLength4 := (binLength + 3) &^ 3                   // 4字节对齐
	totalLength := 12 + 8 + jsonLength4 + 8 + binLength4 // 文件总长

	// Header
	writer.WriteString("glTF")                           // Magic
	writer.Write(cmn.Uint32ToBytes(2))                   // Version
	writer.Write(cmn.Uint32ToBytes(uint32(totalLength))) // Length

	// JSON Chunk
	writer.Write(cmn.Uint32ToBytes(uint32(jsonLength4))) // Length（对齐长度，对于json填充的空格不影响解析）
	writer.WriteString("JSON")                           // Type: 0x4E4F534A ("JSON")
	writer.Write(cmn.StringToBytes(strJson))             // JSON Data
	for range jsonLength4 - jsonLength {
		writer.WriteByte(' ') // 填充
	}

	// BIN Chunk
	writer.Write(cmn.Uint32ToBytes(uint32(binLength))) // Length（非对齐长度）
	writer.Write([]byte{'B', 'I', 'N', 0})             // Type: 0x004E4942 ("BIN")

	// BIN Data
	for _, d := range rows {
		writer.Write(cmn.Float32ToBytes(d.PositionX))
		writer.Write(cmn.Float32ToBytes(d.PositionY))
		writer.Write(cmn.Float32ToBytes(d.PositionZ))
	}
	for _, d := range rows {
		writer.WriteByte(d.ColorR)
		writer.WriteByte(d.ColorG)
		writer.WriteByte(d.ColorB)
		writer.WriteByte(d.ColorA)
	}
	for _, d := range rows {
		writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatRotation(d.RotationX)))
		writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatRotation(d.RotationY)))
		writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatRotation(d.RotationZ)))
		writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatRotation(d.RotationW)))
	}
	for _, d := range rows {
		writer.Write(cmn.Float32ToBytes(d.ScaleX))
		writer.Write(cmn.Float32ToBytes(d.ScaleY))
		writer.Write(cmn.Float32ToBytes(d.ScaleZ))
	}
	for _, d := range rows {
		writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatOpacity(d.ColorA)))
	}
	for _, d := range rows {
		writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatColor(d.ColorR)))
		writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatColor(d.ColorG)))
		writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatColor(d.ColorB)))
	}

	if outputShDegree > 0 {
		for _, d := range rows {
			if len(d.SH45) == 0 {
				d.SH45 = InitZeroSH45()
			}
		}

		for i := range 3 {
			for _, d := range rows {
				writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatSH(d.SH45[i*3+0])))
				writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatSH(d.SH45[i*3+1])))
				writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatSH(d.SH45[i*3+2])))
			}
		}
	}
	if outputShDegree > 1 {
		for i := range 5 {
			for _, d := range rows {
				writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatSH(d.SH45[(i+3)*3+0])))
				writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatSH(d.SH45[(i+3)*3+1])))
				writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatSH(d.SH45[(i+3)*3+2])))
			}
		}
	}
	if outputShDegree > 2 {
		for i := range 7 {
			for _, d := range rows {
				writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatSH(d.SH45[(i+8)*3+0])))
				writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatSH(d.SH45[(i+8)*3+1])))
				writer.Write(cmn.Float32ToBytes(cmn.DecodeSplatSH(d.SH45[(i+8)*3+2])))
			}
		}
	}

	// 填充
	for range binLength4 - binLength {
		writer.WriteByte(0)
	}

	err = writer.Flush()
	cmn.ExitOnError(err)
	return int64(totalLength)
}

func writeGlb_genGlbJson_KHR_gaussian_splatting_compression_spz_2(glbFile string, rows []*SplatData) int64 {
	file, err := os.Create(glbFile)
	cmn.ExitOnError(err)
	defer file.Close()

	writer := bufio.NewWriter(file)
	bts := genSpzVer3Bytes(rows)
	strJson := genJson_KHR_gaussian_splatting_compression_spz_2(len(bts))

	jsonLength := len(strJson)
	binLength := len(bts)
	jsonLength4 := (len(strJson) + 3) &^ 3               // 4字节对齐
	binLength4 := (binLength + 3) &^ 3                   // 4字节对齐
	totalLength := 12 + 8 + jsonLength4 + 8 + binLength4 // 文件总长

	// Header
	writer.WriteString("glTF")                           // Magic
	writer.Write(cmn.Uint32ToBytes(2))                   // Version
	writer.Write(cmn.Uint32ToBytes(uint32(totalLength))) // Length

	// JSON Chunk
	writer.Write(cmn.Uint32ToBytes(uint32(jsonLength4))) // Length（对齐长度）
	writer.WriteString("JSON")                           // Type: 0x4E4F534A ("JSON")
	writer.Write(cmn.StringToBytes(strJson))             // JSON Data
	for range jsonLength4 - jsonLength {
		writer.WriteByte(' ') // 填充
	}

	// BIN Chunk
	writer.Write(cmn.Uint32ToBytes(uint32(binLength))) // Length（非对齐长度）
	writer.Write([]byte{'B', 'I', 'N', 0})             // Type: 0x004E4942 ("BIN")
	writer.Write(bts)                                  // BIN Data
	for range binLength4 - binLength {
		writer.WriteByte(0) // 填充
	}

	err = writer.Flush()
	cmn.ExitOnError(err)
	return int64(totalLength)
}

func genSpzVer3Bytes(rows []*SplatData) []byte {

	outputShDegree := GetArgShDegree()
	ver := 3
	log.Println("[Info] output spz version:", ver)
	log.Println("[Info] output shDegree:", outputShDegree)

	h := &SpzHeader{
		Magic:          SPZ_MAGIC,
		Version:        uint32(ver),
		NumPoints:      uint32(len(rows)),
		ShDegree:       uint8(outputShDegree),
		FractionalBits: 12,
		Flags:          0,
		Reserved:       0,
	}

	if outputShDegree > 0 {
		log.Println("[Info] quality level:", oArg.Quality, "(range 1~9)")
		if oArg.Quality < 9 {
			// 小于9级使用聚类，第9级按通常处理
			_, _, paletteSize := ReWriteShByKmeans(rows)
			log.Println("[Info] sh palette size", paletteSize)
		}
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
		if h.Version >= 3 {
			bts = append(bts, cmn.SpzEncodeRotationsV3(rows[i].RotationW, rows[i].RotationX, rows[i].RotationY, rows[i].RotationZ)...)
		} else {
			bts = append(bts, cmn.SpzEncodeRotations(rows[i].RotationW, rows[i].RotationX, rows[i].RotationY, rows[i].RotationZ)...)
		}
	}

	switch outputShDegree {
	case 1:
		for i := range rows {
			if len(rows[i].SH45) > 0 {
				for j := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH45[j]))
				}
			} else {
				for range 9 {
					bts = append(bts, 128)
				}
			}
		}
	case 2:
		for i := range rows {
			if len(rows[i].SH45) > 0 {
				for j := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH45[j]))
				}
				for j := 9; j < 24; j++ {
					bts = append(bts, cmn.SpzEncodeSH23(rows[i].SH45[j]))
				}
			} else {
				for range 24 {
					bts = append(bts, 128)
				}
			}
		}
	case 3:
		for i := range rows {
			if len(rows[i].SH45) > 0 {
				for j := range 9 {
					bts = append(bts, cmn.SpzEncodeSH1(rows[i].SH45[j]))
				}
				for j := 9; j < 45; j++ {
					bts = append(bts, cmn.SpzEncodeSH23(rows[i].SH45[j]))
				}
			} else {
				for range 45 {
					bts = append(bts, 128)
				}
			}
		}
	}

	gzipDatas, err := cmn.CompressGzip(bts)
	cmn.ExitOnError(err)
	return gzipDatas
}

func genJson_KHR_gaussian_splatting_compression_spz_2(binLenght int) string {
	tmpl := `{
		"asset": {
			"generator": "gsbox $version",
			"version": "2.0"
		},
		"bufferViews": [
			{
			"buffer": 0,
			"byteLength": $byteLength
			}
		],
		"buffers": [
			{
			"byteLength": $byteLength
			}
		],
		"extensions": {},
		"extensionsRequired": [
			"KHR_gaussian_splatting",
			"KHR_gaussian_splatting_compression_spz_2"
		],
		"extensionsUsed": [
			"KHR_gaussian_splatting",
			"KHR_gaussian_splatting_compression_spz_2"
		],
		"meshes": [
			{
			"primitives": [
				{
				"attributes": {},
				"extensions": {
					"KHR_gaussian_splatting": {
					"extensions": {
						"KHR_gaussian_splatting_compression_spz_2": {
						"bufferView": 0
						}
					}
					}
				},
				"mode": 0
				}
			]
			}
		],
		"nodes": [
			{
			"mesh": 0
			}
		],
		"scene": 0,
		"scenes": [
			{
			"nodes": [
				0
			]
			}
		]
	}`

	tmpl = cmn.ReplaceAll(tmpl, "$version", cmn.VER)
	tmpl = cmn.ReplaceAll(tmpl, "$byteLength", cmn.IntToString(binLenght))
	return cmn.JsonStringify(tmpl)
}

func genJson_KHR_gaussian_splatting(count int, mm *V3MinMax) string {
	tmpl := `{
		"asset": {
			"generator": "gsbox $version",
			"version": "2.0"
		},
		"accessors": [
			{
			"bufferView": 0,
			"byteOffset": 0,
			"componentType": 5126,
			"count": $count,
			"min": [ $minX, $minY, $minZ ],
			"max": [ $maxX, $maxY, $maxZ ],
			"type": "VEC3"
			},
			{
			"bufferView": 1,
			"byteOffset": 0,
			"componentType": 5121,
			"count": $count,
			"normalized": true,
			"type": "VEC4"
			},
			{
			"bufferView": 2,
			"byteOffset": 0,
			"componentType": 5126,
			"count": $count,
			"type": "VEC4"
			},
			{
			"bufferView": 3,
			"byteOffset": 0,
			"componentType": 5126,
			"count": $count,
			"type": "VEC3"
			},
			{
			"bufferView": 4,
			"byteOffset": 0,
			"componentType": 5126,
			"count": $count,
			"type": "SCALAR"
			},
			{
			"bufferView": 5,
			"byteOffset": 0,
			"componentType": 5126,
			"count": $count,
			"type": "VEC3"
			}
			$accessorsSH$
		],
		"bufferViews": [
			{
			"buffer": 0,
			"byteLength": $byteLength0$,
			"byteOffset": $byteOffset0$,
			"target": 34962
			},
			{
			"buffer": 0,
			"byteLength": $byteLength1$,
			"byteOffset": $byteOffset1$,
			"target": 34962
			},
			{
			"buffer": 0,
			"byteLength": $byteLength2$,
			"byteOffset": $byteOffset2$,
			"target": 34962
			},
			{
			"buffer": 0,
			"byteLength": $byteLength3$,
			"byteOffset": $byteOffset3$,
			"target": 34962
			},
			{
			"buffer": 0,
			"byteLength": $byteLength4$,
			"byteOffset": $byteOffset4$,
			"target": 34962
			},
			{
			"buffer": 0,
			"byteLength": $byteLength5$,
			"byteOffset": $byteOffset5$,
			"target": 34962
			}
			$bufferViewsSH$
		],
		"buffers": [
			{
			"byteLength": $byteLengthTotal
			}
		],
		"extensionsUsed": [
			"KHR_gaussian_splatting"
		],
		"meshes": [
			{
			"primitives": [
				{
				"attributes": {
					"POSITION": 0,
					"COLOR_0": 1,
					"KHR_gaussian_splatting:ROTATION": 2,
					"KHR_gaussian_splatting:SCALE": 3,
					"KHR_gaussian_splatting:OPACITY": 4,
					"KHR_gaussian_splatting:SH_DEGREE_0_COEF_0": 5
					$SH_DEGREE$
				},
				"extensions": {
					"KHR_gaussian_splatting": {
						"colorSpace": "srgb_rec709_display",
						"kernel": "ellipse",
						"projection": "perspective",
						"sortingMethod": "cameraDistance"
					}
				},
				"mode": 0
				}
			]
			}
		],
		"nodes": [
			{
			"mesh": 0
			}
		],
		"scene": 0,
		"scenes": [
			{
			"nodes": [
				0
			]
			}
		]
	}`

	maxIndex := 6
	const TPL1 = `{"bufferView": $index, "byteOffset": 0, "componentType": 5126, "count": $count, "type": "VEC3"}`
	const TPL2 = `{"buffer": 0, "byteLength": $byteLength$index$, "byteOffset": $byteOffset$index$, "target": 34962}`
	const TPLSH1 = `,"KHR_gaussian_splatting:SH_DEGREE_1_COEF_0": 6, "KHR_gaussian_splatting:SH_DEGREE_1_COEF_1": 7, "KHR_gaussian_splatting:SH_DEGREE_1_COEF_2": 8`
	const TPLSH2 = TPLSH1 + `,"KHR_gaussian_splatting:SH_DEGREE_2_COEF_0": 9, "KHR_gaussian_splatting:SH_DEGREE_2_COEF_1": 10, "KHR_gaussian_splatting:SH_DEGREE_2_COEF_2": 11, "KHR_gaussian_splatting:SH_DEGREE_2_COEF_3": 12, "KHR_gaussian_splatting:SH_DEGREE_2_COEF_4": 13`
	const TPLSH3 = TPLSH2 + `,"KHR_gaussian_splatting:SH_DEGREE_3_COEF_0": 14, "KHR_gaussian_splatting:SH_DEGREE_3_COEF_1": 15, "KHR_gaussian_splatting:SH_DEGREE_3_COEF_2": 16, "KHR_gaussian_splatting:SH_DEGREE_3_COEF_3": 17, "KHR_gaussian_splatting:SH_DEGREE_3_COEF_4": 18, "KHR_gaussian_splatting:SH_DEGREE_3_COEF_5": 19, "KHR_gaussian_splatting:SH_DEGREE_3_COEF_6": 20`
	outputShDegree := GetArgShDegree()
	if outputShDegree < 1 {
		tmpl = cmn.ReplaceAll(tmpl, "$accessorsSH$", "")
		tmpl = cmn.ReplaceAll(tmpl, "$bufferViewsSH$", "")
		tmpl = cmn.ReplaceAll(tmpl, "$SH_DEGREE$", "")
	} else if outputShDegree == 1 {
		accessorsSH := ""
		bufferViewsSH := ""
		for i := 6; i < 9; i++ {
			accessorsSH = accessorsSH + "," + cmn.ReplaceAll(TPL1, "$index", cmn.IntToString(i))
			bufferViewsSH = bufferViewsSH + "," + cmn.ReplaceAll(TPL2, "$index", cmn.IntToString(i))
		}
		tmpl = cmn.ReplaceAll(tmpl, "$accessorsSH$", accessorsSH)
		tmpl = cmn.ReplaceAll(tmpl, "$bufferViewsSH$", bufferViewsSH)
		tmpl = cmn.ReplaceAll(tmpl, "$SH_DEGREE$", TPLSH1)
		maxIndex = 9
	} else if outputShDegree == 2 {
		accessorsSH := ""
		bufferViewsSH := ""
		for i := 6; i < 14; i++ {
			accessorsSH = accessorsSH + "," + cmn.ReplaceAll(TPL1, "$index", cmn.IntToString(i))
			bufferViewsSH = bufferViewsSH + "," + cmn.ReplaceAll(TPL2, "$index", cmn.IntToString(i))
		}
		tmpl = cmn.ReplaceAll(tmpl, "$accessorsSH$", accessorsSH)
		tmpl = cmn.ReplaceAll(tmpl, "$bufferViewsSH$", bufferViewsSH)
		tmpl = cmn.ReplaceAll(tmpl, "$SH_DEGREE$", TPLSH2)
		maxIndex = 14
	} else {
		accessorsSH := ""
		bufferViewsSH := ""
		for i := 6; i < 21; i++ {
			accessorsSH = accessorsSH + "," + cmn.ReplaceAll(TPL1, "$index", cmn.IntToString(i))
			bufferViewsSH = bufferViewsSH + "," + cmn.ReplaceAll(TPL2, "$index", cmn.IntToString(i))
		}
		tmpl = cmn.ReplaceAll(tmpl, "$accessorsSH$", accessorsSH)
		tmpl = cmn.ReplaceAll(tmpl, "$bufferViewsSH$", bufferViewsSH)
		tmpl = cmn.ReplaceAll(tmpl, "$SH_DEGREE$", TPLSH3)
		maxIndex = 21
	}

	var byteLengths []int
	byteLengths = append(byteLengths, count*12) // POSITION, 3 × 4
	byteLengths = append(byteLengths, count*4)  // COLOR_0, 4 × 1
	byteLengths = append(byteLengths, count*16) // ROTATION, 4 × 4
	byteLengths = append(byteLengths, count*12) // SCALE, 3 × 4
	byteLengths = append(byteLengths, count*4)  // OPACITY, 1 × 4
	byteLengths = append(byteLengths, count*12) // SH_DEGREE_0_COEF_0, 3 × 4
	for range 15 {
		byteLengths = append(byteLengths, count*12) // sh1 ~ sh3, 3 × 4
	}

	tmpl = cmn.ReplaceAll(tmpl, "$version", cmn.VER)
	tmpl = cmn.ReplaceAll(tmpl, "$minX", cmn.FormatFloat32(mm.MinX))
	tmpl = cmn.ReplaceAll(tmpl, "$minY", cmn.FormatFloat32(mm.MinY))
	tmpl = cmn.ReplaceAll(tmpl, "$minZ", cmn.FormatFloat32(mm.MinZ))
	tmpl = cmn.ReplaceAll(tmpl, "$maxX", cmn.FormatFloat32(mm.MaxX))
	tmpl = cmn.ReplaceAll(tmpl, "$maxY", cmn.FormatFloat32(mm.MaxY))
	tmpl = cmn.ReplaceAll(tmpl, "$maxZ", cmn.FormatFloat32(mm.MaxZ))
	tmpl = cmn.ReplaceAll(tmpl, "$count", cmn.IntToString(count))
	offset := 0

	for i := range maxIndex {
		tmpl = cmn.ReplaceAll(tmpl, "$byteLength"+cmn.IntToString(i)+"$", cmn.IntToString(byteLengths[i]))
		tmpl = cmn.ReplaceAll(tmpl, "$byteOffset"+cmn.IntToString(i)+"$", cmn.IntToString(offset))
		offset += byteLengths[i]
	}
	tmpl = cmn.ReplaceAll(tmpl, "$byteLengthTotal", cmn.IntToString(offset))

	return cmn.JsonStringify(tmpl)
}
