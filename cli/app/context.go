package app

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/MusicDin/kubitect/cli/env"
	"github.com/MusicDin/kubitect/cli/ui"
)

const (
	defaultHomeDir     = ".kubitect"
	defaultShareDir    = "share"
	defaultClustersDir = "clusters"
)

type AppContextOptions struct {
	// Automatically approve user prompts
	AutoApprove bool

	// Show debug messages
	Debug bool

	// Disable color in output
	NoColor bool

	// Local deployment. Use working dir as project home dir.
	Local bool

	// Show terraform plan
	ShowTerraformPlan bool

	// AppContext instance.
	appContext AppContext
}

type (
	AppContext interface {
		// WorkingDir returns path of the directory from where
		// the application was executed.
		WorkingDir() string

		// HomeDir returns path of the directory where cluster
		// configuration files and related dependencies are
		// stored.
		//
		// Default is "$HOME/.kubitect".
		HomeDir() string

		// ShareDir returns path of the directory where binaries
		// that are shared among all clusters are stored.
		//
		// Default is "$HOME/.kubitect/share"
		ShareDir() string

		// ClustersDir returns path of the directory where
		// clusters are created.
		//
		// Default is "$HOME/.kubitect/clusters".
		ClustersDir() string

		// LocalClustersDir returns path of the directory where
		// local clusters are created.
		//
		// Default is "./.kubitect/clusters".
		LocalClustersDir() string

		// Local indicates that all cluster actions should be
		// executed within working directory.
		Local() bool

		// ShowTerraformPlan indicates that terraform plan should
		// be always shown.
		ShowTerraformPlan() bool
	}

	appContext struct {
		workingDir string
		homeDir    string
		local      bool
		showTfPlan bool
	}
)

// NewAppContext creates new application context and initializes
// a global UI.
func (o AppContextOptions) AppContext() AppContext {
	if o.appContext != nil {
		return o.appContext
	}

	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	home := filepath.Join(wd, defaultHomeDir)

	if !o.Local {
		userHomeDir, err := os.UserHomeDir()

		if err != nil {
			panic(err)
		}

		home = filepath.Join(userHomeDir, defaultHomeDir)
	}

	uiOpts := ui.UiOptions{
		Debug:       o.Debug,
		NoColor:     o.NoColor,
		AutoApprove: o.AutoApprove,
	}

	// Initialize global ui
	ui.GlobalUi(uiOpts)

	o.appContext = &appContext{
		homeDir:    home,
		workingDir: wd,
		local:      o.Local,
		showTfPlan: o.ShowTerraformPlan,
	}

	return o.appContext
}

func (c *appContext) Local() bool {
	return c.local
}

func (c *appContext) ShowTerraformPlan() bool {
	return c.showTfPlan
}

func (c *appContext) WorkingDir() string {
	return c.workingDir
}

func (c *appContext) HomeDir() string {
	return c.homeDir
}

func (c *appContext) ShareDir() string {
	return path.Join(c.homeDir, defaultShareDir)
}

func (c *appContext) ClustersDir() string {
	return filepath.Join(c.homeDir, defaultClustersDir)
}

func (c *appContext) LocalClustersDir() string {
	return filepath.Join(c.workingDir, defaultHomeDir, defaultClustersDir)
}

func (c appContext) VerifyRequirements() error {
	var missing []string

	for _, app := range env.ProjectRequiredApps {
		if !appExists(app) {
			missing = append(missing, app)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("Some requirements are not met: %v", missing)
	}

	return nil
}

// appExists returns true if command with
// a given name is found in PATH.
func appExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
