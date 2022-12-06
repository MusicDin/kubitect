package ui

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUi_Global(t *testing.T) {
	ui := GlobalUi()

	assert.NotNil(t, ui)
	assert.NotNil(t, ui.Streams.In.File)
	assert.NotNil(t, ui.Streams.Out.File)
	assert.NotNil(t, ui.Streams.Err.File)

	if ui != GlobalUi() {
		t.Error("GlobalUi is not a singleton")
	}
}

func TestUiOptions(t *testing.T) {
	opts := UiOptions{
		Debug:       true,
		AutoApprove: true,
		NoColor:     true,
	}

	ui := NewUi(opts)

	assert.NotNil(t, ui)
	assert.NotNil(t, ui.Streams.In.File)
	assert.NotNil(t, ui.Streams.Out.File)
	assert.NotNil(t, ui.Streams.Err.File)
	assert.True(t, ui.autoApprove)
	assert.True(t, ui.NoColor)
	assert.True(t, ui.Debug)
}

func TestUi_AskAutoApprove(t *testing.T) {
	ui := MockUi(t)
	ui.autoApprove = true

	assert.NoError(t, ui.Ask())
}

func TestUi_Print(t *testing.T) {
	ui := MockUi(t)

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
	ui := MockUi(t)

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
