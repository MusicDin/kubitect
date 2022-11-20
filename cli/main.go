package main

import (
	"cli/cmd"
	"cli/ui"
	"os"
)

func main() {
	err := cmd.Execute()

	if err != nil {
		ui.PrintBlock(err)
		os.Exit(1)
	}
}
