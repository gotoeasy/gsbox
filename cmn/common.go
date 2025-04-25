package cmn

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const lastestVerUrl = "https://reall3d.com/gsbox/open-lastest.json?v=" + VER

var NewVersionMessage = ""

const COLOR_SCALE = 0.15
const SH_C0 float64 = 0.28209479177387814

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

func ExitOnConditionError(condition bool, err error) {
	if condition {
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

// float32 转 3字节长度的[]byte
func EncodeFloat32ToBytes3(f float32) []byte {
	fixed32 := int32(math.Round(float64(f) * 4096))

	// 将固定点数拆分为3字节
	return []byte{
		byte(fixed32 & 0xFF),         // 最低字节
		byte((fixed32 >> 8) & 0xFF),  // 中间字节
		byte((fixed32 >> 16) & 0xFF), // 最高字节
	}
}

// 3字节长度的[]byte 转 float32
func DecodeBytes3ToFloat32(bytes []byte) float32 {
	fixed32 := int32(bytes[0]) | int32(bytes[1])<<8 | int32(bytes[2])<<16
	if fixed32&0x800000 != 0 {
		fixed32 |= int32(-1) << 24 // 如果符号位为1，将高8位填充为1
	}
	return float32(fixed32) / 4096 // 将固定点数转换回浮点数
}

// float32 编码成 byte
func EncodeFloat32ToByte(f float32) byte {
	if f <= 0 {
		return 0
	}
	encoded := math.Round((math.Log(float64(f)) + 10.0) * 16.0) // 编码公式
	// 确保结果在0-255范围内
	if encoded < 0 {
		return 0
	} else if encoded > 255 {
		return 255
	}
	return byte(encoded)

}

// byte 解码成 float32
func DecodeByteToFloat32(encodedByte byte) float32 {
	return float32(math.Exp(float64(encodedByte)/16.0 - 10.0)) // 解码公式
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

// uint32 转 string
func Uint32ToString(num uint32) string {
	return strconv.FormatUint(uint64(num), 10)
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

/** 删除非ASCII字符，不可见字符替换为空格 */
func RemoveNonASCII(s string) (bool, string) {
	var result string
	remove := false
	for _, r := range s {
		if r <= 127 {
			if r == '\t' || r == '\n' || r == '\r' || r == '\f' || r == '\v' {
				result += " "
			} else {
				result += string(r)
			}
		} else {
			remove = true
		}
	}
	return remove, result
}

func GetSystemDateYYYYMMDD() uint32 {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	day := now.Day()

	date := uint32(year*10000 + month*100 + day) // 组合成 yyyymmdd 格式
	return date
}

func BytesToInt32(bs []byte) int32 {
	return int32(binary.LittleEndian.Uint32(bs))
}

func BytesToFloat32(bs []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(bs))
}

func BytesToUint32(bs []byte) uint32 {
	return binary.LittleEndian.Uint32(bs)
}

func Int32ToBytes(intNum int32) []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.LittleEndian, intNum)
	return bytebuf.Bytes()
}

func StringToBytes(s string) []byte {
	return []byte(s)
}

func Uint32ToBytes(intNum uint32) []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.LittleEndian, intNum)
	return bytebuf.Bytes()
}

/*
*【注意】要用于专有校验时，应修改初始值或添加自定义的前缀后缀参与计算，且不公开
 */
func HashBytes(bts []byte) uint32 {
	var rs uint32 = 53653
	for i := 0; i < len(bts); i++ {
		rs = (rs * 33) ^ uint32(bts[i])
	}
	return rs
}

