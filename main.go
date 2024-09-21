package main

import (
	"errors"
	"fmt"
	"gsbox/cmn"
	"gsbox/gsplat"
	"time"
)

const VER = "v2.1.1"

func main() {

	args := cmn.ParseArgs("-v", "-version", "--version", "-h", "-help", "--help")
	if args.HasCmd("-v", "-version", "--version") && args.ArgCount == 2 {
		version()
	} else if args.HasCmd("-h", "-help", "--help") && args.ArgCount == 2 {
		usage()
	} else if args.HasCmd("ply2splat") {
		ply2splat(args)
	} else if args.HasCmd("splat2ply") {
		splat2ply(args)
	} else if args.HasCmd("splat2splat") {
		splat2splat(args)
	} else if args.HasCmd("ply2ply") {
		ply2ply(args)
	} else if args.HasCmd("info") {
		plyInfo(args)
	} else {
		input := args.GetArgIgnorecase("-i", "--input")
		output := args.GetArgIgnorecase("-o", "--output")
		if (!cmn.Endwiths(input, ".ply", true) && !cmn.Endwiths(input, ".splat", true)) ||
			(!cmn.Endwiths(output, ".ply", true) && !cmn.Endwiths(output, ".splat", true)) {
			usage()
		} else {
			// 按后缀名识别格式判断装换
			if cmn.Endwiths(input, ".ply", true) {
				if cmn.Endwiths(output, ".splat", true) {
					ply2splat(args)
				} else {
					ply2ply(args)
				}
			} else {
				if cmn.Endwiths(output, ".splat", true) {
					splat2splat(args)
				} else {
					splat2ply(args)
				}
			}
		}
	}
}

func version() {
	fmt.Println("")
	fmt.Println("gsbox", VER)
	fmt.Println("homepage", "https://github.com/gotoeasy/gsbox")
}
func usage() {
	fmt.Println("")
	fmt.Println("gsbox", VER)
	fmt.Println("homepage", "https://github.com/gotoeasy/gsbox")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  gsbox [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  ply2splat                convert ply to splat")
	fmt.Println("  splat2ply                convert splat to ply")
	fmt.Println("  simple-ply               simple mode to write ply")
	fmt.Println("  info <plyfile>           display the ply header")
	fmt.Println("  -i, --input <file>       specify the input file")
	fmt.Println("  -o, --output <file>      specify the output file")
	fmt.Println("  -c, --comment <text>     output ply with comment")
	fmt.Println("  -h, --help               display help information")
	fmt.Println("  -v, --version            display version information")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  gsbox ply2splat -i /path/to/input.ply -o /path/to/output.splat")
	fmt.Println("  gsbox splat2ply -i /path/to/input.splat -o /path/to/output.ply")
	fmt.Println("  gsbox -i /path/to/input.ply -o /path/to/output.splat")
	fmt.Println("  gsbox -i /path/to/input.splat -o /path/to/output.ply")
	fmt.Println("  gsbox -i /path/to/input.splat -o /path/to/output.ply simple-ply")
	fmt.Println("  gsbox -i /path/to/input.splat -o /path/to/output.ply -c \"your comment\"")
	fmt.Println("  gsbox info -i /path/to/file.ply")
	fmt.Println("  gsbox info /path/to/file.ply")
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

	if !cmn.Endwiths(input, ".ply", true) {
		cmn.ExitOnError(errors.New("the input file must be ply format"))
	}
	if !cmn.IsExistFile(input) {
		cmn.ExitOnError(errors.New("file not found: " + input))
	}
	header, err := gsplat.ReadPlyHeaderString(input, 1024)
	cmn.ExitOnError(err)
	fmt.Print(header)
}
func ply2splat(args *cmn.OsArgs) {
	startTime := time.Now()
	checkInputFileExists(args)
	createOutputDir(args)
	datas := gsplat.ReadPly(args.GetArgIgnorecase("-i", "--input"), "ply-3dgs")
	gsplat.Sort(datas)
	gsplat.WriteSplat(args.GetArgIgnorecase("-o", "--output"), datas)
	fmt.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func splat2ply(args *cmn.OsArgs) {
	startTime := time.Now()
	checkInputFileExists(args)
	createOutputDir(args)
	datas := gsplat.ReadSplat(args.GetArgIgnorecase("-i", "--input"))
	gsplat.WritePly(args.GetArgIgnorecase("-o", "--output"), datas, args.GetArgIgnorecase("-c", "--comment"), args.HasCmd("simple-ply"))
	fmt.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}
func splat2splat(args *cmn.OsArgs) {
	startTime := time.Now()
	checkInputFileExists(args)
	createOutputDir(args)

	datas := gsplat.ReadSplat(args.GetArgIgnorecase("-i", "--input"))
	gsplat.Sort(datas)
	gsplat.WriteSplat(args.GetArgIgnorecase("-o", "--output"), datas)
	fmt.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}
func ply2ply(args *cmn.OsArgs) {
	startTime := time.Now()
	checkInputFileExists(args)
	createOutputDir(args)
	datas := gsplat.ReadPly(args.GetArgIgnorecase("-i", "--input"), "ply-3dgs")
	gsplat.WritePly(args.GetArgIgnorecase("-o", "--output"), datas, args.GetArgIgnorecase("-c", "--comment"), args.HasCmd("simple-ply"))
	fmt.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func checkInputFileExists(args *cmn.OsArgs) {
	input := args.GetArgIgnorecase("-i", "--input")
	if input == "" {
		cmn.ExitOnError(errors.New(`please specify the input file`))
	}
	if !cmn.IsExistFile(input) {
		cmn.ExitOnError(errors.New("file not found: " + input))
	}
}

func createOutputDir(args *cmn.OsArgs) {
	output := args.GetArgIgnorecase("-o", "--output")
	if output == "" {
		cmn.ExitOnError(errors.New(`please specify the output file`))
	}
	cmn.ExitOnError(cmn.MkdirAll(cmn.Dir(output)))
}
