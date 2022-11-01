package actions

import (
	"cli/env"
	"cli/tools/ansible"
	"cli/tools/git"
	"cli/tools/virtualenv"
	"cli/utils"
	"fmt"
	"os"
	"path/filepath"
)

type Action string

const (
	UNKNOWN Action = ""
	CREATE  Action = "create"
	SCALE   Action = "scale"
	UPGRADE Action = "upgrade"
)

func Apply(userCfgPath string, action env.ApplyAction) error {
	c, err := NewCluster(userCfgPath)

	if err != nil {
		return err
	}

	var events []*OnChangeEvent

	if c.OldCfg != nil {
		events, err := plan(c, action)

		if err != nil {
			return err
		}

		if len(events) == 0 {
			return fmt.Errorf("No changes detected. Aborting...")
		}
	}

	if err := c.Prepare(userCfgPath); err != nil {
		return err
	}

	switch action {
	case env.CREATE:
		err = c.Create()
	case env.UPGRADE:
		err = c.Upgrade()
	case env.SCALE:
		err = c.Scale(events)
	}

	if err != nil {
		return err
	}

	return applyNewConfig(c.Path)
}

// prepare prepares cluster's directory. It ensures that Kubitect project
// files are present in the directory, new configuration file is stored in
// the temporary location and that virtual environment is created.
func (c Cluster) Prepare(userCfgPath string) error {
	if err := c.Init(); err != nil {
		return err
	}

	if err := storeNewConfig(c.Path, userCfgPath); err != nil {
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

// Init ensures cluster directory exists and that all required files
// are copied from the Kubitect git project. If local flag is used,
// project files are copied from the current directory.
func (c Cluster) Init() error {
	cfg := c.NewCfg

	url := env.ConstProjectUrl
	version := env.ConstProjectVersion

	if cfg.Kubitect.Url != nil {
		url = string(*cfg.Kubitect.Url)
	}

	if cfg.Kubitect.Version != nil {
		version = string(*cfg.Kubitect.Version)
	}

	if env.DebugMode {
		utils.PrintDebug("kubitect.url: %s", url)
		utils.PrintDebug("kubitect.version: %s", version)
	}

	if env.Local {
		wd, err := os.Getwd()

		if err != nil {
			return err
		}

		return copyReqFiles(wd, c.Path)
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

// setupMainVE creates main (Kubitect) virtual environment.
func (c Cluster) SetupMainVE() error {
	ktVer := env.ConstProjectVersion

	if c.NewCfg.Kubitect.Version != nil {
		ktVer = string(*c.NewCfg.Kubitect.Version)
	}

	return virtualenv.Env.Main.Setup(c.Path, ktVer)
}

// setupKubesprayVE creates Kubespray virtual environment.
func (c Cluster) SetupKubesprayVE() error {
	ksVer := string(*c.NewCfg.Kubernetes.Kubespray.Version)
	return virtualenv.Env.Kubespray.Setup(c.Path, ksVer)
}
