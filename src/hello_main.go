package main

import (
	_ "example"
	"fmt"
	"hera"
)

func main() {
	fmt.Printf("hera starting...\n")
	r := hera.NewRouter()

	hera.Logger.Init("hera", hera.LevelDebug)
	hera.Logger.Info("hera start 8083...")
	r.Start(":8083")
}
