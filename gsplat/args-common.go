package gsplat

import (
	"errors"
	"gsbox/cmn"
	"log"
)

var Args *cmn.OsArgs
var shDegreeFrom uint8
var oArg *ArgValues

type ArgValues struct {
	Quality int // 质量级别（1~9，默认5），越大越精确质量越好
	KI      int // 聚类计算时的迭代次数（5~50，默认10），越大越精确耗时越长
	KN      int // 聚类计算时的查找最邻近节点数量（10~200，默认15），越大越精确耗时越长

	webpQuality int // webp压缩的质量参数, 80~99，根据质量级别自动选定

	hasQuality bool
	hasKI      bool
	hasKN      bool
}

func InitArgs() *cmn.OsArgs {
	Args = cmn.ParseArgs("-v", "-version", "--version", "-h", "-help", "--help")

	oArg = &ArgValues{}
	oArg.Quality = max(1, min(cmn.StringToInt(Args.GetArgIgnorecase("-q", "--quality"), 5), 9))
	oArg.KI = max(5, min(cmn.StringToInt(Args.GetArgIgnorecase("-ki", "--kmeans-iterations"), 10), 50))
	oArg.KN = max(10, min(cmn.StringToInt(Args.GetArgIgnorecase("-kn", "--kmeans-nearest-nodes"), 15), 200))

	oArg.hasQuality = Args.HasArgIgnorecase("-q", "--quality")
	oArg.hasKI = Args.HasArgIgnorecase("-ki", "--kmeans-iterations")
	oArg.hasKN = Args.HasArgIgnorecase("-kn", "--kmeans-nearest-nodes")

	// 按摩质量级别自动调整相关参数
	kis := []int{5, 7, 9, 10, 10, 10, 12, 15, 20}
	kns := []int{10, 12, 14, 15, 15, 20, 30, 50, 100}
	if oArg.hasQuality && !oArg.hasKI {
		oArg.KI = kis[oArg.Quality-1]
	}
	if oArg.hasQuality && !oArg.hasKN {
		oArg.KN = kns[oArg.Quality-1]
	}
	wqs := []int{80, 84, 86, 88, 90, 92, 94, 96, 99}
	oArg.webpQuality = wqs[oArg.Quality-1]

	return Args
}

func SetShDegreeFrom(shDegree uint8) {
	shDegreeFrom = shDegree
}

func GetAndCheckInputFiles() []string {
	inputs := Args.GetArgsIgnorecase("-i", "--input")
	for i := range inputs {
		cmn.ExitOnConditionError(inputs[i] == "", errors.New(`please specify the input file`))
		if !cmn.Startwiths(inputs[i], "http://") && !cmn.Startwiths(inputs[i], "https://") {
			cmn.ExitOnConditionError(!cmn.IsExistFile(inputs[i]), errors.New("file not found: "+inputs[i]))
		}
	}
	return inputs
}

func GetAndCheckInputFile() string {
	input := Args.GetArgIgnorecase("-i", "--input")
	cmn.ExitOnConditionError(input == "", errors.New(`please specify the input file`))
	if !cmn.IsNetFile(input) {
		cmn.ExitOnConditionError(!cmn.IsExistFile(input), errors.New("file not found: "+input))
	}
	return input
}

func CreateOutputDir() string {
	output := Args.GetArgIgnorecase("-o", "--output")
	cmn.ExitOnConditionError(output == "", errors.New(`please specify the output file`))
	cmn.ExitOnError(cmn.MkdirAll(cmn.Dir(output)))
	return output
}

func GetArgShDegree() uint8 {
	shDegree := shDegreeFrom
	if Args.HasArg("-sh", "--shDegree") {
		sh := cmn.StringToInt(Args.GetArgIgnorecase("-sh", "--shDegree"), 3) // 默认满级输出
		sh = max(0, min(sh, 3))                                              // 限制越界参数值到边界
		shDegree = min(shDegreeFrom, uint8(sh))                              // 来源和输出取其小，不必输出更高级别
	}
	return shDegree
}

func GetArgFlag(defaultFlag uint8, arg1 string, arg2 string) uint8 {
	flag := defaultFlag
	if Args.HasArg(arg1, arg2) {
		val := cmn.StringToInt(Args.GetArgIgnorecase(arg1, arg2), -1)
		if val >= 0 && val < 256 {
			flag = uint8(val)
		}
	}
	return flag
}

func GetRotateArgs() (bool, float32, float32, float32) {
	has := Args.HasArgIgnorecase("-rx", "--rotateX", "-ry", "--rotateY", "-rz", "--rotateZ")
	var rx, ry, rz float32
	if has {
		rx = cmn.StringToFloat32(Args.GetArgIgnorecase("-rx", "--rotateX"), 0)
		ry = cmn.StringToFloat32(Args.GetArgIgnorecase("-ry", "--rotateY"), 0)
		rz = cmn.StringToFloat32(Args.GetArgIgnorecase("-rz", "--rotateZ"), 0)
	}
	return has, rx, ry, rz
}

func GetScaleArgs() (bool, float32) {
	has := Args.HasArgIgnorecase("-s", "--scale")
	var scale float32 = 1.0
	if has {
		scale = min(max(cmn.StringToFloat32(Args.GetArgIgnorecase("-s", "--scale"), 1.0), 0.001), 1000.0)
	}
	return has, scale
}

func GetTranslateArgs() (bool, float32, float32, float32) {
	has := Args.HasArgIgnorecase("-tx", "--translateX", "-ty", "--translateY", "-tz", "--translateZ")
	var tx, ty, tz float32
	if has {
		tx = cmn.StringToFloat32(Args.GetArgIgnorecase("-tx", "--translateX"), 0)
		ty = cmn.StringToFloat32(Args.GetArgIgnorecase("-ty", "--translateY"), 0)
		tz = cmn.StringToFloat32(Args.GetArgIgnorecase("-tz", "--translateZ"), 0)
	}
	return has, tx, ty, tz
}

func GetArgFlagValue(arg1 string, arg2 string) uint16 {
	flag := uint16(0)
	if Args.HasArg(arg1, arg2) {
		val := cmn.StringToInt(Args.GetArgIgnorecase(arg1, arg2), -1)
		if val > 0 && val < 32768 {
			flag = uint16(val)
		} else {
			log.Println("[Warn] ignore invalid flag value:", val)
		}
	}
	return flag
}
