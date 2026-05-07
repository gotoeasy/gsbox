package gsplat

import (
	"encoding/json"
	"errors"
	"gsbox/cmn"
	"os"
)

func ReadGlb(glbFile string) (shDegree uint8, datas []*SplatData) {
	file, err := os.Open(glbFile)
	cmn.ExitOnError(err)
	defer file.Close()

	bs := make([]byte, 20)
	_, err = file.Read(bs)
	cmn.ExitOnError(err)

	jsonLength := int(cmn.BytesToUint32(bs[12:16])) // 实际长度
	jsonLength4 := (jsonLength + 3) &^ 3            // 4字节对齐长度

	bs = make([]byte, jsonLength4)
	_, err = file.Read(bs)
	cmn.ExitOnError(err)

	strJson := cmn.BytesToString(bs)
	if !cmn.Contains(strJson, KHR_gaussian_splatting) {
		cmn.ExitOnError(errors.New("GLB without KHR_gaussian_splatting extension is not supported"))
	}

	if cmn.Contains(strJson, KHR_gaussian_splatting_compression_spz_2) {
		return readGlb_KHR_gaussian_splatting_compression_spz_2(file, strJson)
	}
	return readGlb_KHR_gaussian_splatting(file, strJson)
}

func readGlb_KHR_gaussian_splatting_compression_spz_2(file *os.File, strJson string) (shDegree uint8, datas []*SplatData) {
	var oJson any
	cmn.ExitOnError(json.Unmarshal([]byte(strJson), &oJson))

	bs := make([]byte, 8)
	_, err := file.Read(bs)
	cmn.ExitOnError(err)

	byteLength := getGsBufferViewsByteLength(oJson)

	bs = make([]byte, byteLength) // spz
	_, err = file.Read(bs)
	cmn.ExitOnError(err)

	ungzipDatas, err := cmn.DecompressGzip(bs)
	cmn.ExitOnConditionError(err != nil, errors.New("[SPZ ERROR] UnGzip failed"))
	cmn.ExitOnConditionError(len(ungzipDatas) < HeaderSizeSpzV3, errors.New("[SPZ ERROR] Invalid spz header"))

	header := ParseSpzHeader(ungzipDatas[0:HeaderSizeSpzV3])
	shDegree = header.ShDegree
	datas = readSpzDatasV2V3(ungzipDatas[HeaderSizeSpzV3:], header)
	return
}

