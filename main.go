package main

import (
	"fmt"
	"gsbox/cmn"
	"gsbox/gsplat"
	"time"
)

const VER = "v1.1.0"

func main() {

	args := cmn.ParseArgs()
	startTime := time.Now()
	if args.Command == "ply2splat" {
		// ply2splat
		datas := gsplat.ReadPly(args.Input)
		gsplat.Sort(datas)
		gsplat.WriteSplat(args.Output, datas)
		fmt.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
	} else if args.Command == "splat2ply" {
		// splat2ply
		datas := gsplat.ReadSplat(args.Input)
		gsplat.WritePly(args.Output, datas)
		fmt.Println("Processing time conversion:", cmn.GetTimeInfo(time.Since(startTime).Milliseconds()))
	} else if args.Command == "info" {
		// info
		header, err := gsplat.ReadPlyHeader(args.Input)
		cmn.ExitOnError(err)
		fmt.Print(header.ToString())
	} else {
		// usage
		fmt.Println("gsbox", VER)
		fmt.Println("")
		fmt.Println("Usage")
		fmt.Println("gsbox ply2splat from.ply to.splat")
		fmt.Println("gsbox splat2ply from.splat to.ply")
		fmt.Println("gsbox info file.ply")
		fmt.Println("")
	}

}
