package cmn

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
)

// 打包成一个sog文件
func ZipSogFiles(zipAsSogFile string, files []string, printLogs ...bool) {
	printLog := len(printLogs) == 0 || printLogs[0]
	if printLog {
		log.Println("[Info] save as sog")
	}
	os.MkdirAll(filepath.Dir(zipAsSogFile), 0666)

	zipFile, err := os.Create(zipAsSogFile)
	ExitOnError(err)
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		addFileToZip(zipWriter, file)
	}
}

func addFileToZip(zipWriter *zip.Writer, filePath string) {
	file, err := os.Open(filePath)
	ExitOnError(err)
	defer file.Close()

	info, err := file.Stat()
	ExitOnError(err)

	// 只保留文件名，去掉目录路径
	header, err := zip.FileInfoHeader(info)
	ExitOnError(err)
	header.Name = filepath.Base(filePath)

	writer, err := zipWriter.CreateHeader(header)
	ExitOnError(err)

	// 复制文件内容到ZIP
	_, err = io.Copy(writer, file)
	ExitOnError(err)
}
