package cmn

import (
	"os"
)

// 命令行解析结果
type OsArgs struct {
	String        string // 原命令
	ArgCount      int    // 参数个数(含命令本身)
	mapIndexValue map[int]string
	mapCmd        map[string]bool
	mapParam      map[string]string
	LastParam     string // 最后一个参数
	mapCustomCmd  map[string]bool
}

// 命令行解析器
// 约定：
// 参数名总是以“-”作为前缀，参数值紧跟参数名，不支持重复参数名
// 指令默认总是非“-”前缀，但也可以通过参数自定义指令，指令总是忽略大小写
func ParseArgs(customCmds ...string) *OsArgs {
	args := &OsArgs{}
	args.mapIndexValue = make(map[int]string)
	args.mapCmd = make(map[string]bool)
	args.mapParam = make(map[string]string)
	args.mapCustomCmd = make(map[string]bool)
	args.String = Join(os.Args, " ")
	args.ArgCount = len(os.Args)

	for _, cmd := range customCmds {
		args.mapCustomCmd[ToLower(Trim(cmd))] = true
	}

	for index, arg := range os.Args {
		args.mapIndexValue[index] = arg
		args.LastParam = arg // 最后一个参数作为命令的输入看待
	}

	for index, arg := range os.Args {
		if index == 0 {
			continue // 跳过命令本身
		}
		if index == 1 {
			if !Startwiths(arg, "-") || args.mapCustomCmd[ToLower(arg)] {
				args.mapCmd[ToLower(arg)] = true // 是指令
			}
			continue // 不可能是参数值
		}

		if !Startwiths(args.mapIndexValue[index-1], "-") || args.mapCustomCmd[ToLower(args.mapIndexValue[index-1])] {
			// 上一个参数是指令，当前参数可能是指令或参数
			if !Startwiths(arg, "-") || args.mapCustomCmd[ToLower(arg)] {
				args.mapCmd[ToLower(arg)] = true // 是指令
			}
		} else {
			// 上一个参数是参数，则当前参数是参数值
			val := ReplaceAll(arg, "\\r", "\r")
			val = ReplaceAll(val, "\\n", "\n")
			val = ReplaceAll(val, "\\t", "\t")
			args.mapParam[args.mapIndexValue[index-1]] = val
			args.mapParam["\n"+ToLower(args.mapIndexValue[index-1])] = val
		}
	}

	return args
}

// 取指定参数名对应的值
// 例如命令 test -d /abc 用GetArg("-d", "--dir")取得/abc
func (o *OsArgs) GetArg(names ...string) string {
	for i := 0; i < len(names); i++ {
		if o.mapParam[names[i]] != "" {
			return o.mapParam[names[i]]
		}
	}
	return ""
}

// 取指定参数名对应的值(忽略参数名大小写)
func (o *OsArgs) GetArgIgnorecase(names ...string) string {
	for i := 0; i < len(names); i++ {
		v := o.mapParam["\n"+ToLower(names[i])]
		if v != "" {
			return v
		}
	}
	return ""
}

// 判断是否含有指定参数名
func (o *OsArgs) HasArg(names ...string) bool {
	return o.GetArg(names...) != ""
}

// 判断是否含有指定参数名(忽略大小写)
func (o *OsArgs) HasArgIgnorecase(names ...string) bool {
	return o.GetArgIgnorecase(names...) != ""
}

// 判断是否含有指定指令(忽略大小写)
// 例如命令 docker run ... HasCmd("Run")返回true
func (o *OsArgs) HasCmd(names ...string) bool {
	for i := 0; i < len(names); i++ {
		if o.mapCmd[ToLower(Trim(names[i]))] {
			return true
		}
	}
	return false
}

func (o *OsArgs) GetArgByIndex(index int) string {
	return o.mapIndexValue[index]
}
