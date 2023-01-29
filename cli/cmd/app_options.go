package cmd

import (
	"cli/ui"
	"os"
	"path/filepath"
)

const (
	EnvHome            = "KUBITECT_HOME"
	DefaultHomeDir     = ".kubitect"
	DefaultShareDir    = "share"
	DefaultClustersDir = "clusters"
)

type AppOptions struct {
	// Automatically approve user prompts
	AutoApprove bool

	// // Show debug messages
	Debug bool

	// // Disable color in output
	NoColor bool

	// Local deployment. Use working dir as project home dir.
	Local bool

	// Show terraform plan
	ShowTerraformPlan bool
}

func (o *AppOptions) AppContext() *AppContext {
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	hd := filepath.Join(wd, DefaultHomeDir)

	if !o.Local {
		userHomeDir, err := os.UserHomeDir()

		if err != nil {
			panic(err)
		}

		def := filepath.Join(userHomeDir, DefaultHomeDir)
		hd = EnvVar(EnvHome, def)
	}

	uiOpts := ui.UiOptions{
		Debug:       o.Debug,
		NoColor:     o.NoColor,
		AutoApprove: o.AutoApprove,
	}

	// Initialize global ui
	ui.GlobalUi(uiOpts)

	return &AppContext{
		homeDir:    hd,
		workingDir: wd,
		local:      o.Local,
		showTfPlan: o.ShowTerraformPlan,
	}
}
