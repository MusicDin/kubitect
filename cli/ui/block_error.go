package ui

type ErrorBlock struct {
	BasicBlock

	Severity Level
}

func (e ErrorBlock) Error() string {
	return e.Format(nil, Colors.NONE)
}

func NewErrorBlock(level Level, content []Content) ErrorBlock {
	return ErrorBlock{
		Severity:   level,
		BasicBlock: BasicBlock{content},
	}
}

// NewErrorLine contains a title and lines to be printed within
// a block. When formatted, a colored title is printed in the
// same line as title.
func NewErrorLine(title string, lines ...string) Content {
	return Content{
		title:   title,
		lines:   lines,
		compact: true,
	}
}

// NewErrorSection contains a title and lines to be printed
// within a block. When formatted, a colored title and lines
// are printed each in a new line. Lines are also additionally
// indented.
func NewErrorSection(title string, lines ...string) Content {
	return Content{
		title:         title,
		lines:         lines,
		linesIndent:   2,
		linesRequired: true,
	}
}
