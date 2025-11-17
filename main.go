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
	args := gsplat.InitArgs()
	startTime := time.Now()
	if args.HasCmd("-v", "-version", "--version") && args.ArgCount == 2 {
		version()
	} else if args.HasCmd("-h", "-help", "--help") && args.ArgCount == 2 {
		usage()
	} else if args.HasCmd("p2s", "ply2splat") {
		ply2splat()
	} else if args.HasCmd("p2x", "ply2spx") {
		ply2spx()
	} else if args.HasCmd("p2z", "ply2spz") {
		ply2spz()
	} else if args.HasCmd("p2g", "ply2sog") {
		ply2sog()
	} else if args.HasCmd("p2p", "ply2ply") {
		ply2ply()
	} else if args.HasCmd("s2p", "splat2ply") {
		splat2ply()
	} else if args.HasCmd("s2x", "splat2spx") {
		splat2spx()
	} else if args.HasCmd("s2z", "splat2spz") {
		splat2spz()
	} else if args.HasCmd("s2g", "splat2sog") {
		splat2sog()
	} else if args.HasCmd("s2s", "splat2splat") {
		splat2splat()
	} else if args.HasCmd("x2p", "spx2ply") {
		spx2ply()
	} else if args.HasCmd("x2s", "spx2splat") {
		spx2splat()
	} else if args.HasCmd("x2z", "spx2spz") {
		spx2spz()
	} else if args.HasCmd("x2g", "spx2sog") {
		spx2sog()
	} else if args.HasCmd("x2x", "spx2spx") {
		spx2spx()
	} else if args.HasCmd("z2p", "spz2ply") {
		spz2ply()
	} else if args.HasCmd("z2s", "spz2splat") {
		spz2splat()
	} else if args.HasCmd("z2x", "spz2spx") {
		spz2spx()
	} else if args.HasCmd("z2g", "spz2sog") {
		spz2sog()
	} else if args.HasCmd("z2z", "spz2spz") {
		spz2spz()
	} else if args.HasCmd("k2p", "ksplat2ply") {
		ksplat2ply()
	} else if args.HasCmd("k2s", "ksplat2splat") {
		ksplat2splat()
	} else if args.HasCmd("k2x", "ksplat2spx") {
		ksplat2spx()
	} else if args.HasCmd("k2z", "ksplat2spz") {
		ksplat2spz()
	} else if args.HasCmd("k2g", "ksplat2sog") {
		ksplat2sog()
	} else if args.HasCmd("g2p", "sog2ply") {
		sog2ply()
	} else if args.HasCmd("g2s", "sog2splat") {
		sog2splat()
	} else if args.HasCmd("g2x", "sog2spx") {
		sog2spx()
	} else if args.HasCmd("g2z", "sog2spz") {
		sog2spz()
	} else if args.HasCmd("g2g", "sog2sog") {
		sog2sog()
	} else if args.HasCmd("ps", "printSplat") {
		printSplat()
	} else if args.HasCmd("join") {
		join()
	} else if args.HasCmd("info") {
		info(args)
	} else {
		usage()
	}
	dur := time.Since(startTime).Milliseconds()
	if dur < 1000 {
		time.Sleep((1000 - time.Duration(dur)) * time.Millisecond) // wait 1 second to get latest version
	}
	fmt.Print(cmn.NewVersionMessage)
	os.Exit(0)
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
	fmt.Println("  p2s, ply2splat                     convert ply to splat")
	fmt.Println("  p2x, ply2spx                       convert ply to spx")
	fmt.Println("  p2z, ply2spz                       convert ply to spz")
	fmt.Println("  p2g, ply2sog                       convert ply to sog")
	fmt.Println("  p2p, ply2ply                       convert ply to ply")
	fmt.Println("  s2p, splat2ply                     convert splat to ply")
	fmt.Println("  s2x, splat2spx                     convert splat to spx")
	fmt.Println("  s2z, splat2spz                     convert splat to spz")
	fmt.Println("  s2g, splat2sog                     convert splat to sog")
	fmt.Println("  s2s, splat2splat                   convert splat to splat")
	fmt.Println("  x2p, spx2ply                       convert spx to ply")
	fmt.Println("  x2s, spx2splat                     convert spx to splat")
	fmt.Println("  x2z, spx2spz                       convert spx to spz")
	fmt.Println("  x2g, spx2sog                       convert spx to sog")
	fmt.Println("  x2x, spx2spx                       convert spx to spx")
	fmt.Println("  z2p, spz2ply                       convert spz to ply")
	fmt.Println("  z2s, spz2splat                     convert spz to splat")
	fmt.Println("  z2x, spz2spx                       convert spz to spx")
	fmt.Println("  z2g, spz2sog                       convert spz to sog")
	fmt.Println("  z2z, spz2spz                       convert spz to spz")
	fmt.Println("  k2p, ksplat2ply                    convert ksplat to ply")
	fmt.Println("  k2s, ksplat2splat                  convert ksplat to splat")
	fmt.Println("  k2x, ksplat2spx                    convert ksplat to spx")
	fmt.Println("  k2z, ksplat2spx                    convert ksplat to spz")
	fmt.Println("  k2g, ksplat2sog                    convert ksplat to sog")
	fmt.Println("  g2p, sog2ply                       convert sog to ply")
	fmt.Println("  g2s, sog2splat                     convert sog to splat")
	fmt.Println("  g2x, sog2spx                       convert sog to spx")
	fmt.Println("  g2z, sog2spz                       convert sog to spz")
	fmt.Println("  g2g, sog2sog                       convert sog to sog")
	fmt.Println("  ps,  printsplat                    print data to text file like splat format layout")
	fmt.Println("  join                               join the input model files into a single output file")
	fmt.Println("  info <file>                        display the model file information")
	fmt.Println("  -i,  --input <file>                specify the input file")
	fmt.Println("  -o,  --output <file>               specify the output file")
	fmt.Println("  -q,  --quality <num>               specify the quality(1~9) for spx|spz|sog output, default is 5")
	fmt.Println("  -ct, --compression-type <type>     specify the compression type(0:gzip,1:xz) for spx output")
	fmt.Println("  -c,  --comment <text>              specify the comment for ply/spx output")
	fmt.Println("  -a,  --alpha <num>                 specify the minimum alpha(0~255) to filter the output splat data")
	fmt.Println("  -bs, --block-size <num>            specify the block size(4096~524288) for spx output (default is 65536)")
	fmt.Println("  -bf, --block-format <num>          specify the block data format(19~20) for spx output (default is 19)")
	fmt.Println("  -sh, --shDegree <num>              specify the SH degree(0~3) for output")
	fmt.Println("  -f1, --is-inverted <bool>          specify the header flag1(IsInverted) for spx output, default is false")
	fmt.Println("  -rx, --rotateX <num>               specify the rotation angle in degrees about the x-axis for transform")
	fmt.Println("  -ry, --rotateY <num>               specify the rotation angle in degrees about the y-axis for transform")
	fmt.Println("  -rz, --rotateZ <num>               specify the rotation angle in degrees about the z-axis for transform")
	fmt.Println("  -s,  --scale <num>                 specify a uniform scaling factor(0.001~1000) for transform")
	fmt.Println("  -tx, --translateX <num>            specify the translation value about the x-axis for transform")
	fmt.Println("  -ty, --translateY <num>            specify the translation value about the y-axis for transform")
	fmt.Println("  -tz, --translateZ <num>            specify the translation value about the z-axis for transform")
	fmt.Println("  -to, --transform-order <RST>       specify the transform order (RST/RTS/SRT/STR/TRS/TSR), default is RST")
	fmt.Println("  -ov, --output-version <num>        specify the output versions for spx|spz|sog, default is newest")
	fmt.Println("  -ki, --kmeans-iterations <num>     specify the kmeans iterations, default is 10")
	fmt.Println("  -kn, --kmeans-nearest-nodes <num>  specify the kmeans nearest nodes, default is 20")
	fmt.Println("  -v,  --version                     display version information")
	fmt.Println("  -h,  --help                        display help information")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  gsbox ply2sog -i /path/to/input.ply -o /path/to/output.sog")
	fmt.Println("  gsbox s2x -i /path/to/input.splat -o /path/to/output.spx -c \"your comment\" -bs 10240 -ct xz")
	fmt.Println("  gsbox x2z -i /path/to/input.spx -o /path/to/output.spz -sh 0 -rz 90 -s 0.9 -tx 0.1 -to TRS")
	fmt.Println("  gsbox z2p -i /path/to/input.spz -o /path/to/output.ply -c \"your comment\"")
	fmt.Println("  gsbox k2z -i /path/to/input.ksplat -o /path/to/output.spz -ov 3")
	fmt.Println("  gsbox g2x -i /path/to/input.sog -o /path/to/output.spx")
	fmt.Println("  gsbox g2x -i /path/to/meta.json -o /path/to/output.spx")
	fmt.Println("  gsbox join -i a.ply -i b.splat -i c.spx -i d.spz -i e.ksplat -i f.sog -i meta.json -o output.spx")
	fmt.Println("  gsbox ps -i /path/to/input.spx -o /path/to/output.txt")
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

	cmn.ExitOnConditionError(input == "", errors.New("please specify the input file"))
	cmn.ExitOnConditionError(!cmn.IsExistFile(input), errors.New("file not found: "+input))

	isPly := cmn.Endwiths(input, ".ply", true)
	isSpx := cmn.Endwiths(input, ".spx", true)
	isSplat := cmn.Endwiths(input, ".splat", true)
	isSpz := cmn.Endwiths(input, ".spz", true)
	isKsplat := cmn.Endwiths(input, ".ksplat", true)
	isSog := cmn.Endwiths(input, ".sog", true)
	isMetaJson := cmn.FileName(input) == "meta.json"
	count := 0

	shDegree := uint8(0)
	var sogFileSize int64
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
		shDegree = header.ShDegree
	} else if isSpz {
		header, _ := gsplat.ReadSpz(input)
		fmt.Println(header.ToString())
		count = int(header.NumPoints)
		shDegree = header.ShDegree
	} else if isKsplat {
		secHeader, mainHeader, _ := gsplat.ReadKsplat(input)
		fmt.Println(mainHeader.ToString())
		fmt.Println("[Section 0]")
		fmt.Println(secHeader.ToString())
		count = mainHeader.SplatCount
		shDegree = mainHeader.ShDegree
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
	} else if isSog || isMetaJson {
		version, cnt, degree, fileSize := gsplat.ReadSogInfo(input)
		count = cnt
		shDegree = degree
		sogFileSize = fileSize
		fmt.Println("Sog Version :", version)
		fmt.Println("Splat Count :", count)
		fmt.Println("SH Degree   :", shDegree)
	} else {
		cmn.ExitOnError(errors.New("the input file must be (ply | splat | spx | spz | ksplat | sog) format"))
	}
	gsplat.SetShDegreeFrom(shDegree)

	if count > 0 {
		fmt.Println("\n[Info]", gsplat.CompressionInfo(input, count, sogFileSize))
	}
}

