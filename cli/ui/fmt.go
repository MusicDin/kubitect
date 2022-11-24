package ui

import "strings"

// format formats given string into multiple lines that fit the
// width (columns) of the output stream.
// Argument startAt defines where in the current row should the
// string be printed from.
//
// It returns formatted lines and number of columns left in the
// current line.
func Format(o *OutputStream, str string, indent, startAt int) (Lines, int) {
	var lines Lines

	if o == nil {
		lines = strings.Split(str, "\n")

		if startAt > 0 {
			lines[0] = " " + lines[0]
		}

		return lines, 1
	}

	width := o.Columns()

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
