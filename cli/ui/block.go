package ui

import (
	"fmt"
	"strings"
)

// Error line symbols
const (
	lineInitial = "\u250C" // .-
	lineMiddle  = "\u2502" // |
	lineFinal   = "\u2514" // '-
)

type Block interface {
	Format(*OutputStream, Level)
}

type BasicBlock struct {
	content []Content
}

func (b BasicBlock) Format(stream *OutputStream, level Level) string {
	var lines Lines

	color := level.Color()
	prefix := lineMiddle + " "
	indent := len([]rune(prefix))

	for _, c := range b.content {
		for _, l := range c.format(stream, color, indent) {
			lines.Append(color(prefix) + l)
		}
	}

	lines.Prepend(color(lineInitial))
	lines.Append(color(lineFinal))

	return strings.Join(lines, "\n")
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

func (c Content) format(s *OutputStream, color Color, indent int) []string {
	if c.linesRequired && len(c.lines) == 0 {
		return nil
	}

	out, colsLeft := Format(s, c.title, indent, 0)
	out.Color(color)

	if c.compact {
		c.linesIndent = 0
	} else {
		colsLeft = 0
	}

	prefix := fmt.Sprintf("%*s", c.linesIndent, "")

	for i, l := range c.lines {
		lines, _ := Format(s, l, len(prefix)+indent, colsLeft)
		lines.Prefix(prefix)

		if c.compact && i == 0 {
			out.Glue(lines)
			continue
		}

		out.Append(lines...)
	}

	return out
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
