package main

import (
	"fmt"
	"os"

	"github.com/parashmaity/go-builder/build"
)

func main() {
	if err := build.BuildCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
