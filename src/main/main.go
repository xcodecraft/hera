package main

import (
	_ "example"
	"fmt"
	"hera"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("hera starting...\n")
	n := hera.Classic()
	hera.Logger.Init("hera", hera.LevelDebug)
	n.Run(8083)
}
