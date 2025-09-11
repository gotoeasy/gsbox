package gsplat

import (
	"fmt"
	"gsbox/cmn"
	"sort"
)

type DataFreq struct {
	Bytes3 [3]uint8 // R, G, B
	Count  int
	Splat  *SplatData
}

type FrequencyCounter struct {
	freqMap    map[[3]uint8]int
	TotalCount int
}

func NewFrequencyCounter() *FrequencyCounter {
	return &FrequencyCounter{
		freqMap:    make(map[[3]uint8]int),
		TotalCount: 0,
	}
}

func (fc *FrequencyCounter) CountByScale(datas []*SplatData) {
	fc.TotalCount = len(datas)
	for _, splat := range datas {
		key := [3]uint8{cmn.EncodeSpxScale(splat.ScaleX), cmn.EncodeSpxScale(splat.ScaleY), cmn.EncodeSpxScale(splat.ScaleZ)}
		fc.freqMap[key]++
	}
}

func (fc *FrequencyCounter) CountByColor(datas []*SplatData) {
	fc.TotalCount = len(datas)
	for _, splat := range datas {
		key := [3]uint8{splat.ColorR, splat.ColorG, splat.ColorB}
		fc.freqMap[key]++
	}
}

func (fc *FrequencyCounter) CountByRotation(datas []*SplatData) {
	fc.TotalCount = len(datas)
	for _, splat := range datas {
		key := [3]uint8{splat.RotationX, splat.RotationY, splat.RotationZ}
		fc.freqMap[key]++
	}
}

func (fc *FrequencyCounter) GetTopN(topCnt int) []*DataFreq {
	// 将所有颜色频率对转换为切片
	freqs := make([]*DataFreq, 0, len(fc.freqMap))
	for bytes3, count := range fc.freqMap {
		freqs = append(freqs, &DataFreq{Bytes3: bytes3, Count: count})
	}

	// 按频率降序排序
	sort.Slice(freqs, func(i, j int) bool {
		return freqs[i].Count > freqs[j].Count
	})

	// 取前N个
	if len(freqs) > topCnt {
		freqs = freqs[:topCnt]
	}

	cnt := 0
	for _, cf := range freqs {
		cnt += cf.Count
	}

	oldVal := fc.TotalCount * 19
	newVal := len(freqs)*3 + cnt*1 + (fc.TotalCount-cnt)*4 + fc.TotalCount*16
	fmt.Printf("合计重复 %v/%v,  %.1f%%,  %.1f%%\n", cnt, fc.TotalCount, 100.0*float64(cnt)/float64(fc.TotalCount), 100.0*float64(newVal)/float64(oldVal))

	return freqs
}
