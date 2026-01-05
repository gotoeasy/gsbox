package gsplat

import (
	"encoding/json"
	"gsbox/cmn"
	"log"
	"math"
	"sort"
)

const FileSizeThreshold = 400000 // 合并阈值

type SplatTiles struct {
	Version          uint                  `json:"version"`
	Magic            string                `json:"magic"`
	Comment          string                `json:"comment,omitempty"`
	LodLevels        uint16                `json:"lodLevels"`
	TotalCount       int                   `json:"totalCount"`
	Environment      string                `json:"environment,omitempty"`
	Files            map[string]*SplatFile `json:"files"`
	Tree             *SplatNode            `json:"tree"`
	ShDegree         uint8                 `json:"shDegree"`
	ShCentroids      []uint8               `json:"-"`
	PaletteSize      int                   `json:"-"`
	EnvironmentDatas []*SplatData          `json:"-"`
}

type SplatNode struct {
	Center   []float32       `json:"center"`
	Radius   float32         `json:"radius"`
	Children *[]*SplatNode   `json:"children,omitempty"`
	Lods     *[]*TileMapping `json:"lods,omitempty"`
	bound    *Bound          `json:"-"`
}

type SplatFile struct {
	FileKey     string       `json:"-"`
	Index       int          `json:"-"`
	Url         string       `json:"url"`
	Lod         uint16       `json:"lod"`
	Seq         int          `json:"-"`
	Count       int          `json:"-"`
	Datas       []*SplatData `json:"-"`
	ShCentroids []uint8      `json:"-"`
}

type SplatTile struct {
	Center []float32      `json:"center"`
	Radius float32        `json:"radius"`
	Lods   []*TileMapping `json:"lods"`
	bound  *Bound         `json:"-"`
}

type TileMapping struct {
	FileKey string       `json:"fileKey"`
	Offset  int          `json:"offset"`
	Count   int          `json:"count"`
	Datas   []*SplatData `json:"-"`
}

type BTreeNode struct {
	datas     []*SplatData
	file      *SplatFile
	children  []*BTreeNode
	isLeaf    bool
	lodCounts []int
	lods      []*TileMapping
	mm        *V3MinMax
	bound     *Bound
	level     int
}

func BuildLodMetaSplatTiles(datas []*SplatData) (*SplatTiles, *LodMeta) {
	splatTiles := buildSplatTilesLodMeta(datas)
	return copyToLodMeta(splatTiles) // 包围盒、文件名、JSON结构
}

func buildSplatTilesLodMeta(datas []*SplatData) *SplatTiles {
	log.Println("[Info] (parameter) cs:", oArg.CutSize, "(cut size)")

	// 调色板
	var shCentroids []uint8
	var paletteSize int
	outputShDegree := GetArgShDegree()

	if outputShDegree > 0 {
		shCentroids, _, paletteSize = ReWriteShByKmeans(datas)

		// 根据输出级别相应的置零
		if outputShDegree < 3 {
			idxs := []int{0, 3, 8}
			cnt := len(shCentroids) / 60 // 60=15*4
			for i := range cnt {
				for d := idxs[outputShDegree]; d < 15; d++ {
					shCentroids[i*60+d*4+0] = 128
					shCentroids[i*60+d*4+1] = 128
					shCentroids[i*60+d*4+2] = 128
				}
			}
		}
	}

	lodLevels := uint16(0)
	for _, d := range datas {
		lodLevels = max(lodLevels, d.Lod)
	}
	lodLevels++

	btree := buildBTree(datas, lodLevels)
	files := buildBTreeFile(btree, lodLevels)

	fileMap := make(map[string]*SplatFile)

	tree := &SplatNode{}
	copyToSplatTreeLod(btree, tree)

	for _, f := range files {
		fileMap[f.FileKey] = f
	}

	del, comment := cmn.RemoveNonASCII(Args.GetArgIgnorecase("-c", "--comment"))
	if del {
		log.Println("[Warn] The existing non-ASCII characters in the comment have been removed!")
	}
	if comment == "" {
		comment = DefaultSpxComment()
	}

	splatTiles := &SplatTiles{
		Version:     1,
		Magic:       "splat-lod",
		Comment:     comment,
		LodLevels:   lodLevels,
		TotalCount:  len(datas),
		Files:       fileMap,
		Tree:        tree,
		ShDegree:    outputShDegree,
		ShCentroids: shCentroids,
		PaletteSize: paletteSize,
	}

	return splatTiles
}

