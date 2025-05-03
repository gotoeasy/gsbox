package main

import (
	"errors"
	"fmt"
	"gsbox/cmn"
	"gsbox/gsplat"
	"log"
	"os"
	"time"
)

func main() {
	args := cmn.ParseArgs("-v", "-version", "--version", "-h", "-help", "--help")
	gsplat.Args = args
	if args.HasCmd("-v", "-version", "--version") && args.ArgCount == 2 {
		version()
	} else if args.HasCmd("-h", "-help", "--help") && args.ArgCount == 2 {
		usage()
	} else if args.HasCmd("p2s", "ply2splat") {
		ply2splat(args)
	} else if args.HasCmd("p2x", "ply2spx") {
		ply2spx(args)
	} else if args.HasCmd("p2z", "ply2spz") {
		ply2spz(args)
	} else if args.HasCmd("p2p", "ply2ply") {
		ply2ply(args)
	} else if args.HasCmd("s2p", "splat2ply") {
		splat2ply(args)
	} else if args.HasCmd("s2x", "splat2spx") {
		splat2spx(args)
	} else if args.HasCmd("s2z", "splat2spz") {
		splat2spz(args)
	} else if args.HasCmd("s2s", "splat2splat") {
		splat2splat(args)
	} else if args.HasCmd("x2p", "spx2ply") {
		spx2ply(args)
	} else if args.HasCmd("x2s", "spx2splat") {
		spx2splat(args)
	} else if args.HasCmd("x2z", "spx2spz") {
		spx2spz(args)
	} else if args.HasCmd("x2x", "spx2spx") {
		spx2spx(args)
	} else if args.HasCmd("z2p", "spz2ply") {
		spz2ply(args)
	} else if args.HasCmd("z2s", "spz2splat") {
		spz2splat(args)
	} else if args.HasCmd("z2x", "spz2spx") {
		spz2spx(args)
	} else if args.HasCmd("z2z", "spz2spz") {
		spz2spz(args)
	} else if args.HasCmd("info") {
		plyInfo(args)
	} else {
		usage()
	}
	fmt.Print(cmn.NewVersionMessage)
}

