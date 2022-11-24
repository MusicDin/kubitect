package ui

import (
	"cli/env"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockUi(t *testing.T) *Ui {
	return &Ui{
		Streams: &Streams{
			Out: mockTerminalStream(t, "stdout"),
			Err: mockTerminalStream(t, "stderr"),
			In:  mockInputStream(t),
		},
	}
}

func mockTerminalStream(t *testing.T, fName string) *OutputStream {
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
	isTerm := func(*os.File) bool {
		return false
	}

	return &InputStream{
		File:       tmpFile(t, "stdin"),
		isTerminal: isTerm,
	}
}

func tmpFile(t *testing.T, name string) *os.File {
	f, err := os.CreateTemp(t.TempDir(), name)

	if err != nil {
		t.Errorf("failed creating tmp file (%s): %v", name, err)
	}

	return f
}

func readFile(t *testing.T, f *os.File) string {
	if f == nil {
		t.Errorf("failed to read nil stream: %v", f.Name())
	}

	file, err := ioutil.ReadFile(f.Name())

	if err != nil {
		t.Errorf("failed reading tmp file (%s): %v", f.Name(), err)
	}

	return string(file)
}

// func TestLevel_Color(t *testing.T) {
// 	assert.Equal(t, Colors.NONE, INFO.Color())
// 	assert.Equal(t, Colors.NONE, DEBUG.Color())
// 	assert.Equal(t, Colors.YELLOW, WARN.Color())
// 	assert.Equal(t, Colors.RED, ERROR.Color())

// 	env.NoColor = true
// 	assert.Equal(t, Colors.NONE, WARN.Color())
// 	assert.Equal(t, Colors.NONE, ERROR.Color())
// 	env.NoColor = false
// }

func TestUi_Global(t *testing.T) {
	if GlobalUi() != GlobalUi() {
		t.Error("GlobalUi is not a singleton")
	}
}

func TestUi_AskAutoApprove(t *testing.T) {
	ui := mockUi(t)

	assert.NoError(t, ui.Ask())

	env.AutoApprove = true
	err := ui.Ask()
	env.AutoApprove = false
	assert.NoError(t, err)
}

func TestUi_Print(t *testing.T) {
	ui := mockUi(t)

	ui.Print(DEBUG, "not-printed")

	ui.Print(INFO, "test1\n")
	ui.Printf(INFO, "%s\n", "test2")
	ui.Println(INFO, "test3")

	ui.Print(WARN, "test1\n")
	ui.Printf(ERROR, "%s\n", "test2")
	ui.Println(ERROR, "test3")

	out := readFile(t, ui.Streams.Out.File)
	err := readFile(t, ui.Streams.Err.File)

	assert.Equal(t, "test1\ntest2\ntest3\n", string(out))
	assert.Equal(t, "test1\ntest2\ntest3\n", string(err))
}

func TestUi_PrintBlock(t *testing.T) {
	ui := mockUi(t)

	expect := "┌\n│ Error: \n│ test\n└\n"

	eb := NewErrorBlock(ERROR,
		[]Content{
			NewErrorLine("Error:", "test"),
		},
	)

	ui.PrintBlockE(fmt.Errorf("test"))
	ui.PrintBlockE(eb)

	err := readFile(t, ui.Streams.Err.File)

	assert.Equal(t, expect+expect, string(err))

}
