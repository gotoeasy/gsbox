package cmn

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Trim(str string) string {
	return strings.TrimSpace(str)
}

// 字符串切割
func Split(str string, sep string) []string {
	return strings.Split(str, sep)
}

// []byte 转 string
func BytesToString(b []byte) string {
	return string(b)
}

// 判断是否包含（区分大小写）
func Contains(str string, substr string) bool {
	return strings.Contains(str, substr)
}

// 判断是否指定前缀
func Startwiths(str string, startstr string, ignoreCase ...bool) bool {
	lstr := Left(str, len(startstr))
	if len(ignoreCase) > 0 && ignoreCase[0] {
		return EqualsIngoreCase(lstr, startstr)
	}
	return lstr == startstr
}

func Endwiths(str string, endstr string, ignoreCase ...bool) bool {
	if len(ignoreCase) > 0 && ignoreCase[0] {
		return strings.HasSuffix(ToLower(str), ToLower(endstr))
	}
	return strings.HasSuffix(str, endstr)
}

// 取左文字
func Left(str string, length int) string {
	srune := []rune(str)
	lenr := len(srune)
	if lenr <= length {
		return str
	}

	var rs string
	for i := 0; i < length; i++ {
		rs += string(srune[i])
	}
	return rs
}

// 判断是否相同（忽略大小写）
func EqualsIngoreCase(str1 string, str2 string) bool {
	return ToLower(str1) == ToLower(str2)
}

// 转小写
func ToLower(str string) string {
	return strings.ToLower(str)
}

// string 转 int
func StringToInt(s string, defaultVal ...int) int {
	var defaultValue int
	if len(defaultVal) > 0 {
		defaultValue = defaultVal[0]
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return v
}

// 全部替换
func ReplaceAll(str string, old string, new string) string {
	return strings.ReplaceAll(str, old, new)
}

func ExitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// float32 转 []byte
func Float32ToBytes(f float32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, math.Float32bits(f))
	return b
}

// 限制范围
func Clip(f float64, min float64, max float64) float64 {
	if f < min {
		return min
	} else if f > max {
		return max
	}
	return f
}

// 限制范围
func ClipUint8(f float64) uint8 {
	if f < 0 {
		return 0
	} else if f > 255 {
		return 255
	}
	return uint8(f)
}

// 限制范围
func ClipFloat32(f float64) float32 {
	if f < -math.MaxFloat32 {
		return -math.MaxFloat32
	} else if f > math.MaxFloat32 {
		return math.MaxFloat32
	}
	return float32(f)
}

// 强制转换float64 -> float32，避免NaN
func ToFloat32(f float64) float32 {
	return ClipFloat32(f)
}

// 强制转换float64 -> uint8，避免NaN
func ToUint8(f float64) uint8 {
	return ClipUint8(f)
}

// 字符串数组拼接为字符串
func Join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

// int 转 string
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// 强制转换float64 -> float32，避免NaN，最后再转成[]byte
func ToFloat32Bytes(f float64) []byte {
	return Float32ToBytes(ToFloat32(f))
}

// 判断文件是否存在
func IsExistFile(file string) bool {
	defer func() {
		if err := recover(); err != nil {
			ExitOnError(errors.New("invalid file path: " + file))
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

// 返回目录，同filepath.Dir(path)
func Dir(file string) string {
	return filepath.Dir(file)
}

// 创建多级目录（存在时不报错）
func MkdirAll(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// 获取时间信息
func GetTimeInfo(milliseconds int64) string {
	seconds := milliseconds / 1000
	minutes := seconds / 60
	sMinutes := "minute"
	sSeconds := "second"
	sMilliseconds := "millisecond"

	if minutes > 1 {
		sMinutes += "s"
	}
	if seconds > 1 {
		sSeconds += "s"
	}
	if milliseconds > 1 {
		sMilliseconds += "s"
	}

	if minutes > 0 {
		seconds %= 60
		return fmt.Sprintf("%d %s %d %s", minutes, sMinutes, seconds, sSeconds)
	} else if seconds > 0 {
		ms := milliseconds % 1000
		return fmt.Sprintf("%d %s %d %s", seconds, sSeconds, ms, sMilliseconds)
	}

	return fmt.Sprintf("%d %s", milliseconds, sMilliseconds)
}
