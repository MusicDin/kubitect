package main

import (
	"cli/cmd"
)

func main() {

	err := cmd.Execute()

	if err != nil {
		// fmt.Fprintln(os.Stderr, "[ ERROR ]", err)
		// os.Exit(1)
	}
}