func init() {
	req, err := http.NewRequest("GET", lastestVerUrl, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{Timeout: 2 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	var data struct {
		Ver string `json:"ver"`
	}
	err = json.Unmarshal(bts, &data)
	if err != nil {
		return
	}

	if data.Ver != VER {
		NewVersionMessage = "\nNotice: the latest version (" + data.Ver + ") is now available.\n"
	}
}

func ClipUint8Round(x float64) uint8 {
	return uint8(math.Max(0, math.Min(255, math.Round(x))))
}

// 固定24位编码
func SpzEncodePosition(val float32) []byte {
	return EncodeFloat32ToBytes3(val)
}

func SpzDecodePosition(bts []byte, fractionalBits uint8) float32 {
	scale := 1.0 / float64(int(1<<fractionalBits))
	fixed32 := int32(bts[0]) | int32(bts[1])<<8 | int32(bts[2])<<16
	if fixed32&0x800000 != 0 {
		// fixed32 |= int32(-1) << 24 // 如果符号位为1，将高8位填充为1
		fixed32 |= -0x1000000
	}
	return ClipFloat32(float64(fixed32) * scale)
}

func SpzEncodeScale(val float32) uint8 {
	return ClipUint8Round((float64(val) + 10.0) * 16.0)
}

func SpzDecodeScale(val uint8) float32 {
	return float32(val)/16.0 - 10.0
}

func SpzEncodeRotations(rw uint8, rx uint8, ry uint8, rz uint8) []byte {
	r0 := float64(rw)/128.0 - 1.0
	r1 := float64(rx)/128.0 - 1.0
	r2 := float64(ry)/128.0 - 1.0
	r3 := float64(rz)/128.0 - 1.0
	if r0 < 0 {
		r0, r1, r2, r3 = -r0, -r1, -r2, -r3
	}
	qlen := math.Sqrt(r0*r0 + r1*r1 + r2*r2 + r3*r3)
	return []byte{ClipUint8Round((r1/qlen)*127.5 + 127.5), ClipUint8Round((r2/qlen)*127.5 + 127.5), ClipUint8Round((r3/qlen)*127.5 + 127.5)}
}

func SpzDecodeRotations(rx uint8, ry uint8, rz uint8) (uint8, uint8, uint8, uint8) {
	r1 := float64(rx)/127.5 - 1.0
	r2 := float64(ry)/127.5 - 1.0
	r3 := float64(rz)/127.5 - 1.0
	r0 := math.Sqrt(math.Max(0.0, 1.0-(r1*r1+r2*r2+r3*r3)))
	return ClipUint8(r0*128.0 + 128.0), ClipUint8(r1*128.0 + 128.0), ClipUint8(r2*128.0 + 128.0), ClipUint8(r3*128.0 + 128.0)
}

func SpzEncodeColor(val uint8) uint8 {
	fColor := (float64(val)/255.0 - 0.5) / SH_C0                      // 解码为原值
	return ClipUint8Round(fColor*(COLOR_SCALE*255.0) + (0.5 * 255.0)) // 按spz方式编码
}

func SpzDecodeColor(val uint8) uint8 {
	fColor := (float64(val) - (0.5 * 255.0)) / (COLOR_SCALE * 255.0)
	return ClipUint8((0.5 + SH_C0*fColor) * 255.0)
}

func SpzEncodeSH1(encodeSHval uint8) uint8 {
	q := math.Floor((float64(encodeSHval)+4.0)/8.0) * 8.0
	return ClipUint8(q)
}

func SpzEncodeSH23(encodeSHval uint8) uint8 {
	q := math.Floor((float64(encodeSHval)+8.0)/16.0) * 16.0
	return ClipUint8(q)
}

// ------------------ Splat ------------------
func EncodeSplatScale(val float32) float32 {
	return ClipFloat32(math.Exp(float64(val)))
}
func DecodeSplatScale(encodedVal float32) float32 {
	return ClipFloat32(math.Log(float64(encodedVal)))
}

func EncodeSplatColor(val float64) uint8 {
	return ClipUint8((0.5 + SH_C0*val) * 255.0)
}
func DecodeSplatColor(val uint8) float32 {
	return ClipFloat32((float64(val)/255.0 - 0.5) / SH_C0)
}

func EncodeSplatOpacity(val float64) uint8 {
	return ClipUint8((1.0 / (1.0 + math.Exp(-val))) * 255.0)
}
func DecodeSplatOpacity(val uint8) float32 {
	return ClipFloat32(-math.Log((1.0 / (float64(val) / 255.0)) - 1.0))
}

func EncodeSplatRotation(val float64) uint8 {
	return ClipUint8(val*128.0 + 128.0)
}
func DecodeSplatRotation(val uint8) float32 {
	return (float32(val) - 128.0) / 128.0
}

func EncodeSplatSH(val float64) uint8 {
	return ClipUint8(math.Round(val*128.0) + 128.0)
}
func DecodeSplatSH(val uint8) float32 {
	return (float32(val) - 128.0) / 128.0
}

// ------------------ SPX ------------------
func EncodeSpxPositionUint24(val float32) []byte {
	fixed32 := int32(math.Round(float64(val) * 4096.0))
	return []byte{
		byte(fixed32 & 0xFF),         // 最低字节
		byte((fixed32 >> 8) & 0xFF),  // 中间字节
		byte((fixed32 >> 16) & 0xFF), // 最高字节
	}
}
func DecodeSpxPositionUint24(b0 uint8, b1 uint8, b2 uint8) float32 {
	i32 := int32(b0) | (int32(b1) << 8) | (int32(b2) << 16)
	if i32&0x800000 > 0 {
		i32 |= -0x1000000
	}
	return float32(i32) / 4096.0
}

func EncodeSpxScale(val float32) uint8 {
	return ClipUint8Round((float64(val) + 10.0) * 16.0)
}

func DecodeSpxScale(val uint8) float32 {
	return float32(val)/16.0 - 10.0
}

func EncodeSpxSH(encodeSHval uint8) uint8 {
	q := math.Floor((float64(encodeSHval)+4.0)/8.0) * 8.0
	return ClipUint8(q)
}
