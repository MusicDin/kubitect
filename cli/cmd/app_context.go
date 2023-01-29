package cmd

import (
	"cli/ui"
	"path"
	"path/filepath"
)

type AppContext struct {
	workingDir string
	homeDir    string
	local      bool
	showTfPlan bool
	ui         *ui.Ui
}

func (c *AppContext) Local() bool {
	return c.local
}

func (c *AppContext) ShowTerraformPlan() bool {
	return c.showTfPlan
}

func (c *AppContext) WorkingDir() string {
	return c.workingDir
}

func (c *AppContext) HomeDir() string {
	return c.homeDir
}

func (c *AppContext) ShareDir() string {
	return path.Join(c.homeDir, DefaultShareDir)
}

func (c *AppContext) ClustersDir() string {
	return filepath.Join(c.homeDir, DefaultClustersDir)
}

func (c *AppContext) LocalClustersDir() string {
	return filepath.Join(c.workingDir, DefaultHomeDir, DefaultClustersDir)
}

func (c *AppContext) Ui() *ui.Ui {
	return c.ui
}
