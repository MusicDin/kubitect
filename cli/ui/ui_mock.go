package ui

import (
	"io/ioutil"
	"os"
	"testing"
)

func mockUi(t *testing.T) *Ui {
	t.Helper()

	return &Ui{
		Streams: &Streams{
			Out: mockTerminalStream(t, "stdout"),
			Err: mockTerminalStream(t, "stderr"),
			In:  mockInputStream(t),
		},
	}
}

func mockTerminalStream(t *testing.T, fName string) *OutputStream {
	t.Helper()

	isTerm := func(*os.File) bool {
		return true
	}

	cols := func(*os.File) int {
		return 10
	}

	return &OutputStream{
		File:       tmpFile(t, fName),
		isTerminal: isTerm,
		columns:    cols,
	}
}

func mockNonTerminalStream(t *testing.T, fName string) *OutputStream {
	t.Helper()

	isTerm := func(*os.File) bool {
		return false
	}

	cols := func(*os.File) int {
		return 10
	}

	return &OutputStream{
		File:       tmpFile(t, fName),
		isTerminal: isTerm,
		columns:    cols,
	}
}

func mockInputStream(t *testing.T) *InputStream {
	t.Helper()

	isTerm := func(*os.File) bool {
		return false
	}

	return &InputStream{
		File:       tmpFile(t, "stdin"),
		isTerminal: isTerm,
	}
}

func tmpFile(t *testing.T, name string) *os.File {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), name)

	if err != nil {
		t.Errorf("failed creating tmp file (%s): %v", name, err)
	}

	return f
}

func readFile(t *testing.T, f *os.File) string {
	t.Helper()

	if f == nil {
		t.Errorf("failed to read nil stream: %v", f.Name())
	}

	file, err := ioutil.ReadFile(f.Name())

	if err != nil {
		t.Errorf("failed reading tmp file (%s): %v", f.Name(), err)
	}

	return string(file)
}
