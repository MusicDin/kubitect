package utils

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

func PrintError(str ...string) {
	printStamp(red("ERROR"), str)
}

func PrintWarning(str ...string) {
	printStamp(yellow("WARNING"), str)
}

func printStamp(stamp string, str []string) {
	fmt.Printf("[ %v ] %s\n", stamp, strings.Join(str, " "))
}
