package main

import (
	_ "api"
	_ "api1"
	_ "api2"
	"fmt"
	"os"

	hera "github.com/xcodecraft/hera"
)

func main() {
	curentDir, _ := os.Getwd()
	hera.Run(fmt.Sprintf("%s/conf/app.yaml", curentDir))
}
