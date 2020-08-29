///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package main

import (
	"digger/shell"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			args = []string{"--version"}
		}
	}
	newArgs := append([]string{os.Args[0]}, args...)
	shell.Parse(newArgs)
	//app := gin.Default()
}
