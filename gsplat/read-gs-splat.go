package gsplat

import (
	"bytes"
	"encoding/binary"
	"gsbox/cmn"
	"os"
)

func ReadSplat(splatFile string, headers ...string) []*SplatData {

	file, err := os.Open(splatFile)
	cmn.ExitOnError(err)
	defer file.Close()

	var headerSize int64 = 0
	if len(headers) > 0 && hasHeader(file, headers[0]) {
		headerSize = int64(len((headers[0])))
	}

	fileInfo, err := file.Stat()
	cmn.ExitOnError(err)
	fileSize := fileInfo.Size()
	count := (fileSize - headerSize) / SPLAT_DATA_SIZE

	var i int64 = 0
	datas := make([]*SplatData, count)
	for ; i < count; i++ {
		splatBytes := make([]byte, SPLAT_DATA_SIZE)
		_, err = file.ReadAt(splatBytes, headerSize+i*SPLAT_DATA_SIZE)
		cmn.ExitOnError(err)

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
		datas[i] = splatData
	}

	return datas
}

func hasHeader(file *os.File, header string) bool {
	if header == "" {
		return false
	}

	bs := make([]byte, 1024)
	_, err := file.Read(bs)
	cmn.ExitOnError(err)

	str := cmn.BytesToString(bs)
	return cmn.Startwiths(str, header)
}
