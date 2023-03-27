package cluster

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/MusicDin/kubitect/pkg/cluster/event"
	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/MusicDin/kubitect/pkg/tools/git"
	"github.com/MusicDin/kubitect/pkg/ui"
	"github.com/MusicDin/kubitect/pkg/utils/cmp"
	"github.com/MusicDin/kubitect/pkg/utils/file"
	"github.com/MusicDin/kubitect/pkg/utils/keygen"
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

// events returns events of the corresponding action.
func (a ApplyAction) events() event.Events {
	switch a {
	case CREATE:
		return event.ModifyEvents
	case SCALE:
		return event.ScaleEvents
	case UPGRADE:
		return event.UpgradeEvents
	default:
		return nil
	}
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

// Apply either creates new or modifies an existing cluster, based on the
// provided action.
func (c *Cluster) Apply(a string) error {
	action, err := ToApplyActionType(a)

	if err != nil {
		return err
	}

	if c.AppliedConfig == nil && (action == SCALE || action == UPGRADE) {
		ui.Printf(ui.INFO, "Cannot %s cluster '%s'. It has not been created yet.\n\n", action, c.Name)

		err := ui.Ask("Would you like to create it instead?")
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
			ui.Println(ui.INFO, "No changes detected.")
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
		err = c.upgrade()
	case SCALE:
		err = c.scale(events)
	}

	if err != nil {
		return err
	}

	return c.ApplyNewConfig()
}

// plan compares new and applied configuration files, and detects
// events based on the apply action.
//
// If applied configuration file does not exist, no events and no
// error is returned.
// If blocking changes are detected, an error is returned.
// If warnings are detected, user is asked for permission to continue.
func (c *Cluster) plan(action ApplyAction) (event.Events, error) {
	if c.AppliedConfig == nil {
		return nil, nil
	}

	comp := cmp.NewComparator()
	comp.Tag = "opt"
	comp.ExtraNameTags = []string{"yaml"}
	comp.IgnoreEmptyChanges = true
	comp.PopulateStructNodes = true

	diff, err := comp.Compare(c.AppliedConfig, c.NewConfig)

	if err != nil || len(diff.Changes()) == 0 {
		return nil, err
	}

	fmt.Printf("Following changes have been detected:\n\n")
	fmt.Println(diff.ToYamlDiff())

	events := event.TriggerEvents(diff, action.events())
	blocking := events.Blocking()

	if len(blocking) > 0 {
		ui.PrintBlockE(blocking.Errors()...)
		return nil, fmt.Errorf("Aborted. Configuration file contains errors.")
	}

	warnings := events.Warns()

	if len(warnings) > 0 {
		ui.PrintBlockE(warnings.Errors()...)
		fmt.Println("Above warnings indicate potentially destructive actions.")
	}

	return events, ui.Ask()
}

// create creates a new cluster or modifies the current
// one if the cluster already exists.
func (c *Cluster) create() error {
	if err := c.generateSshKeys(); err != nil {
		return err
	}

	if err := c.Provisioner().Init(); err != nil {
		return err
	}

	if err := c.Provisioner().Apply(); err != nil {
		return err
	}

	if err := c.Sync(); err != nil {
		return err
	}

	if err := c.Executor().Init(); err != nil {
		return err
	}

	return c.Executor().Create()
}

// upgrade upgrades an existing cluster.
func (c *Cluster) upgrade() error {
	if err := c.Provisioner().Init(); err != nil {
		return err
	}

	if err := c.Provisioner().Apply(); err != nil {
		return err
	}

	if err := c.Sync(); err != nil {
		return err
	}

	if err := c.Executor().Init(); err != nil {
		return err
	}

	return c.Executor().Upgrade()
}

// scale scales an existing cluster.
func (c *Cluster) scale(events event.Events) error {
	if err := c.Executor().Init(); err != nil {
		return err
	}

	if err := c.Executor().ScaleDown(events); err != nil {
		return err
	}

	if err := c.Provisioner().Init(); err != nil {
		return err
	}

	if err := c.Provisioner().Apply(); err != nil {
		return err
	}

	if err := c.Sync(); err != nil {
		return err
	}

	return c.Executor().ScaleUp(events)
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
		tmpDir := filepath.Join(dstDir, "tmp")
		proj := git.NewGitProject(env.ConstProjectUrl, env.ConstProjectVersion)

		ui.Printf(ui.DEBUG, "kubitect.url: %s\n", proj.Url())
		ui.Printf(ui.DEBUG, "kubitect.version: %s\n", proj.Version())

		err = cloneAndCopyReqFiles(proj, tmpDir, c.Path)
	}

	if err != nil {
		e, ok := err.(ui.ErrorBlock)
		if !ok {
			return err
		}

		ui.PrintBlockE(e)
		return fmt.Errorf("cluster directory (%s) is missing some required files", srcDir)
	}

	return c.StoreNewConfig()
}

// generateSshKeys ensures that SSH keys for a cluster a exist in the
// cluster directory.
//
// If the key pair is missing, the keys are either generated or retrieved
// from the location specified by the user in a node template section of
// the configuration file. However, if SSH keys already exist in the
// cluster directory, no action is taken.
func (c *Cluster) generateSshKeys() error {
	ui.Println(ui.INFO, "Ensuring SSH keys are present...")

	kpName := path.Base(c.PrivateSshKeyPath())
	kpDir := path.Dir(c.PrivateSshKeyPath())

	// Stop if key pair already exists
	if keygen.KeyPairExists(kpDir, kpName) {
		return nil
	}

	// Keypair does not exist.
	// - Check if user provided path to the custom key pair
	pkPath := string(c.NewConfig.Cluster.NodeTemplate.SSH.PrivateKeyPath)
	if pkPath != "" {
		kp, err := keygen.ReadKeyPair(path.Dir(pkPath), path.Base(pkPath))
		if err != nil {
			return err
		}

		return kp.Write(kpDir, kpName)
	}

	// Keypair does not exist and user has not provided a custom path.
	// - Generate new key pair
	kp, err := keygen.NewKeyPair(4096)
	if err != nil {
		return err
	}

	return kp.Write(kpDir, kpName)

}

// cloneAndCopyReqFiles first clones a project using git and then
// copies project required files from the cloned directory to the
// destination directory.
func cloneAndCopyReqFiles(proj git.GitProject, tmpDir, dstDir string) error {
	if err := os.RemoveAll(tmpDir); err != nil {
		return err
	}

	if err := proj.Clone(tmpDir); err != nil {
		return err
	}

	if err := copyReqFiles(tmpDir, dstDir); err != nil {
		return err
	}

	return os.RemoveAll(tmpDir)
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
