package ui

import (
	"cli/env"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
