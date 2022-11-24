package main

import (
	"cli/cmd"
	"cli/ui"
	"os"
)

func main() {
	err := cmd.Execute()

	if err != nil {
		ui.GlobalUi().PrintBlockE(err)
		os.Exit(1)
	}
}
