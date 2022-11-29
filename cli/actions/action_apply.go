package actions

import (
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

		err := ui.GlobalUi().Ask("Would you like to create it instead?")

		if err != nil {
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
	var err error

	if c.Local {
		err = copyReqFiles(c.Ctx.WorkingDir(), c.Path)
	} else {
		err = cloneAndCopyReqFiles(c)
	}

	if err != nil {
		return err
	}

	return c.StoreNewConfig()
}

func cloneAndCopyReqFiles(c *Cluster) error {
	url := c.KubitectURL()
	version := c.KubitectVersion()

	ui.GlobalUi().Printf(ui.DEBUG, "kubitect.url: %s\n", url)
	ui.GlobalUi().Printf(ui.DEBUG, "kubitect.version: %s\n", version)

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
