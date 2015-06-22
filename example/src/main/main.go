package main

import (
	_ "api1"
	_ "api2"
	"os"
	"fmt"
	hera  "github.com/xcodecraft/hera"
)

func main() {
	curentDir, _ := os.Getwd()
	hera.Run(fmt.Sprintf("%s/../conf/mate.yaml",curentDir)
}