func buildBTree(datas []*SplatData, lodLevels uint16) *BTreeNode {
	defer OnProgress(PhaseCut, 100, 100)

	L := int(math.Ceil(math.Log2(float64(len(datas)/oArg.CutSize)))) + 1
	MaxCutTimes := (1 << L) // 总分割次数 = 2^L - 1
	times := 0
	OnProgress(PhaseCut, times, MaxCutTimes)

	mm := ComputeXyzMinMax(datas)

	lodCounts := make([]int, lodLevels)
	for _, d := range datas {
		lodCounts[d.Lod]++
	}

	rootNode := &BTreeNode{
		datas:     datas,
		isLeaf:    false,
		mm:        mm,
		level:     1,
		lodCounts: lodCounts,
	}

	todos := []*BTreeNode{rootNode}
	for {
		if len(todos) == 0 {
			break
		}

		todoNode := todos[0]
		todos = todos[1:]

		splitBTreeNode(todoNode)
		if !todoNode.isLeaf {
			todos = append(todos, todoNode.children...)
		}

		OnProgress(PhaseCut, times, MaxCutTimes)
		times++
	}
	return rootNode
}

func traveBTree(node *BTreeNode, fnCallBack func(*BTreeNode) bool) {
	if fnCallBack(node) {
		for _, cld := range node.children {
			traveBTree(cld, fnCallBack)
		}
	}
}

func splitBTreeNode(node *BTreeNode) {
	minVal, length, axis := calcLongestAxis(node.mm)
	totalCnt := len(node.datas)
	leftCnt := len(node.datas) / 2
	rightCnt := totalCnt - leftCnt

	node0 := &BTreeNode{lodCounts: make([]int, len(node.lodCounts)), datas: make([]*SplatData, 0, leftCnt)}
	node1 := &BTreeNode{lodCounts: make([]int, len(node.lodCounts)), datas: make([]*SplatData, 0, rightCnt)}

	if totalCnt < 1000 {
		// 小数据量场景：排序取
		switch axis {
		case 0:
			sort.Slice(node.datas, func(i, j int) bool {
				return node.datas[i].PositionX <= node.datas[j].PositionX
			})
		case 1:
			sort.Slice(node.datas, func(i, j int) bool {
				return node.datas[i].PositionY <= node.datas[j].PositionY
			})
		default:
			sort.Slice(node.datas, func(i, j int) bool {
				return node.datas[i].PositionZ <= node.datas[j].PositionZ
			})
		}

		node0.datas = node.datas[:leftCnt]
		node1.datas = node.datas[leftCnt:]
	} else {
		// 大数据量场景：桶数量固定1000
		bucketCnt := 1000
		buckets := make([][]*SplatData, bucketCnt)
		maxIdx := bucketCnt - 1
		inv := float32(bucketCnt) / length

		for _, d := range node.datas {
			var val float32
			switch axis {
			case 0:
				val = d.PositionX
			case 1:
				val = d.PositionY
			default:
				val = d.PositionZ
			}

			idx := min(max(int(math.Floor(float64((val-minVal)*inv))), 0), maxIdx)
			buckets[idx] = append(buckets[idx], d)
		}

		currentLeftCnt := 0

		for i, rows := range buckets {
			rowLen := len(rows)
			if rowLen == 0 {
				continue
			}

			if currentLeftCnt+rowLen < leftCnt {
				node0.datas = append(node0.datas, rows...)
				currentLeftCnt += rowLen
				continue
			} else if currentLeftCnt+rowLen == leftCnt {
				node0.datas = append(node0.datas, rows...)
				currentLeftCnt += rowLen

				// 填充后续桶到右半部分
				for k := i + 1; k < bucketCnt; k++ {
					if len(buckets[k]) > 0 {
						node1.datas = append(node1.datas, buckets[k]...)
					}
				}
				break
			}

			// 排序目标桶
			switch axis {
			case 0:
				sort.Slice(rows, func(i, j int) bool {
					return rows[i].PositionX <= rows[j].PositionX
				})
			case 1:
				sort.Slice(rows, func(i, j int) bool {
					return rows[i].PositionY <= rows[j].PositionY
				})
			default:
				sort.Slice(rows, func(i, j int) bool {
					return rows[i].PositionZ <= rows[j].PositionZ
				})
			}

			// 切割目标桶
			needCnt := leftCnt - currentLeftCnt
			node0.datas = append(node0.datas, rows[:needCnt]...)
			node1.datas = append(node1.datas, rows[needCnt:]...)
			// 填充后续桶到右半部分
			for k := i + 1; k < bucketCnt; k++ {
				if len(buckets[k]) > 0 {
					node1.datas = append(node1.datas, buckets[k]...)
				}
			}
			break
		}
	}

	node.datas = nil

	node.children = make([]*BTreeNode, 0)
	if len(node0.datas) > 0 {
		mm := ComputeXyzMinMax(node0.datas)
		node0.isLeaf = len(node0.datas) <= oArg.CutSize
		node0.mm = mm
		node0.level = node.level + 1
		for _, d := range node0.datas {
			node0.lodCounts[d.Lod]++
		}
		node.children = append(node.children, node0)
	}
	if len(node1.datas) > 0 {
		mm := ComputeXyzMinMax(node1.datas)
		node1.isLeaf = len(node1.datas) <= oArg.CutSize
		node1.mm = mm
		node1.level = node.level + 1
		for _, d := range node1.datas {
			node1.lodCounts[d.Lod]++
		}
		node.children = append(node.children, node1)
	}

	for _, nd := range node.children {
		if !nd.isLeaf {
			continue
		}

		nd.lods = make([]*TileMapping, len(nd.lodCounts))
		for i, v := range nd.lodCounts {
			if v > 0 {
				nd.lods[i] = &TileMapping{}
			}
		}

		for _, d := range nd.datas {
			nd.lods[d.Lod].Datas = append(nd.lods[d.Lod].Datas, d)
		}

		if oArg.isOutputLodMeta {
			nd.bound = calcLodMetaBound(nd.datas)
		}

		nd.datas = nil
	}
}

