package main

import (
	"flag"
	"fmt"
)

var version = "0.1.0"

func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.Parse()
	if showVersion {
		fmt.Println("version:", version)
		return
	}
}
