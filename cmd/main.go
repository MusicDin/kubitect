package main

import (
	"os"

	"github.com/MusicDin/kubitect/pkg/ui"
)

func main() {
	err := NewRootCmd().Execute()
	if err != nil {
		ui.PrintBlockE(err)
		os.Exit(1)
	}
}