func readGlb_KHR_gaussian_splatting(file *os.File, strJson string) (shDegree uint8, datas []*SplatData) {
	var oJson any
	cmn.ExitOnError(json.Unmarshal([]byte(strJson), &oJson))

	splatCount := getGsAccessors0Count(oJson)
	positionIdx := getGsFieldIndex(oJson, "POSITION")
	scaleIdx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SCALE")
	sh0Idx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SH_DEGREE_0_COEF_0")
	opacityIdx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:OPACITY")
	rotationIdx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:ROTATION")
	sh1Idx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SH_DEGREE_1_COEF_0")
	sh2Idx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SH_DEGREE_1_COEF_1")
	sh3Idx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SH_DEGREE_1_COEF_2")
	sh4Idx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SH_DEGREE_2_COEF_0")
	sh5Idx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SH_DEGREE_2_COEF_1")
	sh6Idx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SH_DEGREE_2_COEF_2")
	sh7Idx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SH_DEGREE_2_COEF_3")
	sh8Idx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SH_DEGREE_2_COEF_4")
	sh9Idx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SH_DEGREE_3_COEF_0")
	sh10Idx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SH_DEGREE_3_COEF_1")
	sh11Idx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SH_DEGREE_3_COEF_2")
	sh12Idx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SH_DEGREE_3_COEF_3")
	sh13Idx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SH_DEGREE_3_COEF_4")
	sh14Idx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SH_DEGREE_3_COEF_5")
	sh15Idx := getGsFieldIndex(oJson, "KHR_gaussian_splatting:SH_DEGREE_3_COEF_6")

	offsetPosition := getGsBufferViewsIndexByteOffset(oJson, positionIdx)
	offsetScale := getGsBufferViewsIndexByteOffset(oJson, scaleIdx)
	offsetRgb := getGsBufferViewsIndexByteOffset(oJson, sh0Idx)
	offsetOpacity := getGsBufferViewsIndexByteOffset(oJson, opacityIdx)
	offsetRotation := getGsBufferViewsIndexByteOffset(oJson, rotationIdx)
	offsetSh1 := getGsBufferViewsIndexByteOffset(oJson, sh1Idx)
	offsetSh2 := getGsBufferViewsIndexByteOffset(oJson, sh2Idx)
	offsetSh3 := getGsBufferViewsIndexByteOffset(oJson, sh3Idx)
	offsetSh4 := getGsBufferViewsIndexByteOffset(oJson, sh4Idx)
	offsetSh5 := getGsBufferViewsIndexByteOffset(oJson, sh5Idx)
	offsetSh6 := getGsBufferViewsIndexByteOffset(oJson, sh6Idx)
	offsetSh7 := getGsBufferViewsIndexByteOffset(oJson, sh7Idx)
	offsetSh8 := getGsBufferViewsIndexByteOffset(oJson, sh8Idx)
	offsetSh9 := getGsBufferViewsIndexByteOffset(oJson, sh9Idx)
	offsetSh10 := getGsBufferViewsIndexByteOffset(oJson, sh10Idx)
	offsetSh11 := getGsBufferViewsIndexByteOffset(oJson, sh11Idx)
	offsetSh12 := getGsBufferViewsIndexByteOffset(oJson, sh12Idx)
	offsetSh13 := getGsBufferViewsIndexByteOffset(oJson, sh13Idx)
	offsetSh14 := getGsBufferViewsIndexByteOffset(oJson, sh14Idx)
	offsetSh15 := getGsBufferViewsIndexByteOffset(oJson, sh15Idx)

	if positionIdx < 0 || scaleIdx < 0 || sh0Idx < 0 || opacityIdx < 0 || rotationIdx < 0 {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	bs := make([]byte, 8)
	_, err := file.Read(bs)
	cmn.ExitOnError(err)

	shDegree = 0
	if sh3Idx > 0 {
		shDegree = 1
	}
	if sh8Idx > 0 {
		shDegree = 2
	}
	if sh15Idx > 0 {
		shDegree = 3
	}

	binLength := cmn.BytesToUint32(bs[0:4]) // 长度
	bs = make([]byte, binLength)
	_, err = file.Read(bs)
	cmn.ExitOnError(err)

	datas = make([]*SplatData, splatCount)
	for i := range splatCount {
		data := &SplatData{}
		data.PositionX = cmn.BytesToFloat32(bs[offsetPosition+i*12+0 : offsetPosition+i*12+4])
		data.PositionY = cmn.BytesToFloat32(bs[offsetPosition+i*12+4 : offsetPosition+i*12+8])
		data.PositionZ = cmn.BytesToFloat32(bs[offsetPosition+i*12+8 : offsetPosition+i*12+12])
		data.ScaleX = cmn.BytesToFloat32(bs[offsetScale+i*12+0 : offsetScale+i*12+4])
		data.ScaleY = cmn.BytesToFloat32(bs[offsetScale+i*12+4 : offsetScale+i*12+8])
		data.ScaleZ = cmn.BytesToFloat32(bs[offsetScale+i*12+8 : offsetScale+i*12+12])
		data.ColorR = cmn.EncodeSplatColor(float64(cmn.BytesToFloat32(bs[offsetRgb+i*12+0 : offsetRgb+i*12+4])))
		data.ColorG = cmn.EncodeSplatColor(float64(cmn.BytesToFloat32(bs[offsetRgb+i*12+4 : offsetRgb+i*12+8])))
		data.ColorB = cmn.EncodeSplatColor(float64(cmn.BytesToFloat32(bs[offsetRgb+i*12+8 : offsetRgb+i*12+12])))
		data.ColorA = cmn.EncodeSplatOpacity(float64(cmn.BytesToFloat32(bs[offsetOpacity+i*4+0 : offsetOpacity+i*4+4])))
		data.RotationX = cmn.EncodeSplatRotation(float64(cmn.BytesToFloat32(bs[offsetRotation+i*16+0 : offsetRotation+i*16+4])))
		data.RotationY = cmn.EncodeSplatRotation(float64(cmn.BytesToFloat32(bs[offsetRotation+i*16+4 : offsetRotation+i*16+8])))
		data.RotationZ = cmn.EncodeSplatRotation(float64(cmn.BytesToFloat32(bs[offsetRotation+i*16+8 : offsetRotation+i*16+12])))
		data.RotationW = cmn.EncodeSplatRotation(float64(cmn.BytesToFloat32(bs[offsetRotation+i*16+12 : offsetRotation+i*16+16])))
		datas[i] = data
	}

	if shDegree > 0 {

		for i, d := range datas {
			d.SH45 = InitZeroSH45()
			d.SH45[0] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh1+i*12+0 : offsetSh1+i*12+4]))))
			d.SH45[1] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh1+i*12+4 : offsetSh1+i*12+8]))))
			d.SH45[2] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh1+i*12+8 : offsetSh1+i*12+12]))))
			d.SH45[3] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh2+i*12+0 : offsetSh2+i*12+4]))))
			d.SH45[4] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh2+i*12+4 : offsetSh2+i*12+8]))))
			d.SH45[5] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh2+i*12+8 : offsetSh2+i*12+12]))))
			d.SH45[6] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh3+i*12+0 : offsetSh3+i*12+4]))))
			d.SH45[7] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh3+i*12+4 : offsetSh3+i*12+8]))))
			d.SH45[8] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh3+i*12+8 : offsetSh3+i*12+12]))))
		}

		if shDegree > 1 {
			for i, d := range datas {
				d.SH45[9] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh4+i*12+0 : offsetSh4+i*12+4]))))
				d.SH45[10] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh4+i*12+4 : offsetSh4+i*12+8]))))
				d.SH45[11] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh4+i*12+8 : offsetSh4+i*12+12]))))
				d.SH45[12] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh5+i*12+0 : offsetSh5+i*12+4]))))
				d.SH45[13] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh5+i*12+4 : offsetSh5+i*12+8]))))
				d.SH45[14] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh5+i*12+8 : offsetSh5+i*12+12]))))
				d.SH45[15] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh6+i*12+0 : offsetSh6+i*12+4]))))
				d.SH45[16] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh6+i*12+4 : offsetSh6+i*12+8]))))
				d.SH45[17] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh6+i*12+8 : offsetSh6+i*12+12]))))
				d.SH45[18] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh7+i*12+0 : offsetSh7+i*12+4]))))
				d.SH45[19] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh7+i*12+4 : offsetSh7+i*12+8]))))
				d.SH45[20] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh7+i*12+8 : offsetSh7+i*12+12]))))
				d.SH45[21] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh8+i*12+0 : offsetSh8+i*12+4]))))
				d.SH45[22] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh8+i*12+4 : offsetSh8+i*12+8]))))
				d.SH45[23] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh8+i*12+8 : offsetSh8+i*12+12]))))
			}
		}
		if shDegree > 2 {
			for i, d := range datas {
				d.SH45[24] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh9+i*12+0 : offsetSh9+i*12+4]))))
				d.SH45[25] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh9+i*12+4 : offsetSh9+i*12+8]))))
				d.SH45[26] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh9+i*12+8 : offsetSh9+i*12+12]))))
				d.SH45[27] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh10+i*12+0 : offsetSh10+i*12+4]))))
				d.SH45[28] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh10+i*12+4 : offsetSh10+i*12+8]))))
				d.SH45[29] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh10+i*12+8 : offsetSh10+i*12+12]))))
				d.SH45[30] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh11+i*12+0 : offsetSh11+i*12+4]))))
				d.SH45[31] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh11+i*12+4 : offsetSh11+i*12+8]))))
				d.SH45[32] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh11+i*12+8 : offsetSh11+i*12+12]))))
				d.SH45[33] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh12+i*12+0 : offsetSh12+i*12+4]))))
				d.SH45[34] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh12+i*12+4 : offsetSh12+i*12+8]))))
				d.SH45[35] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh12+i*12+8 : offsetSh12+i*12+12]))))
				d.SH45[36] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh13+i*12+0 : offsetSh13+i*12+4]))))
				d.SH45[37] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh13+i*12+4 : offsetSh13+i*12+8]))))
				d.SH45[38] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh13+i*12+8 : offsetSh13+i*12+12]))))
				d.SH45[39] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh14+i*12+0 : offsetSh14+i*12+4]))))
				d.SH45[40] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh14+i*12+4 : offsetSh14+i*12+8]))))
				d.SH45[41] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh14+i*12+8 : offsetSh14+i*12+12]))))
				d.SH45[42] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh15+i*12+0 : offsetSh15+i*12+4]))))
				d.SH45[43] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh15+i*12+4 : offsetSh15+i*12+8]))))
				d.SH45[44] = cmn.EncodeSplatSH((float64(cmn.BytesToFloat32(bs[offsetSh15+i*12+8 : offsetSh15+i*12+12]))))
			}
		}
	}

	return shDegree, datas
}

