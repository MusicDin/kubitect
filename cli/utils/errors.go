package utils

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// Error line symbols
const (
	lineInitial = "\u250C" // .-
	lineMiddle  = "\u2502" // |
	lineFinal   = "\u2514" // '-
)

type LevelType string

const (
	WARNING LevelType = "warning"
	ERROR   LevelType = "error"
)

type Color func(a ...interface{}) string

func (t LevelType) Color() Color {
	switch t {
	case WARNING:
		return yellow
	default:
		return red
	}
}

func FormatError(e error) error {
	s := e.Error()

	if strings.ContainsAny(s, lineInitial) {
		return e
	}

	return NewError(e)
}

type Errors []error

func (errs Errors) Error() string {
	var out []string

	for _, e := range errs {
		out = append(out, fmt.Sprint(e))
	}

	return strings.Join(out, "\n")
}

type Error struct {
	level   LevelType
	content []ErrorContent
}

type ErrorContent interface {
	Format(Color) string
}

type ErrorLine struct {
	title string
	line  string
}

// ErrorLine is ErrorContent that contains a title and an error line.
// When formatted, the title is colored based on the error type.
// Similarly, a vertical line that is prepended to each terminal line
// is colored based on the error type. A line containing line breaks
// is split into several lines accordingly.
func NewErrorLine(title string, line string) ErrorContent {
	return ErrorLine{
		title: title,
		line:  line,
	}
}

type ErrorSection struct {
	title string
	lines []string
}

// ErrorSection is ErrorContent that contains a title and multiple error
// lines. When formatted, the title is colored based on the error type.
// Similarly, a vertical line that is prepended to each terminal line
// is colored based on the error type. Error lines are also prepended with
// an additional indentation.
func NewErrorSection(title string, lines ...string) ErrorContent {
	return ErrorSection{
		title: title,
		lines: lines,
	}
}

func (e Error) Error() string {
	var out string

	c := e.level.Color()

	for _, cont := range e.content {
		out += cont.Format(c)
	}

	return fmt.Sprintf("%s\n%s%s", c(lineInitial), out, c(lineFinal))
}

func (e Error) String() string {
	return e.Error()
}

// NewError returns new error.
func NewError(msg interface{}) error {
	content := []ErrorContent{
		NewErrorLine("Error:", fmt.Sprint(msg)),
	}

	return NewErrorB(content)
}

// NewErrorB returns new error block.
func NewErrorB(content []ErrorContent) error {
	return Error{
		level:   ERROR,
		content: content,
	}
}

// NewWarn returns new warning.
func NewWarn(msg interface{}) error {
	content := []ErrorContent{
		NewErrorLine("Warning:", fmt.Sprint(msg)),
	}

	return NewWarnB(content)
}

// NewWarnB returns new warning block.
func NewWarnB(content []ErrorContent) error {
	return Error{
		level:   WARNING,
		content: content,
	}
}

// Format formats an error section.
func (e ErrorSection) Format(c Color) string {
	if len(e.lines) == 0 {
		return ""
	}

	titleLine := ErrorLine{
		title: e.title,
	}

	out := titleLine.Format(c)

	for _, line := range e.lines {
		for _, s := range strings.Split(line, "\n") {
			out += format(s, c, 2, true)
		}
	}

	return out
}

// Format formats an error line.
func (e ErrorLine) Format(c Color) string {
	if len(e.line) == 0 {
		out := format(e.title, c, 0, true)
		return strings.Replace(out, e.title, c(e.title), 1)
	}

	var out string

	for i, l := range strings.Split(e.line, "\n") {
		if i == 0 && len(e.title) > 0 {
			out = fmt.Sprintf("%s %s", e.title, l)
			out = format(out, c, 0, true)
			out = strings.Replace(out, e.title, c(e.title), 1)
			continue
		}

		out += format(l, c, 0, true)
	}

	return out
}

// format formats a line according to the terminal width.
func format(msg string, c Color, spacing int, newLine bool) string {
	prefix := fmt.Sprintf("%s %*s", c(lineMiddle), spacing, "")
	indent := len([]rune(lineMiddle)) + spacing + 2

	width, _, err := terminal.GetSize(0)

	if err != nil || width < indent+1 {
		width = indent + 1
	}

	lw := width - indent // line width
	cw := lw             // current line width

	out := prefix

	for _, s := range strings.Split(msg, " ") {
		sw := len(s)

		// word fits into current line
		if sw <= cw {
			if cw < lw {
				out += " "
				cw -= 1
			}

			out += s
			cw -= sw
			continue
		}

		// word fits into new line
		if sw <= lw {
			out += "\n" + prefix + s
			cw = lw - sw
			continue
		}

		// type word char by char

		if cw < lw {
			out += " "
			cw -= 1
		}

		for _, r := range s {
			if cw == 0 {
				out += "\n" + prefix
				cw = lw
			}

			out += string(r)
			cw--
		}
	}

	if newLine {
		out += "\n"
	}

	return out
}
