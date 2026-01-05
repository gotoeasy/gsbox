package gsplat

import (
	"encoding/json"
	"errors"
	"gsbox/cmn"
	"path/filepath"
)

func ReadLodMeta(input string) (dataShDegree uint8, lodLevels uint16, lodDatas []*SplatData, envDatas []*SplatData) {
	cmn.ExitOnConditionError(!cmn.Endwiths(input, "lod-meta.json"), errors.New("invalid lod-meta.json: "+input))

	jsonStr, err := cmn.ReadFileString(input)
	cmn.ExitOnError(err)

	var lodMeta LodMeta
	err = json.Unmarshal([]byte(jsonStr), &lodMeta)
	cmn.ExitOnError(err)

	lodMetaDir := cmn.Dir(input)

	// 汇总文件
	mapFiles := make(map[string]*SplatFile)
	for i, filename := range lodMeta.Filenames {
		splatFile := &SplatFile{FileKey: cmn.IntToString(i), Url: filepath.Join(lodMetaDir, filename)}
		mapFiles[splatFile.FileKey] = splatFile
	}

	// 设定文件LOD级别
	traveTree(lodMeta.Tree, func(node *LodNode) bool {
		if node.Lods != nil {
			for k, v := range *node.Lods {
				fileKey := cmn.IntToString(v.File)
				mapFiles[fileKey].Lod = uint16(cmn.StringToInt(k))
			}
		}
		return true
	})

	// 计算LOD最大级别
	lodLevels = 0
	for _, splatFile := range mapFiles {
		lodLevels = max(lodLevels, splatFile.Lod+1)
	}

	// 数据部分只支持 *.sog 或 meta.json
	lodDatas = make([]*SplatData, 0)
	for _, splatFile := range mapFiles {
		datas, header := ReadSog(splatFile.Url)
		dataShDegree = max(dataShDegree, header.ShDegree)
		for _, d := range datas {
			d.Lod = splatFile.Lod
		}
		lodDatas = append(lodDatas, datas...)
	}

	// 环境天空盒，默认 *.sog，保险起见支持多种格式
	if lodMeta.Environment != "" {
		fileEnv := filepath.Join(lodMetaDir, lodMeta.Environment)
		if cmn.Endwiths(fileEnv, ".sog", true) || cmn.FileName(fileEnv) == "meta.json" {
			envDatas, _ = ReadSog(fileEnv)
		} else if cmn.Endwiths(fileEnv, ".spx", true) {
			_, envDatas = ReadSpx(fileEnv)
		} else if cmn.Endwiths(fileEnv, ".spz", true) {
			_, envDatas = ReadSpz(fileEnv)
		} else if cmn.Endwiths(fileEnv, ".ply", true) {
			_, envDatas = ReadPly(fileEnv)
		} else if cmn.Endwiths(fileEnv, ".splat", true) {
			envDatas = ReadSplat(fileEnv)
		}
	}

	return
}

func traveTree(node *LodNode, fnCallBack func(nd *LodNode) bool) {
	if fnCallBack(node) && node.Children != nil {
		for _, cld := range *node.Children {
			traveTree(cld, fnCallBack)
		}
	}
}
