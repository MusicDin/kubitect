package utils

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	red    = color.New(color.FgHiRed).SprintFunc()
	yellow = color.New(color.FgHiYellow).SprintFunc()
	blue   = color.New(color.FgHiCyan).SprintFunc()
)

func PrintError(msg ...any) {
	printStamp(red("ERROR"), msg)
}

func PrintWarning(msg ...any) {
	printStamp(yellow("WARNING"), msg)
}

func PrintDebug(msg ...any) {
	printStamp(blue("DEBUG"), msg)
}

func printStamp(stamp string, msg []any) {

	stamp = fmt.Sprintf("[ %v ] ", stamp)

	if len(msg) == 0 {
		fmt.Println(stamp)
		return
	}

	format := stamp + fmt.Sprint(msg[0])
	args := msg[1:]

	fmt.Println(fmt.Sprintf(format, args...))
}
