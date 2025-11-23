package gsplat

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"gsbox/cmn"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
)

const SH_C0 float64 = 0.28209479177387814

func ReadPlyHeader(plyFile string) (*PlyHeader, error) {
	file, err := os.Open(plyFile)
	cmn.ExitOnError(err)
	defer file.Close()
	return getPlyHeader(file, 2048)
}

func ReadPly(plyFile string) (*PlyHeader, []*SplatData) {
	isNetFile := cmn.IsNetFile(plyFile)
	if isNetFile {
		tmpdir, err := cmn.CreateTempDir()
		cmn.ExitOnError(err)
		downloadFile := filepath.Join(tmpdir, cmn.FileName(plyFile))
		log.Println("[Info]", "download start,", plyFile)
		err = cmn.HttpDownload(plyFile, downloadFile, nil)
		cmn.RemoveAllFileIfError(err, tmpdir)
		cmn.ExitOnError(err)
		log.Println("[Info]", "download finish")
		plyFile = downloadFile
	}

	file, err := os.Open(plyFile)
	cmn.ExitOnError(err)
	defer func() {
		file.Close()
		if isNetFile {
			cmn.RemoveAllFile(cmn.Dir(plyFile))
		}
	}()

	header, err := getPlyHeader(file, 2048)
	cmn.ExitOnError(err)

	if !header.IsOfficialPly() && !header.IsCompressedPly() {
		cmn.ExitOnError(errors.New("unsupported ply file: " + plyFile))
	}

	datas := make([]*SplatData, header.VertexCount)
	if header.ChunkCount == 0 {
		// 标准3dgs的ply
		_, err = file.Seek(int64(header.HeaderLength), 0)     // 定位到数据开始位置
		cmn.ExitOnError(err)                                  // 可能出错
		reader := bufio.NewReaderSize(file, 2*1024*1024)      // 2MB 缓冲区
		const batchSize = 4096                                // 每次读取的最大点数
		dataBytes := make([]byte, batchSize*header.RowLength) // 预分配缓冲区

		for i := 0; i < header.VertexCount; i += batchSize {
			// 计算本次读取的点数
			batchCount := batchSize
			if i+batchCount > header.VertexCount {
				batchCount = header.VertexCount - i
			}

			// 一次性读取一批数据到预分配缓冲区
			readSize := batchCount * header.RowLength
			_, err := io.ReadFull(reader, dataBytes[:readSize])
			cmn.ExitOnError(err)

			// 按点处理这批数据
			for j := 0; j < batchCount; j++ {
				offset := j * header.RowLength
				data := &SplatData{}
				rowBytes := dataBytes[offset:]
				data.PositionX = float32(readValue(header, "x", rowBytes))
				data.PositionY = float32(readValue(header, "y", rowBytes))
				data.PositionZ = float32(readValue(header, "z", rowBytes))
				data.ScaleX = float32(readValue(header, "scale_0", rowBytes))
				data.ScaleY = float32(readValue(header, "scale_1", rowBytes))
				data.ScaleZ = float32(readValue(header, "scale_2", rowBytes))
				data.ColorR = cmn.EncodeSplatColor(readValue(header, "f_dc_0", rowBytes))
				data.ColorG = cmn.EncodeSplatColor(readValue(header, "f_dc_1", rowBytes))
				data.ColorB = cmn.EncodeSplatColor(readValue(header, "f_dc_2", rowBytes))
				data.ColorA = cmn.EncodeSplatOpacity(readValue(header, "opacity", rowBytes))
				data.RotationW = cmn.EncodeSplatRotation(readValue(header, "rot_0", rowBytes))
				data.RotationX = cmn.EncodeSplatRotation(readValue(header, "rot_1", rowBytes))
				data.RotationY = cmn.EncodeSplatRotation(readValue(header, "rot_2", rowBytes))
				data.RotationZ = cmn.EncodeSplatRotation(readValue(header, "rot_3", rowBytes))

				datas[i+j] = data

				shDim := 0
				maxShDegree := header.MaxShDegree()
				switch maxShDegree {
				case 1:
					shDim = 3
				case 2:
					shDim = 8
				case 3:
					shDim = 15
				}

				shs := make([]byte, 45)
				n := 0
				for k := 0; k < shDim; k++ {
					for c := range 3 {
						shs[n] = cmn.EncodeSplatSH(readValue(header, "f_rest_"+cmn.IntToString(k+c*shDim), rowBytes))
						n++
					}
				}
				for ; n < 45; n++ {
					shs[n] = 128 // cmn.EncodeSplatSH(0) = 128
				}

				switch maxShDegree {
				case 3:
					data.SH2 = shs[:24]
					data.SH3 = shs[24:]
				case 2:
					data.SH2 = shs[:24]
				case 1:
					data.SH1 = shs[:9]
				}
			}
			OnProgress(PhaseRead, i, header.VertexCount)
		}

	} else {
		// 压缩的 .compressed.ply
		readCompressedPlyDatas(file, header, datas)
	}

	OnProgress(PhaseRead, 100, 100)
	return header, datas
}

func readValue(header *PlyHeader, property string, splatDataBytes []byte) float64 {
	offset, typename := header.Property(property)
	switch typename {
	case "float":
		return float64(math.Float32frombits(binary.LittleEndian.Uint32(splatDataBytes[offset : offset+4]))) // 实际全是 float 类型
	case "double":
		var v float64
		cmn.ExitOnError(binary.Read(bytes.NewReader(splatDataBytes[offset:offset+8]), binary.LittleEndian, &v))
		return v
	case "int":
		var v int32
		cmn.ExitOnError(binary.Read(bytes.NewReader(splatDataBytes[offset:offset+4]), binary.LittleEndian, &v))
		return float64(v)
	case "uint":
		var v uint32
		cmn.ExitOnError(binary.Read(bytes.NewReader(splatDataBytes[offset:offset+4]), binary.LittleEndian, &v))
		return float64(v)
	case "short":
		var v int16
		cmn.ExitOnError(binary.Read(bytes.NewReader(splatDataBytes[offset:offset+2]), binary.LittleEndian, &v))
		return float64(v)
	case "ushort":
		var v uint16
		cmn.ExitOnError(binary.Read(bytes.NewReader(splatDataBytes[offset:offset+2]), binary.LittleEndian, &v))
		return float64(v)
	case "uchar":
		var v uint8
		cmn.ExitOnError(binary.Read(bytes.NewReader(splatDataBytes[offset:offset+1]), binary.LittleEndian, &v))
		return float64(v)
	}

	if cmn.Startwiths(property, "f_rest_") {
		return 0 // 球谐系数读取不到时，默认为0
	}

	log.Println("[Error] unsupported property:", "property", typename, property)
	os.Exit(1)
	return 0
}