// oJson.meshes[0].primitives[0].attributes['POSITION']
func getGsFieldIndex(oJson any, fieldName string) int {
	// 顶层必须是 map
	root, ok := oJson.(map[string]any)
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	// 获取 meshes 字段，应为数组
	meshesAny, ok := root["meshes"]
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}
	meshes, ok := meshesAny.([]any)
	if !ok || len(meshes) == 0 {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	// 获取第一个 mesh，应为对象
	mesh0, ok := meshes[0].(map[string]any)
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	// 获取 primitives 字段，应为数组
	primitivesAny, ok := mesh0["primitives"]
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}
	primitives, ok := primitivesAny.([]any)
	if !ok || len(primitives) == 0 {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	// 获取第一个 primitive，应为对象
	prim0, ok := primitives[0].(map[string]any)
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	// 获取 attributes 字段，应为对象
	attrsAny, ok := prim0["attributes"]
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}
	attrs, ok := attrsAny.(map[string]any)
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	// 获取字段
	posAny, ok := attrs[fieldName]
	if !ok {
		return -1 // 指定属性不存在，用-1表示
	}

	// 断言为数值（JSON 数值默认是 float64）
	posIdxFloat, ok := posAny.(float64)
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}
	return int(posIdxFloat)
}

// oJson.bufferViews[index].byteOffset
func getGsBufferViewsIndexByteOffset(oJson any, index int) int {
	root, ok := oJson.(map[string]any)
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	bvsAny, ok := root["bufferViews"]
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}
	bvs, ok := bvsAny.([]any)
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	bv, ok := bvs[index].(map[string]any)
	if !ok {
		return -1 // 读取失败，不报错
	}

	offsetAny, ok := bv["byteOffset"]
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}
	offset, ok := offsetAny.(float64)
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	return int(offset)
}

