package cmn

import (
	"os"
)

// 命令行参数解析结果
type OsArgs struct {
	String        string // 原命令
	mapIndexValue map[int]string
	mapCmd        map[string]int
	mapParam      map[string]string
	LastParam     string // 最后一个参数
}

// 解析命令行参数
func ParseArgs() *OsArgs {
	args := &OsArgs{}
	args.mapIndexValue = make(map[int]string)
	args.mapCmd = make(map[string]int)
	args.mapParam = make(map[string]string)
	args.String = Join(os.Args, " ")

	for index, arg := range os.Args {
		if index == 0 {
			continue // 跳过命令本身
		}
		args.mapIndexValue[index] = arg
		args.LastParam = arg // 最后一个参数作为命令的输入看待
	}

	for index, arg := range os.Args {
		if index == 0 {
			continue // 跳过命令本身
		}
		if index == 1 {
			if !Startwiths(arg, "-") {
				args.mapCmd[arg] = index          // 非“-”前缀，是命令
				args.mapCmd[ToLower(arg)] = index // 转小写存一份
			}
			continue // 不可能是参数值
		}

		if Startwiths(args.mapIndexValue[index-1], "-") {
			// 上一个参数是参数，则当前参数是参数值
			val := ReplaceAll(arg, "\\r", "\r")
			val = ReplaceAll(val, "\\n", "\n")
			val = ReplaceAll(val, "\\t", "\t")
			args.mapParam[args.mapIndexValue[index-1]] = val
			args.mapParam["\n"+ToLower(args.mapIndexValue[index-1])] = val
		} else {
			// 上一个参数是命令，当前参数可能是命令或参数
			if !Startwiths(arg, "-") {
				args.mapCmd[arg] = index          // 非“-”前缀，是命令
				args.mapCmd[ToLower(arg)] = index // 转小写存一份
			}
		}
	}

	return args
}

// 取指定参数对应的值，参数总是“-”前缀，总是紧跟参数值
// 例如命令 test -d /abc 用GetArg("-d", "--dir")取得/abc
func (o *OsArgs) GetArg(names ...string) string {
	for i := 0; i < len(names); i++ {
		if o.mapParam[names[i]] != "" {
			return o.mapParam[names[i]]
		}
	}
	return ""
}
func (o *OsArgs) GetArgIgnorecase(names ...string) string {
	for i := 0; i < len(names); i++ {
		v := o.mapParam["\n"+ToLower(names[i])]
		if v != "" {
			return v
		}
	}
	return ""
}

// 判断是否含有指定参数，参数总是“-”前缀，总是紧跟参数值
func (o *OsArgs) HasArg(names ...string) bool {
	return o.GetArg(names...) != ""
}
func (o *OsArgs) HasArgIgnorecase(names ...string) bool {
	return o.GetArgIgnorecase(names...) != ""
}

// 判断是否含有指定命令，命令总是非“-”前缀，总是单个字符串
// 例如命令 docker run ... HasCmd("run")返回true
func (o *OsArgs) HasCmd(names ...string) bool {
	for i := 0; i < len(names); i++ {
		if o.mapCmd[names[i]] > 0 {
			return true
		}
	}
	return false
}

// 判断是否含有指定命令(忽略大小写)，命令总是非“-”前缀，总是单个字符串
// 例如命令 docker run ... HasCmd("Run")返回true
func (o *OsArgs) HasCmdIgnorecase(names ...string) bool {
	for i := 0; i < len(names); i++ {
		if o.mapCmd[ToLower(names[i])] > 0 {
			return true
		}
	}
	return false
}
