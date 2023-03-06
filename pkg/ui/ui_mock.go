package ui

import (
	"os"
	"testing"

	"github.com/MusicDin/kubitect/pkg/ui/streams"

	"github.com/stretchr/testify/assert"
)

type (
	UiMock interface {
		Ui

		WriteStdin(t *testing.T, content string)
		ReadStdout(t *testing.T) string
		ReadStderr(t *testing.T) string
	}

	uiMock struct {
		*ui

		outOffset int64
		errOffset int64
	}
)

func setUiOptions(ui *uiMock, uiOptions ...UiOptions) *uiMock {
	var o UiOptions

	if len(uiOptions) > 0 {
		o = uiOptions[0]
	}

	ui.autoApprove = o.AutoApprove
	ui.debug = o.Debug
	ui.noColor = o.NoColor

	return ui
}

func MockUi(t *testing.T, uiOptions ...UiOptions) UiMock {
	t.Helper()

	mock := uiMock{
		ui: &ui{
			streams: streams.MockStreams(t),
		},
	}

	return setUiOptions(&mock, uiOptions...)
}

func MockTerminalUi(t *testing.T, uiOptions ...UiOptions) UiMock {
	t.Helper()

	mock := uiMock{
		ui: &ui{
			streams: streams.MockTerminalStreams(t),
		},
	}

	return setUiOptions(&mock, uiOptions...)
}

func MockGlobalUi(t *testing.T, uiOptions ...UiOptions) UiMock {
	t.Helper()

	mock := uiMock{
		ui: &ui{
			streams: streams.MockStreams(t),
		},
	}

	instance = mock.ui

	return setUiOptions(&mock, uiOptions...)
}

func MockGlobalTerminalUi(t *testing.T, uiOptions ...UiOptions) UiMock {
	t.Helper()

	instance = &ui{
		streams: streams.MockTerminalStreams(t),
	}

	mock := uiMock{
		ui: instance,
	}

	return setUiOptions(&mock, uiOptions...)
}

// ReadStderr reads output stream file content since last read.
func (ui *uiMock) ReadStdout(t *testing.T) string {
	t.Helper()

	f := ui.Streams().Out().File()
	f.Seek(ui.errOffset, 0)

	bytes, err := os.ReadFile(f.Name())
	assert.NoError(t, err)
	f.Seek(0, 1)

	ui.errOffset += int64(len(bytes))

	return string(bytes)
}

// ReadStderr reads error stream file content since last read.
func (ui *uiMock) ReadStderr(t *testing.T) string {
	t.Helper()

	f := ui.Streams().Err().File()
	f.Seek(ui.outOffset, 0)

	bytes, err := os.ReadFile(f.Name())
	assert.NoError(t, err)
	f.Seek(0, 1)

	ui.outOffset += int64(len(bytes))

	return string(bytes)
}

// WriteStdin write given content to the input stream.
func (ui *uiMock) WriteStdin(t *testing.T, content string) {
	t.Helper()

	err := os.Truncate(ui.Streams().In().File().Name(), 0)
	assert.NoError(t, err)

	_, err = ui.Streams().In().File().WriteString(content)
	assert.NoError(t, err)

	_, err = ui.Streams().In().File().Seek(0, 0)
	assert.NoError(t, err)
}
