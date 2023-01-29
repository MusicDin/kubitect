package ui

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func colorsEqual(c1, c2 Color) bool {

	cf1 := reflect.ValueOf(c1)
	cf2 := reflect.ValueOf(c2)

	return cf1.Pointer() == cf2.Pointer()
}

func TestUi_Global(t *testing.T) {
	assert.Nil(t, instance)

	ui := GlobalUi()

	assert.NotNil(t, ui)
	assert.NotNil(t, ui.Streams().In().File())
	assert.NotNil(t, ui.Streams().Out().File())
	assert.NotNil(t, ui.Streams().Err().File())

	if ui != GlobalUi() {
		t.Error("GlobalUi is not a singleton")
	}
}

func TestUi_Streams(t *testing.T) {
	MockGlobalUi(t)

	assert.NotNil(t, Streams())
	assert.NotNil(t, Streams().In().File())
	assert.NotNil(t, Streams().Out().File())
	assert.NotNil(t, Streams().Err().File())
}

func TestUiOptions_Default(t *testing.T) {
	MockGlobalUi(t)
	assert.True(t, HasColor())
	assert.False(t, Debug())
	assert.False(t, AutoApprove())
}

func TestUiOptions_Custom(t *testing.T) {
	o := UiOptions{
		AutoApprove: true,
		Debug:       true,
		NoColor:     true,
	}

	ui := MockUi(t, o)
	assert.False(t, ui.HasColor())
	assert.True(t, ui.Debug())
	assert.True(t, ui.AutoApprove())
}

func TestOutputStream_Out(t *testing.T) {
	ui := MockGlobalUi(t)
	s := instance.outputStream(INFO)
	assert.Equal(t, ui.Streams().Out(), s)
}

func TestOutputStream_Err(t *testing.T) {
	ui := MockGlobalUi(t)
	s := instance.outputStream(ERROR)
	assert.Equal(t, ui.Streams().Err(), s)
}

func TestOutputColor(t *testing.T) {
	MockGlobalUi(t)

	debugColor := instance.outputColor(DEBUG)
	infoColor := instance.outputColor(INFO)
	warnColor := instance.outputColor(WARN)
	errColor := instance.outputColor(ERROR)

	assert.True(t, colorsEqual(debugColor, Colors.NONE), "Wrong color output for severity level DEBUG.")
	assert.True(t, colorsEqual(infoColor, Colors.NONE), "Wrong color output for severity level INFO.")
	assert.True(t, colorsEqual(warnColor, Colors.YELLOW), "Wrong color output for severity level WARN.")
	assert.True(t, colorsEqual(errColor, Colors.RED), "Wrong color output for severity level ERROR.")
}

func TestOutputColor_NoColor(t *testing.T) {
	o := UiOptions{
		NoColor: true,
	}

	MockGlobalUi(t, o)

	debugColor := instance.outputColor(DEBUG)
	infoColor := instance.outputColor(INFO)
	warnColor := instance.outputColor(WARN)
	errColor := instance.outputColor(ERROR)

	assert.True(t, colorsEqual(debugColor, Colors.NONE), "Wrong color output for severity level DEBUG (NoColor: true).")
	assert.True(t, colorsEqual(infoColor, Colors.NONE), "Wrong color output for severity level INFO (NoColor: true).")
	assert.True(t, colorsEqual(warnColor, Colors.NONE), "Wrong color output for severity level WARN (NoColor: true).")
	assert.True(t, colorsEqual(errColor, Colors.NONE), "Wrong color output for severity level ERROR (NoColor: true).")
}

func TestUi_Ask(t *testing.T) {
	MockGlobalUi(t)
	assert.NoError(t, Ask())
}

func TestUi_Ask_NonTerminal(t *testing.T) {
	ui := MockUi(t)

	assert.NoError(t, ui.Ask()) // Auto-approve if stdin is not a terminal
}

func TestUi_Ask_Terminal(t *testing.T) {
	ui := MockTerminalUi(t)
	ui.WriteStdin(t, "yes")

	assert.NoError(t, ui.Ask())
}

func TestUi_Ask_TerminalReject(t *testing.T) {
	ui := MockTerminalUi(t)
	ui.WriteStdin(t, "no")

	assert.EqualError(t, ui.Ask(), "User aborted...")
}

func TestUi_Ask_TerminalFail(t *testing.T) {
	ui := MockTerminalUi(t)

	assert.EqualError(t, ui.Ask(), "ask: EOF")
}

func TestUi_Ask_TerminalDefaultQuestion(t *testing.T) {
	ui := MockGlobalTerminalUi(t)
	ui.WriteStdin(t, "yes")

	assert.NoError(t, ui.Ask())
	assert.Equal(t, "\nWould you like to continue? (yes/no) ", ui.ReadStdout(t))
}

func TestUi_Ask_TerminalCustomQuestion(t *testing.T) {
	ui := MockTerminalUi(t)
	ui.WriteStdin(t, "yes")

	assert.NoError(t, ui.Ask("Test", "1", "23"))
	assert.Equal(t, "\nTest 1 23 (yes/no) ", ui.ReadStdout(t))
}

func TestUi_Ask_AutoApprove(t *testing.T) {
	o := UiOptions{
		AutoApprove: true,
	}

	ui := MockUi(t, o)

	assert.NoError(t, ui.Ask())
}

func TestUi_Print(t *testing.T) {
	ui := MockGlobalUi(t)
	Print(INFO, "test", "2")

	assert.Equal(t, "test2", ui.ReadStdout(t))
}

func TestUi_Print_Debug(t *testing.T) {
	o := UiOptions{
		Debug: true,
	}

	ui := MockUi(t, o)
	ui.Print(DEBUG, "test", "2")

	assert.Equal(t, "test2", ui.ReadStdout(t))
}

func TestUi_Print_NoDebug(t *testing.T) {
	ui := MockUi(t)
	ui.Print(DEBUG, "test", "2")

	assert.Equal(t, "", ui.ReadStdout(t))
}

func TestUi_Println(t *testing.T) {
	ui := MockGlobalUi(t)
	Println(INFO, "test", "test2")

	assert.Equal(t, "testtest2\n", ui.ReadStdout(t))
}

func TestUi_Printf(t *testing.T) {
	ui := MockGlobalUi(t)
	Printf(INFO, "%s:%d", "test", 2)

	assert.Equal(t, "test:2", ui.ReadStdout(t))
}

func TestUi_PrintBlockE(t *testing.T) {
	ui := MockGlobalUi(t)

	expect := "┌\n│ Error: test\n└\n"

	eb := NewErrorBlock(ERROR,
		[]Content{
			NewErrorLine("Error:", "test"),
		},
	)

	PrintBlockE(fmt.Errorf("test"))
	PrintBlockE(eb)

	err := ui.ReadStderr(t)
	assert.Equal(t, expect+expect, string(err))
}
