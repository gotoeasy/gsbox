package cmn

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// 解压 ZIP 文件到指定目录
func Unzip(zipFilePath string, extractDir string) {
	// 打开 ZIP 文件
	reader, err := zip.OpenReader(zipFilePath)
	ExitOnError(err)
	defer reader.Close()

	// 创建目标目录
	err = MkdirAll(extractDir)
	ExitOnError(err)

	// 遍历 ZIP 文件中的条目
	for _, file := range reader.File {
		// 获取文件的完整路径
		filePath := filepath.Join(extractDir, file.Name)

		// 如果是目录，直接创建目录
		if file.FileInfo().IsDir() {
			err = MkdirAll(filePath)
			ExitOnError(err)
			continue
		}

		// 创建文件
		fileWriter, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		ExitOnError(err)
		defer fileWriter.Close()

		// 打开 ZIP 文件中的文件
		fileReader, err := file.Open()
		ExitOnError(err)
		defer fileReader.Close()

		// 将 ZIP 文件中的文件内容复制到目标文件
		_, err = io.Copy(fileWriter, fileReader)
		ExitOnError(err)
	}
}
