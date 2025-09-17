package cmn

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// 路径分隔符
func PathSeparator() string {
	return string(os.PathSeparator)
}

// 取文件名，如“abc.txt”
func FileName(name string) string {
	return path.Base(ReplaceAll(name, "\\", "/"))
}

// 取不含扩展名的文件名，如“abc.txt时返回abc”
func FileNameWithoutExt(name string) string {
	fileNameWithExtension := FileName(name)
	return strings.TrimSuffix(fileNameWithExtension, FileExtName(fileNameWithExtension))
}

// 取文件扩展名，如“.txt”
func FileExtName(name string) string {
	return path.Ext(name)
}

// 判断文件是否存在
func IsExistFile(file string) bool {
	defer func() {
		if err := recover(); err != nil {
			// 非法路径导致的panic，忽略掉，默认的返回值是false
		}
	}()
	s, err := os.Stat(file)
	if err == nil {
		return !s.IsDir()
	}
	if os.IsNotExist(err) {
		return false
	}
	return !s.IsDir()
}

// 判断文件夹是否存在
func IsExistDir(dir string) bool {
	defer func() {
		if err := recover(); err != nil {
			// 非法路径导致的panic，忽略掉，默认的返回值是false
		}
	}()
	s, err := os.Stat(dir)
	if err == nil {
		return s.IsDir()
	}
	if os.IsNotExist(err) {
		return false
	}
	return s.IsDir()
}

// 删除文件或目录(含全部子目录文件)
func RemoveAllFile(pathorfile string) error {
	return os.RemoveAll(pathorfile)
}

// 复制文件
func CopyFile(srcFilePath string, dstFilePath string) error {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	MkdirAll(Dir(dstFilePath))
	distFile, err := os.Create(dstFilePath)
	if err != nil {
		return err
	}
	defer distFile.Close()

	// 复制文件内容
	_, err = io.Copy(distFile, srcFile)
	if err != nil {
		return err
	}
	return nil

}

// 复制目录（源目录中的文件和子目录，复制到目标目录，目标目录不存在时自动创建）
func CopyDir(srcDir, dstDir string) error {
	// 创建目标目录
	err := MkdirAll(dstDir)
	if err != nil {
		return err
	}

	// 打开源目录
	dir, err := os.Open(srcDir)
	if err != nil {
		return err
	}
	defer dir.Close()

	// 读取源目录中的文件和子目录
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	// 逐个处理文件和子目录
	for _, fileInfo := range fileInfos {
		srcPath := filepath.Join(srcDir, fileInfo.Name())
		dstPath := filepath.Join(dstDir, fileInfo.Name())

		if fileInfo.IsDir() {
			// 如果是子目录，递归复制
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// 如果是文件，复制文件内容
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 返回目录，同filepath.Dir(path)
func Dir(path string) string {
	return filepath.Dir(path)
}

// 临时目录下创建临时子目录
func CreateTempDir() (string, error) {
	dir := filepath.Join(os.TempDir(), "gsbox", RandomString(16))
	return dir, MkdirAll(dir)
}

// 创建多级目录（存在时不报错）
func MkdirAll(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// 写文件（指定目录不存在时先创建，不含目录时存当前目录）
func WriteFileString(filename string, content string) error {
	return WriteFileBytes(filename, StringToBytes(content))
}

// 写文件（指定目录不存在时先创建，不含目录时存当前目录）
func WriteFileBytes(filename string, data []byte) error {
	os.MkdirAll(filepath.Dir(filename), 0777)
	return os.WriteFile(filename, data, 0666)
}

// 一次性读文件（适用于小文件）
func ReadFileBytes(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// 一次性读文件（适用于小文件）
func ReadFileString(filename string) (string, error) {
	by, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return BytesToString(by), nil
}

// 取目录中指定后缀的文件列表(升序)
func GetFiles(dir string, suffix string) ([]string, error) {
	var paths []string

	if !IsExistDir(dir) {
		return paths, nil
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return paths, err
	}

	for _, file := range files {
		if Endwiths(file.Name(), suffix) {
			path := filepath.Join(dir, file.Name())
			paths = append(paths, path)
		}
	}

	// 升序
	sort.Slice(paths, func(i, j int) bool {
		return paths[i] < paths[j]
	})

	return paths, nil
}
