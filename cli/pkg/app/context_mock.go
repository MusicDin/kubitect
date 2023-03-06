package app

import (
	ui2 "github.com/MusicDin/kubitect/cli/pkg/ui"
	"path"
	"testing"
)

type (
	AppContextMock interface {
		AppContext
		Ui() ui2.UiMock
		Options() AppContextOptions
	}

	appContextMock struct {
		appContext
		appContextOptions AppContextOptions
		ui                ui2.UiMock
	}
)

func (m *appContextMock) Ui() ui2.UiMock {
	return m.ui
}

func (m *appContextMock) Options() AppContextOptions {
	return m.appContextOptions
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

	o.appContext = &ctx

	if !o.Local {
		ctx.homeDir = path.Join(tmpDir, "home")
	}

	uOpts := ui2.UiOptions{
		AutoApprove: o.AutoApprove,
		Debug:       o.Debug,
		NoColor:     o.NoColor,
	}

	u := ui2.MockGlobalTerminalUi(t, uOpts)

	return &appContextMock{
		appContext:        ctx,
		appContextOptions: o,
		ui:                u,
	}
}
