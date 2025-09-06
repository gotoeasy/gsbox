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
	} else if args.HasCmd("k2p", "ksplat2ply") {
		ksplat2ply(args)
	} else if args.HasCmd("k2s", "ksplat2splat") {
		ksplat2splat(args)
	} else if args.HasCmd("k2x", "ksplat2spx") {
		ksplat2spx(args)
	} else if args.HasCmd("k2z", "ksplat2spx") {
		ksplat2spz(args)
	} else if args.HasCmd("ps", "printSplat") {
		printSplat(args)
	} else if args.HasCmd("join") {
		join(args)
	} else if args.HasCmd("info") {
		info(args)
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
	fmt.Println("  p2s, ply2splat                  convert ply to splat")
	fmt.Println("  p2x, ply2spx                    convert ply to spx")
	fmt.Println("  p2z, ply2spz                    convert ply to spz")
	fmt.Println("  p2p, ply2ply                    convert ply to ply")
	fmt.Println("  s2p, splat2ply                  convert splat to ply")
	fmt.Println("  s2x, splat2spx                  convert splat to spx")
	fmt.Println("  s2z, splat2spz                  convert splat to spz")
	fmt.Println("  s2s, splat2splat                convert splat to splat")
	fmt.Println("  x2p, spx2ply                    convert spx to ply")
	fmt.Println("  x2s, spx2splat                  convert spx to splat")
	fmt.Println("  x2z, spx2spz                    convert spx to spz")
	fmt.Println("  x2x, spx2spx                    convert spx to spx")
	fmt.Println("  z2p, spz2ply                    convert spz to ply")
	fmt.Println("  z2s, spz2splat                  convert spz to splat")
	fmt.Println("  z2x, spz2spx                    convert spz to spx")
	fmt.Println("  z2z, spz2spz                    convert spz to spz")
	fmt.Println("  k2p, ksplat2ply                 convert ksplat to ply")
	fmt.Println("  k2s, ksplat2splat               convert ksplat to splat")
	fmt.Println("  k2x, ksplat2spx                 convert ksplat to spx")
	fmt.Println("  k2z, ksplat2spx                 convert ksplat to spz")
	fmt.Println("  ps,  printsplat                 print data to text file like splat format layout")
	fmt.Println("  join                            join the input model files into a single output file")
	fmt.Println("  info <file>                     display the model file information")
	fmt.Println("  -i,  --input <file>             specify the input file")
	fmt.Println("  -o,  --output <file>            specify the output file")
	fmt.Println("  -ct, --compression-type <type>  specify the compression type(0:gzip,1:xz) for spx output, default is gzip")
	fmt.Println("  -c,  --comment <text>           specify the comment for ply/spx output")
	fmt.Println("  -a,  --alpha <num>              Specify the minimum alpha(0~255) to filter the output splat data")
	fmt.Println("  -bs, --block-size <num>         specify the block size(64~1024000) for spx output (default 102400)")
	fmt.Println("  -bf, --block-format <num>       specify the block data format for spx output (default 19)")
	fmt.Println("  -sh, --shDegree <num>           specify the SH degree(0~3) for output")
	fmt.Println("  -f2, --is-inverted <bool>       specify the header flag2(IsInverted) for spx output, default is false")
	fmt.Println("  -rx, --rotateX <num>            specify the rotation angle in degrees about the x-axis for transform")
	fmt.Println("  -ry, --rotateY <num>            specify the rotation angle in degrees about the y-axis for transform")
	fmt.Println("  -rz, --rotateZ <num>            specify the rotation angle in degrees about the z-axis for transform")
	fmt.Println("  -s,  --scale <num>              specify a uniform scaling factor(0.001~1000) for transform")
	fmt.Println("  -tx, --translateX <num>         specify the translation value about the x-axis for transform")
	fmt.Println("  -ty, --translateY <num>         specify the translation value about the y-axis for transform")
	fmt.Println("  -tz, --translateZ <num>         specify the translation value about the z-axis for transform")
	fmt.Println("  -to, --transform-order <RST>    specify the transform order (RST/RTS/SRT/STR/TRS/TSR), default is RST")
	fmt.Println("  -ov, --output-version <num>     specify the spz output version(2~3), default is 2")
	fmt.Println("  -v,  --version                  display version information")
	fmt.Println("  -h,  --help                     display help information")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  gsbox ply2splat -i /path/to/input.ply -o /path/to/output.splat")
	fmt.Println("  gsbox s2x -i /path/to/input.splat -o /path/to/output.spx -c \"your comment\" -bs 10240 -ct xz")
	fmt.Println("  gsbox x2z -i /path/to/input.spx -o /path/to/output.spz -sh 0 -rz 90 -s 0.9 -tx 0.1 -to TRS")
	fmt.Println("  gsbox z2p -i /path/to/input.spz -o /path/to/output.ply -c \"your comment\"")
	fmt.Println("  gsbox k2z -i /path/to/input.ksplat -o /path/to/output.spz -ov 3")
	fmt.Println("  gsbox join -i a.ply -i b.splat -i c.spx -i d.spz -o output.spx")
	fmt.Println("  gsbox info -i /path/to/file.spx")
	fmt.Println("")
}

