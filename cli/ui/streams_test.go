package ui

import (
	"os"
	"testing"
)

func TestStreams_Standard(t *testing.T) {
	s := StandardStreams()

	if s.In.File != os.Stdin {
		t.Error("stdin of initialized GlobalUi does not point to os.Stdin")
	}

	if s.Out.File != os.Stdout {
		t.Error("stdout of initialized GlobalUi does not point to os.Stdout")
	}

	if s.Err.File != os.Stderr {
		t.Error("stderr of initialized GlobalUi does not point to os.Stderr")
	}
}

func TestStreams_Empty(t *testing.T) {
	s := Streams{
		In: &InputStream{
			isTerminal: nil,
		},
		Out: &OutputStream{
			isTerminal: nil,
			columns:    nil,
		},
		Err: &OutputStream{
			isTerminal: nil,
			columns:    nil,
		},
	}

	if s.In.IsTerminal() {
		t.Error("isTerminal produced wrong output (expected: false)")
	}

	if s.Out.IsTerminal() {
		t.Error("isTerminal produced wrong output (expected: false)")
	}

	if s.Out.Columns() != defaultColumns {
		t.Errorf("columns produced wrong output (expected: %d)", defaultColumns)
	}

	if s.Out.IsTerminal() {
		t.Error("isTerminal produced wrong output (expected: false)")
	}

	if s.Out.Columns() != defaultColumns {
		t.Errorf("columns produced wrong output (expected: %d)", defaultColumns)
	}
}

func TestStreams_Terminal(t *testing.T) {
	s := MockUi(t).Streams

	if isTerminal(s.Out.File) {
		t.Error("isTerminal produced wrong output (expected: false)")
	}

	if columns(s.Out.File) != defaultColumns {
		t.Errorf("columns produced wrong output (expected: %d)", defaultColumns)
	}
}
