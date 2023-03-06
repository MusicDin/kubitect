package ui

import (
	"fmt"
	"strings"

	"github.com/MusicDin/kubitect/pkg/ui/streams"
)

// Block symbols
const (
	lineInitial = "\u250C" // .-
	lineMiddle  = "\u2502" // |
	lineFinal   = "\u2514" // '-
)

type (
	Block interface {
		Format(streams.OutputStream, Color) string
	}

	block struct {
		content []Content
	}
)

func (b block) Format(stream streams.OutputStream, color Color) string {
	var lines Lines

	prefix := lineMiddle + " "
	indent := len([]rune(prefix))

	for _, c := range b.content {
		for _, l := range c.format(stream, color, indent) {
			lines = lines.append(color(prefix) + l)
		}
	}

	if len(lines) == 0 {
		return ""
	}

	lines = lines.prepend(color(lineInitial))
	lines = lines.append(color(lineFinal))

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

func (c Content) format(s streams.OutputStream, color Color, indent int) []string {
	if c.linesRequired && len(c.lines) == 0 {
		return nil
	}

	out, pivot := Format(s, c.title, indent, 0)
	out = out.setColor(color)

	if c.compact {
		c.linesIndent = 0
	} else {
		pivot = 0
	}

	prefix := fmt.Sprintf("%*s", c.linesIndent, "")

	for i, l := range c.lines {
		lines, _ := Format(s, l, len(prefix)+indent, pivot)
		lines = lines.setPrefix(prefix)

		if c.compact && i == 0 {
			out = out.glue(lines...)
			continue
		}

		out = out.append(lines...)
	}

	return out
}

type Lines []string

func (l Lines) append(elems ...string) Lines {
	return append(l, elems...)
}

func (l Lines) prepend(elems ...string) Lines {
	return append(elems, l...)
}

func (l Lines) glue(elems ...string) Lines {
	if len(l) == 0 {
		return elems
	}

	x := len(l) - 1
	l[x] = l[x] + elems[0]

	return l.append(elems[1:]...)
}

func (l Lines) setColor(c Color) Lines {
	for i := range l {
		l[i] = c(l[i])
	}
	return l
}

func (l Lines) setPrefix(p string) Lines {
	for i := range l {
		l[i] = p + l[i]
	}

	return l
}
