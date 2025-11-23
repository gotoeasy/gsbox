package gsplat

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"gsbox/cmn"
	"io"
	"log"
	"os"
	"path/filepath"
)

func ReadSplat(splatFile string) []*SplatData {
	isNetFile := cmn.IsNetFile(splatFile)
	if isNetFile {
		tmpdir, err := cmn.CreateTempDir()
		cmn.ExitOnError(err)
		downloadFile := filepath.Join(tmpdir, cmn.FileName(splatFile))
		log.Println("[Info]", "download start,", splatFile)
		err = cmn.HttpDownload(splatFile, downloadFile, nil)
		cmn.RemoveAllFileIfError(err, tmpdir)
		cmn.ExitOnError(err)
		log.Println("[Info]", "download finish")
		splatFile = downloadFile
	}

	file, err := os.Open(splatFile)
	cmn.ExitOnError(err)
	defer func() {
		file.Close()
		if isNetFile {
			cmn.RemoveAllFile(cmn.Dir(splatFile))
		}
	}()

	fileInfo, err := file.Stat()
	cmn.ExitOnError(err)
	fileSize := fileInfo.Size()
	count := int(fileSize / SPLAT_DATA_SIZE)

	reader := bufio.NewReaderSize(file, 2*1024*1024)      // 2MB 缓冲区
	const batchSize = 4096                                // 每次读取的点数
	batchBytes := make([]byte, batchSize*SPLAT_DATA_SIZE) // 预分配缓冲区

	datas := make([]*SplatData, count)
	for i := 0; i < count; i += batchSize {
		OnProgress(PhaseRead, int(i), int(count))

		// 计算本次读取的点数
		currentBatchSize := batchSize
		if i+batchSize > count {
			currentBatchSize = int(count - i)
		}
		// 一次性读取一批数据
		readSize := currentBatchSize * SPLAT_DATA_SIZE
		_, err := io.ReadFull(reader, batchBytes[:readSize])
		cmn.ExitOnError(err)

		// 批量处理数据
		for j := 0; j < currentBatchSize; j++ {
			offset := j * SPLAT_DATA_SIZE
			splatBytes := batchBytes[offset : offset+SPLAT_DATA_SIZE]

			splatData := &SplatData{}
			cmn.ExitOnError(binary.Read(bytes.NewReader(splatBytes[0:4]), binary.LittleEndian, &splatData.PositionX))
			cmn.ExitOnError(binary.Read(bytes.NewReader(splatBytes[4:8]), binary.LittleEndian, &splatData.PositionY))
			cmn.ExitOnError(binary.Read(bytes.NewReader(splatBytes[8:12]), binary.LittleEndian, &splatData.PositionZ))
			cmn.ExitOnError(binary.Read(bytes.NewReader(splatBytes[12:16]), binary.LittleEndian, &splatData.ScaleX))
			cmn.ExitOnError(binary.Read(bytes.NewReader(splatBytes[16:20]), binary.LittleEndian, &splatData.ScaleY))
			cmn.ExitOnError(binary.Read(bytes.NewReader(splatBytes[20:24]), binary.LittleEndian, &splatData.ScaleZ))
			cmn.ExitOnError(binary.Read(bytes.NewReader(splatBytes[24:25]), binary.LittleEndian, &splatData.ColorR))
			cmn.ExitOnError(binary.Read(bytes.NewReader(splatBytes[25:26]), binary.LittleEndian, &splatData.ColorG))
			cmn.ExitOnError(binary.Read(bytes.NewReader(splatBytes[26:27]), binary.LittleEndian, &splatData.ColorB))
			cmn.ExitOnError(binary.Read(bytes.NewReader(splatBytes[27:28]), binary.LittleEndian, &splatData.ColorA))
			cmn.ExitOnError(binary.Read(bytes.NewReader(splatBytes[28:29]), binary.LittleEndian, &splatData.RotationW))
			cmn.ExitOnError(binary.Read(bytes.NewReader(splatBytes[29:30]), binary.LittleEndian, &splatData.RotationX))
			cmn.ExitOnError(binary.Read(bytes.NewReader(splatBytes[30:31]), binary.LittleEndian, &splatData.RotationY))
			cmn.ExitOnError(binary.Read(bytes.NewReader(splatBytes[31:32]), binary.LittleEndian, &splatData.RotationZ))

			splatData.ScaleX = cmn.DecodeSplatScale(splatData.ScaleX)
			splatData.ScaleY = cmn.DecodeSplatScale(splatData.ScaleY)
			splatData.ScaleZ = cmn.DecodeSplatScale(splatData.ScaleZ)

			datas[i+j] = splatData
		}
	}

	return datas
}
