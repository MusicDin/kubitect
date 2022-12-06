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

type GenericOptions struct {
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

func (o *GenericOptions) GlobalContext() *GlobalContext {
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

	return &GlobalContext{
		homeDir:    hd,
		workingDir: wd,
		local:      o.Local,
		showTfPlan: o.ShowTerraformPlan,
		ui:         ui.NewUi(uiOpts),
	}
}