func printSplat() {
	log.Println("[Info] print data to text file like splat format layout.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()

	var datas []*gsplat.SplatData
	shDegree := uint8(0)
	if cmn.Endwiths(input, ".ply", true) {
		header, ds := gsplat.ReadPly(input)
		datas = ds
		shDegree = max(header.MaxShDegree(), shDegree)
	} else if cmn.Endwiths(input, ".splat", true) {
		datas = gsplat.ReadSplat(input)
	} else if cmn.Endwiths(input, ".spx", true) {
		header, ds := gsplat.ReadSpx(input)
		datas = ds
		shDegree = max(header.ShDegree, shDegree)
	} else if cmn.Endwiths(input, ".spz", true) {
		header, ds := gsplat.ReadSpz(input)
		datas = ds
		shDegree = max(header.ShDegree, shDegree)
	} else if cmn.Endwiths(input, ".ksplat", true) {
		_, header, ds := gsplat.ReadKsplat(input)
		datas = ds
		shDegree = max(header.ShDegree, shDegree)
	} else if cmn.Endwiths(input, ".sog", true) || cmn.FileName(input) == "meta.json" {
		ds, h := gsplat.ReadSog(input)
		datas = ds
		shDegree = max(h.ShDegree, shDegree)
	} else {
		cmn.ExitOnError(errors.New("the input file must be (ply | splat | spx | spz | ksplat | sog) format"))
	}
	gsplat.SetShDegreeFrom(shDegree)

	gsplat.PrintSplat(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func join() {
	log.Println("[Info] join the input models into one")
	startTime := time.Now()
	inputs := gsplat.GetAndCheckInputFiles()
	output := gsplat.CreateOutputDir()
	isOutPly := cmn.Endwiths(output, ".ply", true)
	isOutSplat := cmn.Endwiths(output, ".splat", true)
	isOutSpx := cmn.Endwiths(output, ".spx", true)
	isOutSpz := cmn.Endwiths(output, ".spz", true)
	isOutSog := cmn.Endwiths(output, ".sog", true) || cmn.FileName(output) == "meta.json"

	ok := isOutPly || isOutSplat || isOutSpx || isOutSpz || isOutSog
	cmn.ExitOnConditionError(!ok, errors.New("output file must be (ply | splat | spx | spz | sog) format"))

	datas := make([]*gsplat.SplatData, 0)
	var maxFromShDegree uint8
	for _, file := range inputs {
		if cmn.Endwiths(file, ".ply", true) {
			header, ds := gsplat.ReadPly(file)
			maxFromShDegree = max(uint8(header.MaxShDegree()), maxFromShDegree)
			datas = append(datas, ds...)
		} else if cmn.Endwiths(file, ".splat", true) {
			ds := gsplat.ReadSplat(file)
			datas = append(datas, ds...)
		} else if cmn.Endwiths(file, ".spx", true) {
			header, ds := gsplat.ReadSpx(file)
			maxFromShDegree = max((header.ShDegree), maxFromShDegree)
			datas = append(datas, ds...)
		} else if cmn.Endwiths(file, ".spz", true) {
			header, ds := gsplat.ReadSpz(file)
			maxFromShDegree = max((header.ShDegree), maxFromShDegree)
			datas = append(datas, ds...)
		} else if cmn.Endwiths(file, ".ksplat", true) {
			_, header, ds := gsplat.ReadKsplat(file)
			maxFromShDegree = max(uint8(header.ShDegree), maxFromShDegree)
			datas = append(datas, ds...)
		} else if cmn.Endwiths(file, ".sog", true) || cmn.FileName(file) == "meta.json" || (cmn.Startwiths(file, "http") && cmn.Endwiths(file, "/meta.json")) {
			ds, h := gsplat.ReadSog(file)
			maxFromShDegree = max(h.ShDegree, maxFromShDegree)
			datas = append(datas, ds...)
		}
	}
	gsplat.SetShDegreeFrom(maxFromShDegree)
	datas = gsplat.ProcessDatas(datas)
	var fileSize int64
	if isOutPly {
		gsplat.WritePly(output, datas)
	} else if isOutSplat {
		gsplat.WriteSplat(output, datas)
	} else if isOutSpx {
		gsplat.WriteSpx(output, datas)
	} else if isOutSpz {
		gsplat.WriteSpz(output, datas)
	} else if isOutSog {
		fileSize = gsplat.WriteSog(output, datas)
	}
	log.Println("[Info]", inputs, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), fileSize))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ply2splat() {
	log.Println("[Info] convert ply to splat.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	header, datas := gsplat.ReadPly(input)
	gsplat.SetShDegreeFrom(header.MaxShDegree())
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSplat(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), 0))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ply2spx() {
	log.Println("[Info] convert ply to spx.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	header, datas := gsplat.ReadPly(input)
	gsplat.SetShDegreeFrom(header.MaxShDegree())
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSpx(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ply2spz() {
	log.Println("[Info] convert ply to spz.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	header, datas := gsplat.ReadPly(input)
	gsplat.SetShDegreeFrom(header.MaxShDegree())
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSpz(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ply2sog() {
	log.Println("[Info] convert ply to sog.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	header, datas := gsplat.ReadPly(input)
	gsplat.SetShDegreeFrom(header.MaxShDegree())
	datas = gsplat.ProcessDatas(datas)
	fileSize := gsplat.WriteSog(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), fileSize))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ply2ply() {
	log.Println("[Info] convert ply to ply.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	header, datas := gsplat.ReadPly(input)
	gsplat.SetShDegreeFrom(header.MaxShDegree())
	datas = gsplat.ProcessDatas(datas)
	gsplat.WritePly(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func splat2ply() {
	log.Println("[Info] convert splat to ply.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	datas := gsplat.ReadSplat(input)
	gsplat.SetShDegreeFrom(0)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WritePly(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func splat2spx() {
	log.Println("[Info] convert splat to spx.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	datas := gsplat.ReadSplat(input)
	gsplat.SetShDegreeFrom(0)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSpx(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func splat2spz() {
	log.Println("[Info] convert splat to spz.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	datas := gsplat.ReadSplat(input)
	gsplat.SetShDegreeFrom(0)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSpz(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func splat2sog() {
	log.Println("[Info] convert splat to sog.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	datas := gsplat.ReadSplat(input)
	gsplat.SetShDegreeFrom(0)
	datas = gsplat.ProcessDatas(datas)
	fileSize := gsplat.WriteSog(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), fileSize))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func splat2splat() {
	log.Println("[Info] convert splat to splat.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	datas := gsplat.ReadSplat(input)
	gsplat.SetShDegreeFrom(0)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSplat(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), 0))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spx2ply() {
	log.Println("[Info] convert spx to ply.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	header, datas := gsplat.ReadSpx(input)
	gsplat.SetShDegreeFrom(header.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WritePly(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spx2splat() {
	log.Println("[Info] convert spx to splat.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	header, datas := gsplat.ReadSpx(input)
	gsplat.SetShDegreeFrom(header.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSplat(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), 0))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spx2spz() {
	log.Println("[Info] convert spx to spz.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	header, datas := gsplat.ReadSpx(input)
	gsplat.SetShDegreeFrom(header.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSpz(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spx2sog() {
	log.Println("[Info] convert spx to sog.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	header, datas := gsplat.ReadSpx(input)
	gsplat.SetShDegreeFrom(header.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	fileSize := gsplat.WriteSog(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), fileSize))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spx2spx() {
	log.Println("[Info] convert spx to spx.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	header, datas := gsplat.ReadSpx(input)
	gsplat.SetShDegreeFrom(header.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSpx(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spz2ply() {
	log.Println("[Info] convert spz to ply.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	header, datas := gsplat.ReadSpz(input)
	gsplat.SetShDegreeFrom(header.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WritePly(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spz2splat() {
	log.Println("[Info] convert spz to splat.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	header, datas := gsplat.ReadSpz(input)
	gsplat.SetShDegreeFrom(header.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSplat(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), 0))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spz2spx() {
	log.Println("[Info] convert spz to spx.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	header, datas := gsplat.ReadSpz(input)
	gsplat.SetShDegreeFrom(header.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSpx(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spz2sog() {
	log.Println("[Info] convert spz to sog.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	header, datas := gsplat.ReadSpz(input)
	gsplat.SetShDegreeFrom(header.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	fileSize := gsplat.WriteSog(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), fileSize))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func spz2spz() {
	log.Println("[Info] convert spz to spz.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	header, datas := gsplat.ReadSpz(input)
	gsplat.SetShDegreeFrom(header.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSpz(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ksplat2ply() {
	log.Println("[Info] convert ksplat to ply.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	_, header, datas := gsplat.ReadKsplat(input)
	gsplat.SetShDegreeFrom(header.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WritePly(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ksplat2splat() {
	log.Println("[Info] convert ksplat to splat.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	_, header, datas := gsplat.ReadKsplat(input)
	gsplat.SetShDegreeFrom(header.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSplat(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), 0))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ksplat2spx() {
	log.Println("[Info] convert ksplat to spx.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	_, header, datas := gsplat.ReadKsplat(input)
	gsplat.SetShDegreeFrom(header.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSpx(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ksplat2spz() {
	log.Println("[Info] convert ksplat to spz.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	_, header, datas := gsplat.ReadKsplat(input)
	gsplat.SetShDegreeFrom(header.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSpz(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func ksplat2sog() {
	log.Println("[Info] convert ksplat to sog.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	_, header, datas := gsplat.ReadKsplat(input)
	gsplat.SetShDegreeFrom(header.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	fileSize := gsplat.WriteSog(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), fileSize))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func sog2ply() {
	log.Println("[Info] convert sog to ply.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	datas, h := gsplat.ReadSog(input)
	gsplat.SetShDegreeFrom(h.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WritePly(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func sog2splat() {
	log.Println("[Info] convert sog to splat.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	datas, h := gsplat.ReadSog(input)
	gsplat.SetShDegreeFrom(h.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSplat(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), 0))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func sog2spx() {
	log.Println("[Info] convert sog to spx.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	datas, h := gsplat.ReadSog(input)
	gsplat.SetShDegreeFrom(h.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSpx(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func sog2spz() {
	log.Println("[Info] convert sog to spz.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	datas, h := gsplat.ReadSog(input)
	gsplat.SetShDegreeFrom(h.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	gsplat.WriteSpz(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas)))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func sog2sog() {
	log.Println("[Info] convert sog to sog.")
	startTime := time.Now()
	input := gsplat.GetAndCheckInputFile()
	output := gsplat.CreateOutputDir()
	datas, h := gsplat.ReadSog(input)
	gsplat.SetShDegreeFrom(h.ShDegree)
	datas = gsplat.ProcessDatas(datas)
	fileSize := gsplat.WriteSog(output, datas)
	log.Println("[Info]", input, " --> ", output)
	log.Println("[Info]", gsplat.CompressionInfo(output, len(datas), fileSize))
	log.Println("[Info] processing time:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}
