package cluster

import (
	"cli/cluster/event"
	"cli/env"
	"cli/file"
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

func ToApplyActionType(a string) (ApplyAction, error) {
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

func (c *Cluster) Apply(a string) error {
	action, err := ToApplyActionType(a)

	if err != nil {
		return err
	}

	if c.AppliedConfig == nil && (action == SCALE || action == UPGRADE) {
		c.Ui().Printf(ui.INFO, "Cannot %s cluster '%s'. It has not been created yet.\n\n", action, c.Name)

		err := c.Ui().Ask("Would you like to create it instead?")

		if err != nil {
			return err
		}

		action = CREATE
	}

	var events event.Events

	if c.AppliedConfig != nil {
		events, err = c.plan(action)

		if err != nil {
			return err
		}

		if len(events) == 0 {
			c.Ui().Println(ui.INFO, "No changes detected.")
			return nil
		}
	}

	if err := c.prepare(); err != nil {
		return err
	}

	switch action {
	case CREATE:
		err = c.create()
	case UPGRADE:
		err = c.create()
	case SCALE:
		err = c.create()
	}

	if err != nil {
		return err
	}

	return c.ApplyNewConfig()
}

// create creates a new cluster or modifies the current
// one if the cluster already exists.
func (c *Cluster) create() error {
	if err := c.NewExecutor().Init(); err != nil {
		return err
	}

	if err := c.Terraform().Apply(); err != nil {
		return err
	}

	if err := c.Sync(); err != nil {
		return err
	}

	exec := c.NewExecutor()

	if err := exec.Init(); err != nil {
		return err
	}

	return exec.Create()
}

// upgrade upgrades an existing cluster.
func (c *Cluster) upgrade() error {
	if err := c.Terraform().Apply(); err != nil {
		return err
	}

	if err := c.Sync(); err != nil {
		return err
	}

	exec := c.NewExecutor()

	if err := exec.Init(); err != nil {
		return err
	}

	return exec.Upgrade()
}

// scale scales an existing cluster.
func (c *Cluster) scale(events event.Events) error {
	exec := c.NewExecutor()

	if err := exec.ScaleDown(events); err != nil {
		return err
	}

	if err := c.Terraform().Apply(); err != nil {
		return err
	}

	if err := c.Sync(); err != nil {
		return err
	}

	return exec.ScaleUp(events)
}

// prepare prepares cluster's directory. It ensures that Kubitect project
// files are present in the directory, new configuration file is stored in
// the temporary location and that main virtual environment is created.
func (c *Cluster) prepare() error {
	var err error

	srcDir := c.WorkingDir()
	dstDir := c.Path

	if c.Local {
		err = copyReqFiles(srcDir, dstDir)
	} else {
		srcDir = filepath.Join(dstDir, "tmp")

		proj := git.GitProject{
			Url:     c.KubitectURL(),
			Version: c.KubitectVersion(),
			Path:    srcDir,
			Ui:      c.Ui(),
		}

		c.Ui().Printf(ui.DEBUG, "kubitect.url: %s\n", proj.Url)
		c.Ui().Printf(ui.DEBUG, "kubitect.version: %s\n", proj.Version)

		err = cloneAndCopyReqFiles(proj, c.Path)
	}

	if err == nil {
		return c.StoreNewConfig()
	}

	e, ok := err.(ui.ErrorBlock)

	if !ok {
		return err
	}

	c.Ui().PrintBlockE(e)

	if srcDir == c.WorkingDir() {
		return fmt.Errorf("current (working) directory is missing some required files\n\nAre you sure you are in the right directory?")
	}

	return fmt.Errorf("cluster directory (%s) is missing some required files", srcDir)
}

// cloneAndCopyReqFiles first clones a project using git and then
// copies project required files from the cloned directory to the
// destination directory.
func cloneAndCopyReqFiles(git git.GitProject, dstDir string) error {
	if err := os.RemoveAll(git.Path); err != nil {
		return err
	}

	if err := git.Clone(); err != nil {
		return err
	}

	if err := copyReqFiles(git.Path, dstDir); err != nil {
		return err
	}

	return os.RemoveAll(git.Path)
}

// copyReqFiles copies project required files from source directory
// to the destination directory.
func copyReqFiles(srcDir, dstDir string) error {
	if err := verifyClusterDir(srcDir); err != nil {
		return err
	}

	for _, path := range env.ProjectRequiredFiles {
		src := filepath.Join(srcDir, path)
		dst := filepath.Join(dstDir, path)

		if err := file.ForceCopy(src, dst); err != nil {
			return err
		}
	}

	return verifyClusterDir(dstDir)
}

// verifyClusterDir verifies if the provided cluster directory
// exists and if it contains all necessary directories.
func verifyClusterDir(clusterPath string) error {
	if !file.Exists(clusterPath) {
		return fmt.Errorf("cluster does not exist on path '%s'", clusterPath)
	}

	var missing []string

	for _, path := range env.ProjectRequiredFiles {
		p := filepath.Join(clusterPath, path)

		if !file.Exists(p) {
			missing = append(missing, path)
		}
	}

	if len(missing) > 0 {
		return NewInvalidClusterDirError(missing)
	}

	return nil
}