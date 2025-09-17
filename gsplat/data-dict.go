package gsplat

import "sort"

// 生成字典
func GenerateDict(data []byte) []byte {
	// 统计每个字节的出现次数
	mapByteCount := make(map[byte]int)
	for _, value := range data {
		mapByteCount[value]++
	}

	// 将字节按出现次数降序排序，如果出现次数相同，则按字节值升序排序
	type valueCount struct {
		Byte  byte
		Count int
	}
	var valueCounts []valueCount
	for byteValue, count := range mapByteCount {
		valueCounts = append(valueCounts, valueCount{byteValue, count})
	}
	sort.Slice(valueCounts, func(i, j int) bool {
		if valueCounts[i].Count == valueCounts[j].Count {
			return valueCounts[i].Byte < valueCounts[j].Byte
		}
		return valueCounts[i].Count > valueCounts[j].Count
	})

	// 生成
	dict := make([]byte, 256)
	for i, bc := range valueCounts {
		dict[i] = bc.Byte
	}
	return dict
}

// 按字典编码
func EncodeByDict(data []byte, dict []byte) []byte {
	// 创建一个映射表，将字节值映射到索引值
	mapByteIndex := make(map[byte]byte)
	for i, byteValue := range dict {
		mapByteIndex[byteValue] = byte(i)
	}

	// 将原始数据编码为索引值数组
	encodedData := make([]byte, len(data))
	for i, byteValue := range data {
		encodedData[i] = mapByteIndex[byteValue]
	}
	return encodedData
}

// 按字典解码
func DecodeByDict(encodedData []byte, dict []byte) []byte {
	// 将索引值数组解码为原始字节数组
	decodedData := make([]byte, len(encodedData))
	for i, dictIdx := range encodedData {
		decodedData[i] = dict[dictIdx]
	}
	return decodedData
}
