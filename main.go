package main

import (
	"fmt"
	"gsbox/cmn"
	"gsbox/gsplat"
	"time"
)

const VER = "v2.0.0"

func main() {

	args := cmn.ParseArgs()
	if args.HasCmdIgnorecase("ply2splat") {
		ply2splat(args)
	} else if args.HasCmdIgnorecase("splat2ply") {
		splat2ply(args)
	} else if args.HasCmdIgnorecase("splat2splat") {
		splat2splat(args)
	} else if args.HasCmdIgnorecase("ply2ply") {
		ply2ply(args)
	} else if args.HasCmdIgnorecase("info") {
		// info
		header, err := gsplat.ReadPlyHeaderString(args.LastParam, 1024)
		cmn.ExitOnError(err)
		fmt.Print(header)
	} else {
		input := args.GetArgIgnorecase("-i", "--input")
		output := args.GetArgIgnorecase("-o", "--output")
		if (!cmn.Endwiths(input, ".ply", true) && !cmn.Endwiths(input, ".splat", true)) ||
			(!cmn.Endwiths(output, ".ply", true) && !cmn.Endwiths(output, ".splat", true)) {
			// usage
			fmt.Println("gsbox", VER)
			fmt.Println("")
			fmt.Println("Usage")
			fmt.Println("gsbox ply2splat -i /path/to/input.ply -o /path/to/output.splat [-h header]")
			fmt.Println("gsbox splat2ply -i /path/to/input.splat -o /path/to/output.ply [-h header]")
			fmt.Println("gsbox -i /path/to/input.ply -o /path/to/output.splat [-h header]")
			fmt.Println("gsbox info /path/to/file.ply")
			fmt.Println("")
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

func ply2splat(args *cmn.OsArgs) {
	startTime := time.Now()
	datas := gsplat.ReadPly(args.GetArgIgnorecase("-i", "--input"))
	gsplat.Sort(datas)
	gsplat.WriteSplat(args.GetArgIgnorecase("-o", "--output"), datas, args.GetArgIgnorecase("-h", "--header", "-oh", "--output-header"))
	fmt.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}

func splat2ply(args *cmn.OsArgs) {
	startTime := time.Now()
	datas := gsplat.ReadSplat(args.GetArgIgnorecase("-i", "--input"), args.GetArgIgnorecase("-h", "--header", "-ih", "--input-header"))
	gsplat.WritePly(args.GetArgIgnorecase("-o", "--output"), datas)
	fmt.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}
func splat2splat(args *cmn.OsArgs) {
	startTime := time.Now()
	ih := args.GetArgIgnorecase("-ih", "--input-header")
	oh := args.GetArgIgnorecase("-oh", "--output-header")
	h := args.GetArgIgnorecase("-h", "--header")
	if ih == "" && oh == "" {
		ih = h
		oh = h
	}
	datas := gsplat.ReadSplat(args.GetArgIgnorecase("-i", "--input"), ih)
	gsplat.WriteSplat(args.GetArgIgnorecase("-o", "--output"), datas, oh)
	fmt.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}
func ply2ply(args *cmn.OsArgs) {
	startTime := time.Now()
	datas := gsplat.ReadPly(args.GetArgIgnorecase("-i", "--input"))
	gsplat.Sort(datas)
	gsplat.WritePly(args.GetArgIgnorecase("-o", "--output"), datas)
	fmt.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
}
