package app

import (
	"cli/ui"
	"path"
	"testing"
)

type (
	AppContextMock interface {
		AppContext
		Ui() ui.UiMock
	}

	appContextMock struct {
		appContext
		ui ui.UiMock
	}
)

func (m *appContextMock) Ui() ui.UiMock {
	return m.ui
}

func MockAppContext(t *testing.T, opts ...AppContextOptions) AppContextMock {
	t.Helper()
	tmpDir := t.TempDir()

	var o AppContextOptions

	if len(opts) > 0 {
		o = opts[0]
	}

	ctx := appContext{
		workingDir: tmpDir,
		homeDir:    tmpDir,
		local:      o.Local,
		showTfPlan: o.ShowTerraformPlan,
	}

	if !o.Local {
		ctx.homeDir = path.Join(tmpDir, "home")
	}

	uOpts := ui.UiOptions{
		AutoApprove: o.AutoApprove,
		Debug:       o.Debug,
		NoColor:     o.NoColor,
	}

	u := ui.MockGlobalTerminalUi(t, uOpts)

	return &appContextMock{ctx, u}
}
