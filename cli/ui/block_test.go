package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlock_Basic(t *testing.T) {
	b := BasicBlock{
		content: []Content{
			{
				title: "Title:",
				lines: []string{"Line"},
			},
		},
	}

	assert.Equal(t, "┌\n│ Title:\n│ Line\n└", b.Format(nil, INFO))
}

func TestBlock_LinesIndent(t *testing.T) {
	b := BasicBlock{
		content: []Content{
			{
				title:       "Title:",
				lines:       []string{"Line"},
				linesIndent: 2,
			},
		},
	}

	assert.Equal(t, "┌\n│ Title:\n│   Line\n└", b.Format(nil, INFO))
}

func TestBlock_LinesRequired(t *testing.T) {
	b := BasicBlock{
		content: []Content{
			{
				title:         "Title:",
				linesRequired: true,
			},
		},
	}

	assert.Equal(t, "", b.Format(nil, INFO))
}

func TestBlock_IndentCompact(t *testing.T) {
	b := BasicBlock{
		content: []Content{
			{
				title:   "Title:",
				lines:   []string{"Line"},
				compact: true,
			},
		},
	}

	assert.Equal(t, "┌\n│ Title: Line\n└", b.Format(nil, INFO))
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
