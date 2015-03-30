package main

import (
	_ "example"
	"fmt"
	"hera"
)

func main() {
	fmt.Printf("hera starting...\n")
	r := hera.NewRouter()
	r.Start(":8083")
}