func info(args *cmn.OsArgs) {
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
	isKsplat := cmn.Endwiths(input, ".ksplat", true)
	count := 0

	shDegree := 0
	if isPly {
		headerString, err := gsplat.ReadPlyHeaderString(input, 1024)
		cmn.ExitOnError(err)
		fmt.Print(headerString)
		header, err := gsplat.ReadPlyHeader(input)
		if err == nil && header.ChunkCount > 0 {
			count = header.VertexCount
			shDegree = header.MaxShDegree()
		}
	} else if isSpx {
		header := gsplat.ParseSpxHeader(input)
		fmt.Println(header.ToStringSpx())
		count = int(header.SplatCount)
		shDegree = int(header.ShDegree)
	} else if isSpz {
		header, _ := gsplat.ReadSpz(input, true)
		fmt.Println(header.ToString())
		count = int(header.NumPoints)
		shDegree = int(header.ShDegree)
	} else if isKsplat {
		secHeader, mainHeader, _ := gsplat.ReadKsplat(input, true)
		fmt.Println(mainHeader.ToString())
		fmt.Println("[Section 0]")
		fmt.Println(secHeader.ToString())
		count = mainHeader.SplatCount
		shDegree = int(mainHeader.ShDegree)
	} else if isSplat {
		fileInfo, err := os.Stat(input)
		cmn.ExitOnError(err)
		fileSize := fileInfo.Size()
		if (fileSize)%32 > 0 {
			cmn.ExitOnError(errors.New("invalid splat format"))
		} else {
			count = int(fileSize / 32)
			fmt.Println("SplatCount :", fileSize/32)
		}
	} else {
		cmn.ExitOnError(errors.New("the input file must be (ply | splat | spx | spz | ksplat) format"))
	}

	if count > 0 {
		fmt.Println("\n[Info]", gsplat.CompressionInfo(input, count, shDegree))
	}
}

