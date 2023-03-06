package main

import (
	"github.com/MusicDin/kubitect/cli/pkg/kubitect/cmd"
	"github.com/MusicDin/kubitect/cli/pkg/ui"
	"os"
)

func main() {
	err := cmd.NewRootCmd().Execute()

	if err != nil {
		ui.PrintBlockE(err)
		os.Exit(1)
	}
}
