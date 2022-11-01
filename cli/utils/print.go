package utils

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	green  = color.New(color.FgHiGreen).SprintFunc()
	red    = color.New(color.FgHiRed).SprintFunc()
	yellow = color.New(color.FgHiYellow).SprintFunc()
	blue   = color.New(color.FgHiCyan).SprintFunc()
)

func PrintError(msg ...any) {
	printStamp(red("E"), msg)
}

func PrintWarning(msg ...any) {
	printStamp(yellow("W"), msg)
}

func PrintDebug(msg ...any) {
	printStamp(blue("D"), msg)
}

func PrintSuccess(msg ...any) {
	printStamp(green("S"), msg)
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
