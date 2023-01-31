package streams

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreams_Standard(t *testing.T) {
	s := StandardStreams()
	assert.Equal(t, os.Stdin, s.In().File(), "Standard input stream of does not point to os.Stdin")
	assert.Equal(t, os.Stdout, s.Out().File(), "Standard output stream of does not point to os.Stdout")
	assert.Equal(t, os.Stderr, s.Err().File(), "Standard error stream of does not point to os.Stderr")
}

func TestStreams_Functions(t *testing.T) {
	s := MockStreams(t)
	assert.False(t, isTerminal(s.Out().File()), "isTerminal produced wrong output for output stream")
	assert.Equal(t, defaultColumns, columns(s.Out().File()), "columns produced wrong output for output stream")
}

func TestStreams_NonTerminal(t *testing.T) {
	s := MockStreams(t)
	assert.False(t, s.In().IsTerminal(), "IsTerminal produced wrong output for a non-terminal input stream")
	assert.False(t, s.Out().IsTerminal(), "IsTerminal produced wrong output for a non-terminal output stream")
	assert.False(t, s.Err().IsTerminal(), "IsTerminal produced wrong output for a non-terminal error stream")
	assert.Equal(t, defaultColumns, s.Out().Columns(), "Columns produced wrong output for a non-terminal output stream")
	assert.Equal(t, defaultColumns, s.Err().Columns(), "Columns produced wrong output for a non-terminal error stream")
}

func TestStreams_Terminal(t *testing.T) {
	s := MockTerminalStreams(t)
	assert.True(t, s.In().IsTerminal(), "IsTerminal produced wrong output for a terminal input stream")
	assert.True(t, s.Out().IsTerminal(), "IsTerminal produced wrong output for a terminal output stream")
	assert.True(t, s.Err().IsTerminal(), "IsTerminal produced wrong output for a terminal error stream")
	assert.Equal(t, 42, s.Out().Columns(), "Columns produced wrong output for a terminal output stream")
	assert.Equal(t, 42, s.Err().Columns(), "Columns produced wrong output for a terminal error stream")
}

func TestStreams_Empty(t *testing.T) {
	s := MockEmptyStreams(t)
	assert.False(t, s.In().IsTerminal(), "isTerminal produced wrong output for <nil> input stream")
	assert.False(t, s.Out().IsTerminal(), "isTerminal produced wrong output for <nil> output stream")
	assert.False(t, s.Err().IsTerminal(), "isTerminal produced wrong output for <nil> error stream")
	assert.Equal(t, defaultColumns, s.Out().Columns(), "columns produced wrong output for <nil> output stream")
	assert.Equal(t, defaultColumns, s.Err().Columns(), "columns produced wrong output for <nil> error stream")
}