func buildBTreeFile(root *BTreeNode, lodLevels uint16) []*SplatFile {

	seq := 0

	fnMergeFile := func(mergeNode *BTreeNode, lod uint16) *SplatFile {
		// 收集合并节点下的全部叶节点
		var leafs []*BTreeNode
		count := 0
		traveBTree(mergeNode, func(node *BTreeNode) bool {
			if node.isLeaf && node.lodCounts[lod] > 0 {
				leafs = append(leafs, node)
				count += node.lodCounts[lod]
			}
			return true
		})

		// 创建文件节点
		splatFile := &SplatFile{
			FileKey: cmn.IntToString(int(lod)) + "_" + cmn.IntToString(seq),
			Lod:     lod,
			Seq:     seq,
			Count:   count,
			Datas:   make([]*SplatData, 0, count),
		}
		seq++

		// 汇总数据，记录偏移位置
		offset := 0
		for _, leaf := range leafs {
			leaf.lods[lod].FileKey = splatFile.FileKey
			leaf.lods[lod].Offset = offset
			leaf.lods[lod].Count = len(leaf.lods[lod].Datas)
			SortMorton(leaf.lods[lod].Datas) // 叶节点数据排序
			splatFile.Datas = append(splatFile.Datas, leaf.lods[lod].Datas...)
			offset += leaf.lods[lod].Count

			leaf.lods[lod].Datas = nil
		}

		// 挂载文件节点
		return splatFile
	}

	// 从粗到细遍历lod，合并文件节点数据
	var files []*SplatFile
	for i := range lodLevels {
		seq = 0
		// 收集待合并节点
		var mergeNodes []*BTreeNode
		traveBTree(root, func(node *BTreeNode) bool {
			if node.lodCounts[i] > 0 && node.lodCounts[i] <= FileSizeThreshold {
				mergeNodes = append(mergeNodes, node)
				return false
			}
			return true
		})
		// 依次合并
		for _, node := range mergeNodes {
			files = append(files, fnMergeFile(node, i))
		}
	}
	return files
}

func copyToSplatTreeLod(treeNode *BTreeNode, splatNode *SplatNode) {
	splatNode.Center = []float32{treeNode.mm.CenterX, treeNode.mm.CenterY, treeNode.mm.CenterZ}
	splatNode.Radius = treeNode.mm.Radius

	if treeNode.isLeaf {
		var tileMappings []*TileMapping
		splatNode.Lods = &tileMappings
		*splatNode.Lods = append(*splatNode.Lods, treeNode.lods...)

		if oArg.isOutputLodMeta {
			splatNode.bound = treeNode.bound
		}
	} else {
		var children []*SplatNode
		splatNode.Children = &children
		for _, tnode := range treeNode.children {
			snode := &SplatNode{}
			*splatNode.Children = append(*splatNode.Children, snode)
			copyToSplatTreeLod(tnode, snode)
		}
	}
}

