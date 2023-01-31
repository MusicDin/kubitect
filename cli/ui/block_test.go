package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlock_Basic(t *testing.T) {
	b := block{
		content: []Content{
			{
				title: "Title:",
				lines: []string{"Line"},
			},
		},
	}

	assert.Equal(t, "┌\n│ Title:\n│ Line\n└", b.Format(nil, Colors.NONE))
}

func TestBlock_LinesIndent(t *testing.T) {
	b := block{
		content: []Content{
			{
				title:       "Title:",
				lines:       []string{"Line"},
				linesIndent: 2,
			},
		},
	}

	assert.Equal(t, "┌\n│ Title:\n│   Line\n└", b.Format(nil, Colors.NONE))
}

func TestBlock_LinesRequired(t *testing.T) {
	b := block{
		content: []Content{
			{
				title:         "Title:",
				linesRequired: true,
			},
		},
	}

	assert.Equal(t, "", b.Format(nil, Colors.NONE))
}

func TestBlock_IndentCompact(t *testing.T) {
	b := block{
		content: []Content{
			{
				title:   "Title:",
				lines:   []string{"Line"},
				compact: true,
			},
		},
	}

	assert.Equal(t, "┌\n│ Title: Line\n└", b.Format(nil, Colors.NONE))
}

func TestBlockError_Line(t *testing.T) {
	b := NewErrorBlock(ERROR,
		[]Content{
			NewErrorLine("Title:", "Line"),
		},
	)

	assert.Equal(t, "┌\n│ Title: Line\n└", b.Error())
}

func TestBlockError_Section(t *testing.T) {
	b := NewErrorBlock(ERROR,
		[]Content{
			NewErrorSection("Title:", "Line"),
		},
	)

	assert.Equal(t, "┌\n│ Title:\n│   Line\n└", b.Error())
}

func TestLines_Append(t *testing.T) {
	lines1 := Lines{"test", "123"}
	lines2 := Lines{"456", "test"}
	expect := Lines{"test", "123", "456", "test"}

	assert.Equal(t, expect, lines1.append(lines2...))
}

func TestLines_Prepend(t *testing.T) {
	lines1 := Lines{"test", "123"}
	lines2 := Lines{"456", "test"}
	expect := Lines{"456", "test", "test", "123"}

	assert.Equal(t, expect, lines1.prepend(lines2...))
}

func TestLines_Glue(t *testing.T) {
	lines1 := Lines{"test", "123"}
	lines2 := Lines{"456", "test"}
	expect := Lines{"test", "123456", "test"}

	assert.Equal(t, expect, lines1.glue(lines2...))
}

func TestLines_Glue_Empty(t *testing.T) {
	lines1 := Lines{}
	lines2 := Lines{}
	expect := Lines{}

	assert.Equal(t, expect, lines1.glue(lines2...))
}

func TestLines_SetPrefix(t *testing.T) {
	lines := Lines{"test", "123"}
	expect := Lines{"atest", "a123"}

	assert.Equal(t, expect, lines.setPrefix("a"))
}

func TestLines_SetColor(t *testing.T) {
	color := Colors.RED
	lines := Lines{"test", "123"}
	expect := Lines{color("test"), color("123")}

	assert.Equal(t, expect, lines.setColor(color))
}