func version() {
	fmt.Println("")
	fmt.Println("gsbox", cmn.VER)
	fmt.Println("homepage", "https://github.com/gotoeasy/gsbox")
}
func usage() {
	fmt.Println("")
	fmt.Println("gsbox", cmn.VER)
	fmt.Println("homepage", "https://github.com/gotoeasy/gsbox")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  gsbox [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  p2s, ply2splat            convert ply to splat")
	fmt.Println("  p2x, ply2spx              convert ply to spx")
	fmt.Println("  p2z, ply2spz              convert ply to spz")
	fmt.Println("  p2p, ply2ply              convert ply to ply")
	fmt.Println("  s2p, splat2ply            convert splat to ply")
	fmt.Println("  s2x, splat2spx            convert splat to spx")
	fmt.Println("  s2z, splat2spz            convert splat to spz")
	fmt.Println("  s2s, splat2splat          convert splat to splat")
	fmt.Println("  x2p, spx2ply              convert spx to ply")
	fmt.Println("  x2s, spx2splat            convert spx to splat")
	fmt.Println("  x2z, spx2spz              convert spx to spz")
	fmt.Println("  x2x, spx2spx              convert spx to spx")
	fmt.Println("  z2p, spz2ply              convert spz to ply")
	fmt.Println("  z2s, spz2splat            convert spz to splat")
	fmt.Println("  z2x, spz2spx              convert spz to spx")
	fmt.Println("  z2z, spz2spz              convert spz to spz")
	fmt.Println("  info <file>               display the model file information")
	fmt.Println("  -i,  --input <file>       specify the input file")
	fmt.Println("  -o,  --output <file>      specify the output file")
	fmt.Println("  -c,  --comment <text>     specify the comment for ply/spx output")
	fmt.Println("  -bs, --block-size <num>   specify the block size for spx output (default 20480)")
	fmt.Println("  -sh, --shDegree <num>     specify the SH degree for ply/spx/spz output")
	fmt.Println("  -f1, --flag1 <num>        specify the header flag1 for spx output")
	fmt.Println("  -f2, --flag2 <num>        specify the header flag2 for spx output")
	fmt.Println("  -f3, --flag3 <num>        specify the header flag3 for spx output")
	fmt.Println("  -rx, --rotateX <num>      specify the rotation angle in degrees about the x-axis for transform")
	fmt.Println("  -ry, --rotateY <num>      specify the rotation angle in degrees about the y-axis for transform")
	fmt.Println("  -rz, --rotateZ <num>      specify the rotation angle in degrees about the z-axis for transform")
	fmt.Println("  -s,  --scale <num>        specify a uniform scaling factor(0.01~100.0) for transform")
	fmt.Println("  -tx, --translateX <num>   specify the translation value about the x-axis for transform")
	fmt.Println("  -ty, --translateY <num>   specify the translation value about the y-axis for transform")
	fmt.Println("  -tz, --translateZ <num>   specify the translation value about the z-axis for transform")
	fmt.Println("  -v,  --version            display version information")
	fmt.Println("  -h,  --help               display help information")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  gsbox ply2splat -i /path/to/input.ply -o /path/to/output.splat")
	fmt.Println("  gsbox s2x -i /path/to/input.splat -o /path/to/output.spx -c \"your comment\" -bs 10240")
	fmt.Println("  gsbox x2z -i /path/to/input.spx -o /path/to/output.spz -sh 0 -rz 90 -s 0.9 -tx 0.1")
	fmt.Println("  gsbox z2p -i /path/to/input.spz -o /path/to/output.ply -c \"your comment\"")
	fmt.Println("  gsbox info -i /path/to/file.spx")
	fmt.Println("")
}
func plyInfo(args *cmn.OsArgs) {
	// info
	input := args.GetArgIgnorecase("-i", "--input")
	if args.ArgCount == 3 {
		input = args.LastParam
		if cmn.EqualsIngoreCase(input, "info") {
			input = args.GetArgByIndex(1)
		}
	}

	if input == "" {
		cmn.ExitOnError(errors.New("please specify the input ply file"))
	}
	if !cmn.IsExistFile(input) {
		cmn.ExitOnError(errors.New("file not found: " + input))
	}

	isPly := cmn.Endwiths(input, ".ply", true)
	isSpx := cmn.Endwiths(input, ".spx", true)
	isSplat := cmn.Endwiths(input, ".splat", true)
	isSpz := cmn.Endwiths(input, ".spz", true)

	if isPly {
		header, err := gsplat.ReadPlyHeaderString(input, 1024)
		cmn.ExitOnError(err)
		fmt.Print(header)
	} else if isSpx {
		header := gsplat.ParseSpxHeader(input)
		fmt.Println(header.ToStringSpx())
	} else if isSpz {
		header, _ := gsplat.ReadSpz(input, true)
		fmt.Println(header.ToString())
	} else if isSplat {
		fileInfo, err := os.Stat(input)
		cmn.ExitOnError(err)
		fileSize := fileInfo.Size()
		if (fileSize)%32 > 0 {
			cmn.ExitOnError(errors.New("invalid splat format"))
		} else {
			fmt.Println("SplatCount :", fileSize/32)
		}
	} else {
		cmn.ExitOnError(errors.New("the input file must be (ply | splat | spx) format"))
	}

}

func ply2splat(args *cmn.OsArgs) {
	log.Println("Convert ply to splat.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	_, datas := gsplat.ReadPly(input, "ply-3dgs")
	gsplat.RstTransformDatas(datas)
	gsplat.Sort(datas)
	gsplat.WriteSplat(output, datas)
	log.Println(input, " --> ", output)
	log.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ply2spx(args *cmn.OsArgs) {
	log.Println("Convert ply to spx.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadPly(input, "ply-3dgs")
	gsplat.RstTransformDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(header.MaxShDegree(), args)
	flag1 := getArgFlag(0, args, "-f1", "--flag1")
	flag2 := getArgFlag(0, args, "-f2", "--flag2")
	flag3 := getArgFlag(0, args, "-f3", "--flag3")
	gsplat.WriteSpx(output, datas, comment, shDegree, flag1, flag2, flag3)
	log.Println(input, " --> ", output)
	log.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ply2spz(args *cmn.OsArgs) {
	log.Println("Convert ply to spz.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadPly(input, "ply-3dgs")
	gsplat.RstTransformDatas(datas)
	gsplat.Sort(datas)
	shDegree := getArgShDegree(header.MaxShDegree(), args)
	gsplat.WriteSpz(output, datas, shDegree)
	log.Println(input, " --> ", output)
	log.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ply2ply(args *cmn.OsArgs) {
	log.Println("Convert ply to ply.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadPly(input, "ply-3dgs")
	gsplat.RstTransformDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(header.MaxShDegree(), args)
	gsplat.WritePly(output, datas, comment, shDegree)
	log.Println(input, " --> ", output)
	log.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func splat2ply(args *cmn.OsArgs) {
	log.Println("Convert splat to ply.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	datas := gsplat.ReadSplat(input)
	gsplat.RstTransformDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(0, args)
	gsplat.WritePly(output, datas, comment, shDegree)
	log.Println(input, " --> ", output)
	log.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}
func splat2spx(args *cmn.OsArgs) {
	log.Println("Convert splat to spx.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	datas := gsplat.ReadSplat(input)
	gsplat.RstTransformDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(0, args)
	flag1 := getArgFlag(0, args, "-f1", "--flag1")
	flag2 := getArgFlag(0, args, "-f2", "--flag2")
	flag3 := getArgFlag(0, args, "-f3", "--flag3")
	gsplat.WriteSpx(output, datas, comment, shDegree, flag1, flag2, flag3)
	log.Println(input, " --> ", output)
	log.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}
func splat2spz(args *cmn.OsArgs) {
	log.Println("Convert splat to spz.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	datas := gsplat.ReadSplat(input)
	gsplat.RstTransformDatas(datas)
	gsplat.Sort(datas)
	gsplat.WriteSpz(output, datas, 0)
	log.Println(input, " --> ", output)
	log.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}
func splat2splat(args *cmn.OsArgs) {
	log.Println("Convert splat to splat.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	datas := gsplat.ReadSplat(input)
	gsplat.RstTransformDatas(datas)
	gsplat.Sort(datas)
	gsplat.WriteSplat(output, datas)
	log.Println(input, " --> ", output)
	log.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spx2ply(args *cmn.OsArgs) {
	log.Println("Convert spx to ply.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadSpx(input)
	gsplat.RstTransformDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(int(header.ShDegree), args)
	gsplat.WritePly(output, datas, comment, shDegree)
	log.Println(input, " --> ", output)
	log.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}
func spx2splat(args *cmn.OsArgs) {
	log.Println("Convert spx to splat.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	_, datas := gsplat.ReadSpx(input)
	gsplat.RstTransformDatas(datas)
	gsplat.Sort(datas)
	gsplat.WriteSplat(output, datas)
	log.Println(input, " --> ", output)
	log.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}
func spx2spz(args *cmn.OsArgs) {
	log.Println("Convert spx to spz.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadSpx(input)
	gsplat.RstTransformDatas(datas)
	gsplat.Sort(datas)
	shDegree := getArgShDegree(int(header.ShDegree), args)
	gsplat.WriteSpz(output, datas, shDegree)
	log.Println(input, " --> ", output)
	log.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}
func spx2spx(args *cmn.OsArgs) {
	log.Println("Convert spx to spx.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadSpx(input)
	gsplat.RstTransformDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(int(header.ShDegree), args)
	flag1 := getArgFlag(0, args, "-f1", "--flag1")
	flag2 := getArgFlag(0, args, "-f2", "--flag2")
	flag3 := getArgFlag(0, args, "-f3", "--flag3")
	gsplat.WriteSpx(output, datas, comment, shDegree, flag1, flag2, flag3)
	log.Println(input, " --> ", output)
	log.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spz2ply(args *cmn.OsArgs) {
	log.Println("Convert spz to ply.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadSpz(input, false)
	gsplat.RstTransformDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(int(header.ShDegree), args)
	gsplat.WritePly(output, datas, comment, shDegree)
	log.Println(input, " --> ", output)
	log.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spz2splat(args *cmn.OsArgs) {
	log.Println("Convert spz to splat.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	_, datas := gsplat.ReadSpz(input, false)
	gsplat.RstTransformDatas(datas)
	gsplat.Sort(datas)
	gsplat.WriteSplat(output, datas)
	log.Println(input, " --> ", output)
	log.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spz2spx(args *cmn.OsArgs) {
	log.Println("Convert spz to spx.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadSpz(input, false)
	gsplat.RstTransformDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(int(header.ShDegree), args)
	flag1 := getArgFlag(0, args, "-f1", "--flag1")
	flag2 := getArgFlag(0, args, "-f2", "--flag2")
	flag3 := getArgFlag(0, args, "-f3", "--flag3")
	gsplat.WriteSpx(output, datas, comment, shDegree, flag1, flag2, flag3)
	log.Println(input, " --> ", output)
	log.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}
func spz2spz(args *cmn.OsArgs) {
	log.Println("Convert spz to spz.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadSpz(input, false)
	gsplat.RstTransformDatas(datas)
	gsplat.Sort(datas)
	shDegree := getArgShDegree(int(header.ShDegree), args)
	gsplat.WriteSpz(output, datas, shDegree)
	log.Println(input, " --> ", output)
	log.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func checkInputFileExists(args *cmn.OsArgs) string {
	input := args.GetArgIgnorecase("-i", "--input")
	cmn.ExitOnConditionError(input == "", errors.New(`please specify the input file`))
	cmn.ExitOnConditionError(!cmn.IsExistFile(input), errors.New("file not found: "+input))
	return input
}

func createOutputDir(args *cmn.OsArgs) string {
	output := args.GetArgIgnorecase("-o", "--output")
	cmn.ExitOnConditionError(output == "", errors.New(`please specify the output file`))
	cmn.ExitOnError(cmn.MkdirAll(cmn.Dir(output)))
	return output
}

func getArgShDegree(dataShDegree int, args *cmn.OsArgs) int {
	shDegree := dataShDegree
	if args.HasArg("-sh", "--shDegree") {
		sh := cmn.StringToInt(args.GetArgIgnorecase("-sh", "--shDegree"), -1)
		if sh >= 0 && sh <= 3 {
			shDegree = sh
		}
	}
	return shDegree
}

func getArgFlag(defaultFlag uint8, args *cmn.OsArgs, arg1 string, arg2 string) uint8 {
	flag := defaultFlag
	if args.HasArg(arg1, arg2) {
		val := cmn.StringToInt(args.GetArgIgnorecase(arg1, arg2), -1)
		if val >= 0 && val < 256 {
			flag = uint8(val)
		}
	}
	return flag
}
