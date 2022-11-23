package ui

import (
	"fmt"

	"github.com/fatih/color"
)

type Color func(a ...interface{}) string

var Colors = struct {
	NONE   Color
	RESET  Color
	GREEN  Color
	RED    Color
	YELLOW Color
	BLUE   Color
}{
	NONE:   colorNone,
	RESET:  color.New(color.Reset).SprintFunc(),
	GREEN:  color.New(color.FgHiGreen).SprintFunc(),
	RED:    color.New(color.FgHiRed).SprintFunc(),
	YELLOW: color.New(color.FgHiYellow).SprintFunc(),
	BLUE:   color.New(color.FgHiCyan).SprintFunc(),
}

func colorNone(a ...interface{}) string {
	return fmt.Sprint(a...)
}
