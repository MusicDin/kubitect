package ui

import (
	"cli/env"
	"fmt"
	"strings"
)

// Error line symbols
const (
	lineInitial = "\u250C" // .-
	lineMiddle  = "\u2502" // |
	lineFinal   = "\u2514" // '-
)

func (t Level) Color() Color {
	switch t {
	case WARN:
		return Colors.YELLOW
	case ERROR:
		return Colors.RED
	default:
		return Colors.NONE
	}
}

type Lines []string

func (l *Lines) Append(lines ...string) {
	*l = append(*l, lines...)
}

func (l *Lines) Prepend(lines ...string) {
	*l = append(lines, *l...)
}

func (l *Lines) Color(c Color) {
	for i := range *l {
		(*l)[i] = c((*l)[i])
	}
}

func (l *Lines) Prefix(p string) {
	for i := range *l {
		(*l)[i] = p + (*l)[i]
	}
}

func (l *Lines) Glue(lines Lines) {
	x := len(*l) - 1
	(*l)[x] = (*l)[x] + lines[0]
	l.Append(lines[1:]...)
}

type Block struct {
	// color      Color
	content []Content
}

type Content struct {
	title string
	lines Lines

	// Print each line into a new Indents each line
	linesIndent int

	// When set to true, title is not formatted if len(lines) == 0.
	linesRequired bool

	// Lines are printed directly after title (without new line).
	// When compact is set to true, indentation is ignored.
	compact bool
}

func NewContent(title string, lines ...string) Content {
	return Content{
		title: title,
		lines: lines,
	}
}

// format formats the content into lines that fit the
// width of the terminal.
func (c Content) format(color Color, indent int) []string {
	if c.linesRequired && len(c.lines) == 0 {
		return nil
	}

	out, colsLeft := format(c.title, indent, 0)
	out.Color(color)

	if c.compact {
		c.linesIndent = 0
	} else {
		colsLeft = 0
	}

	prefix := fmt.Sprintf("%*s", c.linesIndent, "")

	for i, l := range c.lines {
		lines, _ := format(l, len(prefix)+indent, colsLeft)
		lines.Prefix(prefix)

		if c.compact && i == 0 {
			out.Glue(lines)
			continue
		}

		out = append(out, lines...)
	}

	return out
}

// // TODO
func (e ErrorBlock) Error() string {
	// var out string

	// c := e.Level.Color()

	// for _, cont := range e.Content {
	// 	out += cont.format(c, 2)
	// }

	// TODO: c.format(Colors.NONE, 0)

	return "" //fmt.Sprintf("%s\n%s%s", c(lineInitial), out, c(lineFinal))
}

// ErrorBlock contains multiple ErrorContent objects.
// When formatted, a "block symbols" are prepended
// to each line to form a block.
type ErrorBlock struct {
	Level   Level
	Content []Content
}

func (e ErrorBlock) Format() string {
	var lines Lines

	color := e.Level.Color()

	if env.NoColor {
		color = Colors.NONE
	}

	prefix := lineMiddle + " "
	indent := len([]rune(prefix))

	for _, c := range e.Content {
		for _, l := range c.format(color, indent) {
			lines.Append(color(prefix) + l)
		}
	}

	lines.Prepend(color(lineInitial))
	lines.Append(color(lineFinal))

	return strings.Join(lines, "\n")
}

// BlockLine is contains a title and a line. When formatted,
// a title is colored and line is printed in the same line
// as title.
func NewErrorLine(title string, lines ...string) Content {
	return Content{
		title:   title,
		lines:   lines,
		compact: true,
	}
}

// BlockSection contains a title and multiple lines. When
// formatted, a title is colored and lines are printed in
// a new line with additional indentation.
func NewErrorSection(title string, lines ...string) Content {
	return Content{
		title:         title,
		lines:         lines,
		linesIndent:   2,
		linesRequired: true,
	}
}

// format formats given string into multiple lines that fit the
// width of the output stream. Argument startAt defines where the
//
func format(str string, indent, startAt int) (Lines, int) {
	var lines Lines

	width := Ui().streams.Out.Columns()

	if width <= indent {
		width = defaultColumns
	}

	for _, m := range strings.Split(str, "\n") {
		ls, colsLeft := fmtLine(m, width-indent, startAt)
		lines = append(lines, ls...)
		startAt = colsLeft
	}

	return lines, startAt
}

// fmtLine formats a message according to the given width.
// If word cannot fit into a current line, it tries to fit
// it into a new line. If word is still to long, it writes
// the word character by character.
//
// It returns formatted message and number of columns left
// in a line.
func fmtLine(line string, width int, startAt int) ([]string, int) {
	var out string

	lw := width   // total line width
	cw := startAt // current line width

	if cw <= 0 {
		cw = lw
	}

	for _, s := range strings.Split(line, " ") {
		sw := len(s)

		// add space
		if 0 < cw && cw < lw {
			out += " "
			cw -= 1
		}

		// word fits into current line
		if sw <= cw {
			out += s
			cw -= sw
			continue
		}

		// word fits into new line
		if sw <= lw {
			out += "\n" + s
			cw = lw - sw
			continue
		}

		// type word char by char
		for _, c := range s {
			if cw < 1 {
				out += "\n"
				cw = lw
			}

			out += string(c)
			cw--
		}
	}

	return strings.Split(out, "\n"), cw
}