func printSplat(args *cmn.OsArgs) {
	log.Println("[Info] print data to text file like splat format layout.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)

	var datas []*gsplat.SplatData
	shDegree := 0
	if cmn.Endwiths(input, ".ply", true) {
		header, ds := gsplat.ReadPly(input)
		datas = ds
		shDegree = max(header.MaxShDegree(), shDegree)
	} else if cmn.Endwiths(input, ".splat", true) {
		datas = gsplat.ReadSplat(input)
	} else if cmn.Endwiths(input, ".spx", true) {
		header, ds := gsplat.ReadSpx(input)
		datas = ds
		shDegree = max(int(header.ShDegree), shDegree)
	} else if cmn.Endwiths(input, ".spz", true) {
		header, ds := gsplat.ReadSpz(input, false)
		datas = ds
		shDegree = max(int(header.ShDegree), shDegree)
	} else if cmn.Endwiths(input, ".ksplat", true) {
		_, header, ds := gsplat.ReadKsplat(input, false)
		datas = ds
		shDegree = max(int(header.ShDegree), shDegree)
	} else {
		cmn.ExitOnError(errors.New("the input file must be (ply | splat | spx | spz | ksplat) format"))
	}
	sh := cmn.StringToInt(args.GetArgIgnorecase("-sh", "--shDegree"), -1)
	if sh >= 0 && sh <= 3 {
		shDegree = sh // 优先按参数要求级别，未指定时按数据实际级别输出
	}

	gsplat.PrintSplat(output, datas, shDegree)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func join(args *cmn.OsArgs) {
	log.Println("[Info] join the input model files into a single output file")
	startTime := time.Now()
	inputs := checkInputFilesExists(args)
	output := createOutputDir(args)
	isOutPly := cmn.Endwiths(output, ".ply", true)
	isOutSplat := cmn.Endwiths(output, ".splat", true)
	isOutSpx := cmn.Endwiths(output, ".spx", true)
	isOutSpz := cmn.Endwiths(output, ".spz", true)
	cmn.ExitOnConditionError(!isOutPly && !isOutSplat && !isOutSpx && !isOutSpz, errors.New("output file must be (ply | splat | spx | spz) format"))

	datas := make([]*gsplat.SplatData, 0)
	maxFromShDegree := 0
	for _, file := range inputs {
		if cmn.Endwiths(file, ".ply", true) {
			header, ds := gsplat.ReadPly(file)
			maxFromShDegree = max(header.MaxShDegree(), maxFromShDegree)
			datas = append(datas, ds...)
		} else if cmn.Endwiths(file, ".splat", true) {
			ds := gsplat.ReadSplat(file)
			datas = append(datas, ds...)
		} else if cmn.Endwiths(file, ".spx", true) {
			header, ds := gsplat.ReadSpx(file)
			maxFromShDegree = max(int(header.ShDegree), maxFromShDegree)
			datas = append(datas, ds...)
		} else if cmn.Endwiths(file, ".spz", true) {
			header, ds := gsplat.ReadSpz(file, false)
			maxFromShDegree = max(int(header.ShDegree), maxFromShDegree)
			datas = append(datas, ds...)
		} else if cmn.Endwiths(file, ".ksplat", true) {
			_, header, ds := gsplat.ReadKsplat(file, false)
			maxFromShDegree = max(header.ShDegree, maxFromShDegree)
			datas = append(datas, ds...)
		}
	}
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	shDegree := getArgShDegree(maxFromShDegree, args)
	if isOutPly {
		comment := args.GetArgIgnorecase("-c", "--comment")
		gsplat.WritePly(output, datas, comment, shDegree)
	} else if isOutSplat {
		gsplat.WriteSplat(output, datas)
	} else if isOutSpx {
		comment := args.GetArgIgnorecase("-c", "--comment")
		gsplat.WriteSpx(output, datas, comment, shDegree)
	} else if isOutSpz {
		gsplat.WriteSpz(output, datas, shDegree)
	}
	log.Println("[Info]", inputs, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), shDegree))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ply2splat(args *cmn.OsArgs) {
	log.Println("[Info] convert ply to splat.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	_, datas := gsplat.ReadPly(input)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	gsplat.WriteSplat(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), 0))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ply2spx(args *cmn.OsArgs) {
	log.Println("[Info] convert ply to spx.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadPly(input)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(header.MaxShDegree(), args)
	gsplat.WriteSpx(output, datas, comment, shDegree)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), shDegree))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ply2spz(args *cmn.OsArgs) {
	log.Println("[Info] convert ply to spz.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadPly(input)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	shDegree := getArgShDegree(header.MaxShDegree(), args)
	gsplat.WriteSpz(output, datas, shDegree)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), shDegree))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ply2ply(args *cmn.OsArgs) {
	log.Println("[Info] convert ply to ply.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadPly(input)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(header.MaxShDegree(), args)
	gsplat.WritePly(output, datas, comment, shDegree)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), shDegree))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func splat2ply(args *cmn.OsArgs) {
	log.Println("[Info] convert splat to ply.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	datas := gsplat.ReadSplat(input)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(0, args)
	gsplat.WritePly(output, datas, comment, shDegree)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), shDegree))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func splat2spx(args *cmn.OsArgs) {
	log.Println("[Info] convert splat to spx.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	datas := gsplat.ReadSplat(input)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(0, args)
	gsplat.WriteSpx(output, datas, comment, shDegree)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), shDegree))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func splat2spz(args *cmn.OsArgs) {
	log.Println("[Info] convert splat to spz.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	datas := gsplat.ReadSplat(input)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	gsplat.WriteSpz(output, datas, 0)
	shDegree := getArgShDegree(0, args)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), shDegree))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func splat2splat(args *cmn.OsArgs) {
	log.Println("[Info] convert splat to splat.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	datas := gsplat.ReadSplat(input)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	gsplat.WriteSplat(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), 0))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spx2ply(args *cmn.OsArgs) {
	log.Println("[Info] convert spx to ply.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadSpx(input)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(int(header.ShDegree), args)
	gsplat.WritePly(output, datas, comment, shDegree)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), shDegree))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spx2splat(args *cmn.OsArgs) {
	log.Println("[Info] convert spx to splat.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	_, datas := gsplat.ReadSpx(input)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	gsplat.WriteSplat(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), 0))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spx2spz(args *cmn.OsArgs) {
	log.Println("[Info] convert spx to spz.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadSpx(input)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	shDegree := getArgShDegree(int(header.ShDegree), args)
	gsplat.WriteSpz(output, datas, shDegree)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), shDegree))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spx2spx(args *cmn.OsArgs) {
	log.Println("[Info] convert spx to spx.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadSpx(input)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(int(header.ShDegree), args)
	gsplat.WriteSpx(output, datas, comment, shDegree)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), shDegree))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spz2ply(args *cmn.OsArgs) {
	log.Println("[Info] convert spz to ply.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadSpz(input, false)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(int(header.ShDegree), args)
	gsplat.WritePly(output, datas, comment, shDegree)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), shDegree))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spz2splat(args *cmn.OsArgs) {
	log.Println("[Info] convert spz to splat.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	_, datas := gsplat.ReadSpz(input, false)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	gsplat.WriteSplat(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), 0))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spz2spx(args *cmn.OsArgs) {
	log.Println("[Info] convert spz to spx.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadSpz(input, false)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(int(header.ShDegree), args)
	gsplat.WriteSpx(output, datas, comment, shDegree)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), shDegree))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spz2spz(args *cmn.OsArgs) {
	log.Println("[Info] convert spz to spz.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	header, datas := gsplat.ReadSpz(input, false)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	shDegree := getArgShDegree(int(header.ShDegree), args)
	gsplat.WriteSpz(output, datas, shDegree)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), shDegree))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ksplat2ply(args *cmn.OsArgs) {
	log.Println("[Info] convert ksplat to ply.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	_, header, datas := gsplat.ReadKsplat(input, false)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(header.ShDegree, args)
	gsplat.WritePly(output, datas, comment, shDegree)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), shDegree))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ksplat2splat(args *cmn.OsArgs) {
	log.Println("[Info] convert ksplat to splat.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	_, _, datas := gsplat.ReadKsplat(input, false)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	gsplat.WriteSplat(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), 0))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ksplat2spx(args *cmn.OsArgs) {
	log.Println("[Info] convert ksplat to spx.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	_, header, datas := gsplat.ReadKsplat(input, false)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	comment := args.GetArgIgnorecase("-c", "--comment")
	shDegree := getArgShDegree(header.ShDegree, args)
	gsplat.WriteSpx(output, datas, comment, shDegree)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), shDegree))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ksplat2spz(args *cmn.OsArgs) {
	log.Println("[Info] convert ksplat to spz.")
	startTime := time.Now()
	input := checkInputFileExists(args)
	output := createOutputDir(args)
	_, header, datas := gsplat.ReadKsplat(input, false)
	gsplat.TransformDatas(datas)
	datas = gsplat.FilterDatas(datas)
	gsplat.Sort(datas)
	shDegree := getArgShDegree(header.ShDegree, args)
	gsplat.WriteSpz(output, datas, shDegree)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), shDegree))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func checkInputFilesExists(args *cmn.OsArgs) []string {
	inputs := args.GetArgsIgnorecase("-i", "--input")
	for i := 0; i < len(inputs); i++ {
		cmn.ExitOnConditionError(inputs[i] == "", errors.New(`please specify the input file`))
		cmn.ExitOnConditionError(!cmn.IsExistFile(inputs[i]), errors.New("file not found: "+inputs[i]))
	}
	return inputs
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
