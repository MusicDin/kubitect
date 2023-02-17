package main

import (
	"os"

	"github.com/MusicDin/kubitect/cli/cmd"
	"github.com/MusicDin/kubitect/cli/ui"
)

func main() {
	err := cmd.NewRootCmd().Execute()

	if err != nil {
		ui.PrintBlockE(err)
		os.Exit(1)
	}
}