func calcLodMetaBound(datas []*SplatData) *Bound {
	mins := []float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}
	maxs := []float32{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}
	for _, d := range datas {
		x := d.PositionX
		y := d.PositionY
		z := d.PositionZ
		rw := cmn.DecodeSplatRotation(d.RotationW)
		rx := cmn.DecodeSplatRotation(d.RotationX)
		ry := cmn.DecodeSplatRotation(d.RotationY)
		rz := cmn.DecodeSplatRotation(d.RotationZ)
		sx := cmn.EncodeSplatScale(d.ScaleX)
		sy := cmn.EncodeSplatScale(d.ScaleY)
		sz := cmn.EncodeSplatScale(d.ScaleZ)

		aabbMins, aabbMaxs := CalcSplatAABB(x, y, z, rx, ry, rz, rw, sx, sy, sz)
		mins[0] = min(mins[0], aabbMins[0])
		mins[1] = min(mins[1], aabbMins[1])
		mins[2] = min(mins[2], aabbMins[2])
		maxs[0] = max(maxs[0], aabbMaxs[0])
		maxs[1] = max(maxs[1], aabbMaxs[1])
		maxs[2] = max(maxs[2], aabbMaxs[2])
	}
	return &Bound{Min: mins, Max: maxs}
}

func setSplatTreeBound(node *SplatNode) *Bound {
	if node.bound == nil {
		var bounds []*Bound
		for _, cld := range *node.Children {
			bounds = append(bounds, setSplatTreeBound(cld))
		}

		bound := &Bound{
			Min: []float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32},
			Max: []float32{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32},
		}
		for _, bd := range bounds {
			bound.Min[0] = min(bound.Min[0], bd.Min[0])
			bound.Min[1] = min(bound.Min[1], bd.Min[1])
			bound.Min[2] = min(bound.Min[2], bd.Min[2])
			bound.Max[0] = max(bound.Max[0], bd.Max[0])
			bound.Max[1] = max(bound.Max[1], bd.Max[1])
			bound.Max[2] = max(bound.Max[2], bd.Max[2])
		}
		node.bound = bound
	}
	return node.bound
}

func copyToLodMeta(splatTiles *SplatTiles) (*SplatTiles, *LodMeta) {

	setSplatTreeBound(splatTiles.Tree)

	filenames := setSplatFileIndex(splatTiles)

	tree := &LodNode{}
	copyToSplatTreeLodMeta(splatTiles, splatTiles.Tree, tree)

	lodMeta := &LodMeta{
		LodLevels:   int(splatTiles.LodLevels),
		Filenames:   filenames,
		Environment: splatTiles.Environment,
		Tree:        tree,
	}

	return splatTiles, lodMeta
}

func setSplatFileIndex(splatTiles *SplatTiles) []string {
	var files []*SplatFile
	for _, v := range splatTiles.Files {
		files = append(files, v)
	}
	sort.Slice(files, func(i, j int) bool {
		if files[i].Lod == files[j].Lod {
			return files[i].Seq < files[j].Seq
		}
		return files[i].Lod < files[j].Lod
	})

	filenames := make([]string, len(files))
	isMetaJson := Args.GetArgIgnorecase("-of", "--output-format") == "meta.json"
	for i, v := range files {
		v.Index = i
		lod_seq := cmn.IntToString(int(v.Lod)) + "_" + cmn.IntToString(v.Seq)
		if isMetaJson {
			v.Url = lod_seq + "/meta.json"
		} else {
			v.Url = lod_seq + ".sog"
		}
		filenames[i] = v.Url
	}
	return filenames
}

func copyToSplatTreeLodMeta(splatTiles *SplatTiles, splatNode *SplatNode, lodNode *LodNode) {
	lodNode.Bound = splatNode.bound

	if splatNode.Lods != nil {
		lods := make(map[string]*LodMapping)
		for _, d := range *splatNode.Lods {
			if d != nil {
				splatFile := splatTiles.Files[d.FileKey]
				lod := cmn.IntToString(int(splatFile.Lod))
				lods[lod] = &LodMapping{File: splatFile.Index, Offset: d.Offset, Count: d.Count}
			}
		}
		lodNode.Lods = &lods
	} else {
		var children []*LodNode
		lodNode.Children = &children
		for _, snode := range *splatNode.Children {
			lnode := &LodNode{}
			*lodNode.Children = append(*lodNode.Children, lnode)
			copyToSplatTreeLodMeta(splatTiles, snode, lnode)
		}
	}
}

func calcLongestAxis(mm *V3MinMax) (float32, float32, int) {
	if mm.LenX >= mm.LenY && mm.LenX >= mm.LenZ {
		return mm.MinX, mm.LenX, 0
	} else if mm.LenY >= mm.LenX && mm.LenY >= mm.LenZ {
		return mm.MinY, mm.LenY, 1
	}
	return mm.MinZ, mm.LenZ, 2
}

func (s *SplatTiles) ToJson() string {
	// data, err := json.MarshalIndent(s, "", "  ")
	data, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(data)
}
