package cmd

import (
	"cli/env"
	"fmt"
	"path"
	"path/filepath"
)

type AppContext struct {
	workingDir string
	homeDir    string
	local      bool
	showTfPlan bool
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

func (c *AppContext) VerifyRequirements() error {
	var missing []string

	for _, app := range env.ProjectRequiredApps {
		if !AppExists(app) {
			missing = append(missing, app)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("Some requirements are not met: %v", missing)
	}

	return nil
}
