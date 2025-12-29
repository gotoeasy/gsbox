package gsplat

import (
	"encoding/json"
)

func ParseSogMeta(jsonStr string) (*SogMeta, error) {
	var meta SogMeta
	err := json.Unmarshal([]byte(jsonStr), &meta)
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

// 定义整个 JSON 文件的结构
type SogMeta struct {
	/** v2 */
	Version int `json:"version,omitempty"`
	/** v2 */
	Count  int        `json:"count,omitempty"`
	Means  *SogMeans  `json:"means,omitempty"`
	Scales *SogScales `json:"scales,omitempty"`
	Quats  *SogQuats  `json:"quats,omitempty"`
	Sh0    *SogSh0    `json:"sh0,omitempty"`
	ShN    *SogShN    `json:"shN,omitempty"`
}

// 定义 means 字段的结构
type SogMeans struct {
	Shape []int  `json:"shape,omitempty"`
	Dtype string `json:"dtype,omitempty"`
	/** v2 */
	Mins []float32 `json:"mins,omitempty"`
	/** v2 */
	Maxs []float32 `json:"maxs,omitempty"`
	/** v2 */
	Files []string `json:"files,omitempty"`
}

// 定义了 scales 字段的结构
type SogScales struct {
	Shape []int     `json:"shape,omitempty"`
	Dtype string    `json:"dtype,omitempty"`
	Mins  []float32 `json:"mins,omitempty"`
	Maxs  []float32 `json:"maxs,omitempty"`
	/** v2 */
	Codebook []float32 `json:"codebook,omitempty"`
	/** v2 */
	Files []string `json:"files,omitempty"`
}

// 定义了 quats 字段的结构
type SogQuats struct {
	Shape    []int  `json:"shape,omitempty"`
	Dtype    string `json:"dtype,omitempty"`
	Encoding string `json:"encoding,omitempty"`
	/** v2 */
	Files []string `json:"files,omitempty"`
}

// 定义了 sh0 字段的结构
type SogSh0 struct {
	Shape []int     `json:"shape,omitempty"`
	Dtype string    `json:"dtype,omitempty"`
	Mins  []float32 `json:"mins,omitempty"`
	Maxs  []float32 `json:"maxs,omitempty"`
	/** v2 */
	Codebook []float32 `json:"codebook,omitempty"`
	/** v2 */
	Files []string `json:"files,omitempty"`
}

// 定义了 shN 字段的结构
type SogShN struct {
	Shape        []int   `json:"shape,omitempty"`
	Dtype        string  `json:"dtype,omitempty"`
	Mins         float32 `json:"mins,omitempty"`
	Maxs         float32 `json:"maxs,omitempty"`
	Quantization int     `json:"quantization,omitempty"`
	/** v2 */
	Count    int       `json:"count"`
	Bands    uint8     `json:"bands"`
	Codebook []float32 `json:"codebook,omitempty"`
	Files    []string  `json:"files,omitempty"`
}

type LodMeta struct {
	LodLevels   int      `json:"lodLevels"`
	Filenames   []string `json:"filenames"`
	Environment string   `json:"environment"`
	Tree        *LodNode `json:"tree"`
}

type LodNode struct {
	Bound    *Bound                  `json:"bound"`
	Children *[]*LodNode             `json:"children,omitempty"`
	Lods     *map[string]*LodMapping `json:"lods,omitempty"`
}

type Bound struct {
	Min []float32 `json:"min"`
	Max []float32 `json:"max"`
}

type LodMapping struct {
	File   int `json:"file"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
}