// oJson.accessors[0].count
func getGsAccessors0Count(oJson any) int {
	root, ok := oJson.(map[string]any)
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	accsAny, ok := root["accessors"]
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}
	accs, ok := accsAny.([]any)
	if !ok || len(accs) == 0 {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	acc0, ok := accs[0].(map[string]any)
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	countAny, ok := acc0["count"]
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}
	count, ok := countAny.(float64) // JSON 数字解析为 float64
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}
	return int(count)
}

// oJson.bufferViews[0].byteLength
func getGsBufferViewsByteLength(oJson any) int {
	root, ok := oJson.(map[string]any)
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	bvsAny, ok := root["bufferViews"]
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}
	bvs, ok := bvsAny.([]any)
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	bv, ok := bvs[0].(map[string]any)
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	byteLengthAny, ok := bv["byteLength"]
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}
	offset, ok := byteLengthAny.(float64)
	if !ok {
		cmn.ExitOnError(errors.New("unsupported glb"))
	}

	return int(offset)
}

func ReadGlbJson(glbFile string) string {
	file, err := os.Open(glbFile)
	cmn.ExitOnError(err)
	defer file.Close()

	bs := make([]byte, 20)
	_, err = file.Read(bs)
	cmn.ExitOnError(err)

	jsonLength := int(cmn.BytesToUint32(bs[12:16]))

	bs = make([]byte, jsonLength)
	_, err = file.Read(bs)
	cmn.ExitOnError(err)

	strJson := cmn.BytesToString(bs)
	return cmn.JsonStringify(strJson, true)
}
