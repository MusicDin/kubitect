package ui

import (
	"strings"

	"github.com/MusicDin/kubitect/pkg/ui/streams"
)

// format formats given string into multiple lines that fit the
// width (columns) of the output stream.
// Argument startAt defines where in the current row should the
// string be printed from.
//
// It returns formatted lines and index of the current column.
func Format(o streams.OutputStream, str string, indent, startAt int) (Lines, int) {
	if o == nil || !o.IsTerminal() {
		lines := strings.Split(str, "\n")

		if startAt > 0 {
			lines[0] = " " + lines[0]
		}

		return lines, 1
	}

	var lines []string
	var pivot int

	width := o.Columns()

	for _, line := range strings.Split(str, "\n") {
		var ls []string
		ls, pivot = fmtLine(line, width-indent, startAt)
		lines = append(lines, ls...)
		startAt = 0
	}

	return lines, pivot
}

// fmtLine formats a message according to the given width.
// If word cannot fit into a current line, it tries to fit
// it into a new line. If word is still to long, it writes
// the word character by character.
//
// It returns formatted message and index of the current
// column (pivot).
func fmtLine(line string, width int, startAt int) ([]string, int) {
	var out string

	cw := width - startAt // current line width

	if cw <= 0 {
		cw = width
	}

	for _, s := range strings.Split(line, " ") {
		sw := len(s)

		// add space
		if 0 < cw && cw < width {
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
		if sw <= width {
			out += "\n" + s
			cw = width - sw
			continue
		}

		// type word char by char
		for _, c := range s {
			if cw < 1 {
				out += "\n"
				cw = width
			}

			out += string(c)
			cw--
		}
	}

	return strings.Split(out, "\n"), (width - cw)
}
