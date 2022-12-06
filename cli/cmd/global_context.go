package cmd

import (
	"cli/ui"
	"path"
	"path/filepath"
)

type GlobalContext struct {
	workingDir string
	homeDir    string
	local      bool
	showTfPlan bool
	ui         *ui.Ui
}

func (c *GlobalContext) Local() bool {
	return c.local
}

func (c *GlobalContext) ShowTerraformPlan() bool {
	return c.showTfPlan
}

func (c *GlobalContext) WorkingDir() string {
	return c.workingDir
}

func (c *GlobalContext) HomeDir() string {
	return c.homeDir
}

func (c *GlobalContext) ShareDir() string {
	return path.Join(c.homeDir, DefaultShareDir)
}

func (c *GlobalContext) ClustersDir() string {
	return filepath.Join(c.homeDir, DefaultClustersDir)
}

func (c *GlobalContext) LocalClustersDir() string {
	return filepath.Join(c.workingDir, DefaultHomeDir, DefaultClustersDir)
}

func (c *GlobalContext) Ui() *ui.Ui {
	return c.ui
}
