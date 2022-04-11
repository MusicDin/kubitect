package main

import (
	"cli/cmd"
	"fmt"
	"os"
)

func main() {

	err := cmd.Execute()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
