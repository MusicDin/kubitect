package actions

import (
	"cli/env"
	"cli/tools/ansible"
	"cli/tools/git"
	"cli/ui"
	"fmt"
	"os"
	"path/filepath"
)

type ApplyAction string

const (
	UNKNOWN ApplyAction = "unknown"
	CREATE  ApplyAction = "create"
	UPGRADE ApplyAction = "upgrade"
	SCALE   ApplyAction = "scale"
)

func (a ApplyAction) String() string {
	return string(a)
}

func ToApplyAction(a string) (ApplyAction, error) {
	switch a {
	case CREATE.String(), "":
		return CREATE, nil
	case UPGRADE.String():
		return UPGRADE, nil
	case SCALE.String():
		return SCALE, nil
	default:
		return UNKNOWN, fmt.Errorf("unknown cluster action: %s", a)
	}
}

func (c *Cluster) Apply(action string) error {
	a, err := ToApplyAction(action)

	if err != nil {
		return err
	}

	if c.AppliedConfig == nil && (a == SCALE || a == UPGRADE) {
		fmt.Printf("Cannot %s cluster '%s'. It has not been created yet.\n\n", a, c.Name)

		if err := ui.Ask("Would you like to create it instead?"); err != nil {
			return err
		}

		a = CREATE
	}

	var events Events

	if c.AppliedConfig != nil {
		events, err = c.plan(a)

		if err != nil {
			return err
		}

		if len(events) == 0 {
			fmt.Println("No changes detected.")
			return nil
		}
	}

	if err := prepare(c); err != nil {
		return err
	}

	switch a {
	case CREATE:
		err = create(c)
	case UPGRADE:
		err = upgrade(c)
	case SCALE:
		err = scale(c, events)
	}

	if err != nil {
		return err
	}

	return c.ApplyNewConfig()
}

// prepare prepares cluster's directory. It ensures that Kubitect project
// files are present in the directory, new configuration file is stored in
// the temporary location and that main virtual environment is created.
func prepare(c *Cluster) error {
	if err := initCluster(c); err != nil {
		return err
	}

	if err := c.StoreNewConfig(); err != nil {
		return err
	}

	if err := c.SetupMainVE(); err != nil {
		return err
	}

	if err := ansible.KubitectInit(c.Path, ansible.INIT); err != nil {
		return err
	}

	return ansible.KubitectHostsSetup(c.Path)
}

// initCluster ensures cluster directory exists and all required files are
// copied from the Kubitect git project. If local flag is used, project
// files are copied from the current directory.
func initCluster(c *Cluster) error {
	cfg := c.NewConfig

	url := env.ConstProjectUrl
	version := env.ConstProjectVersion

	if cfg.Kubitect.Url != nil {
		url = string(*cfg.Kubitect.Url)
	}

	if cfg.Kubitect.Version != nil {
		version = string(*cfg.Kubitect.Version)
	}

	ui.Printf(ui.DEBUG, "kubitect.url: %s\n", url)
	ui.Printf(ui.DEBUG, "kubitect.version: %s\n", version)

	if c.Local {
		return copyReqFiles(c.Ctx.WorkingDir(), c.Path)
	}

	tmpDir := filepath.Join(c.Path, "tmp")

	if err := os.RemoveAll(tmpDir); err != nil {
		return err
	}

	if err := git.Clone(tmpDir, url, version); err != nil {
		return err
	}

	if err := copyReqFiles(tmpDir, c.Path); err != nil {
		return err
	}

	return os.RemoveAll(tmpDir)
}
