package cluster

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/MusicDin/kubitect/embed"
	"github.com/MusicDin/kubitect/pkg/cluster/event"
	"github.com/MusicDin/kubitect/pkg/env"
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
func (a ApplyAction) rules() []event.Rule {
	switch a {
	case CREATE:
		return event.ModifyRules
	case SCALE:
		return event.ScaleRules
	case UPGRADE:
		return event.UpgradeRules
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
		ui.Printf(ui.INFO, "Cannot %s cluster %q. It has not been created yet.\n\n", action, c.Name)

		err := ui.Ask("Would you like to create it instead?")
		if err != nil {
			return err
		}

		action = CREATE
	}

	events, err := c.plan(action)
	if err != nil {
		return err
	}

	if c.AppliedConfig != nil && len(events) == 0 {
		ui.Println(ui.INFO, "No changes detected.")
		return nil
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

// plan compares an already applied configuration file with the new one, and
// detects events based on the apply action. If cluster has not been
// initialized yet, nil is returned both for an error and events.
func (c *Cluster) plan(action ApplyAction) (event.Events, error) {
	if c.AppliedConfig == nil {
		return nil, nil
	}

	cmpOptions := cmp.Options{
		Tag:                "opt",
		ExtraNameTags:      []string{"yaml"},
		RespectSliceOrder:  false,
		IgnoreEmptyChanges: true,
		PopulateAllNodes:   true,
	}

	// Compare configuration files.
	res, err := cmp.Compare(c.AppliedConfig, c.NewConfig, cmpOptions)
	if err != nil {
		return nil, err
	}

	// Return if there is no changes.
	if !res.HasChanges() {
		return nil, nil
	}

	fmtOptions := cmp.FormatOptions{
		ShowColor:            ui.HasColor(),
		ShowDiffOnly:         true,
		ShowChangeTypePrefix: true,
	}

	fmt.Printf("Following changes have been detected:\n\n")
	fmt.Println(res.ToYaml(fmtOptions))

	// Generate events from detected configuration changes and provided rules.
	events, err := event.GenerateEvents(res.Tree(), action.rules())
	if err != nil {
		return nil, err
	}

	hasError := false
	for _, e := range events {
		if !e.Rule.IsOfType(event.Error) {
			continue
		}

		hasError = true
		if e.Change.Type == cmp.Create || e.Change.Type == cmp.Delete {
			// For create and delete events, only change's
			// path is shown.
			err := NewConfigChangeError(e.Rule.Message, e.Change.Path)
			ui.PrintBlockE(err)
		} else {
			err := NewConfigChangeError(e.Rule.Message, e.MatchedChangePaths...)
			ui.PrintBlockE(err)
		}
	}

	if hasError {
		return nil, fmt.Errorf("Configuration file contains errors.")
	}

	hasWarnings := false
	for _, e := range events {
		if !e.Rule.IsOfType(event.Warn) {
			continue
		}

		hasWarnings = true
		err := NewConfigChangeWarning(e.Rule.Message, e.Change.Path)
		ui.PrintBlockE(err)
	}

	if hasWarnings {
		ui.Println(ui.INFO, "Above warnings indicate potentially dangerous actions.")
	}

	return events, ui.Ask()
}

// create creates a new cluster or modifies the current
// one if the cluster already exists.
func (c *Cluster) create() error {
	if err := c.generateSshKeys(); err != nil {
		return err
	}

	if err := c.Provisioner().Init(nil); err != nil {
		return err
	}

	if err := c.Provisioner().Apply(); err != nil {
		return err
	}

	if err := c.Sync(); err != nil {
		return err
	}

	if err := c.Manager().Init(); err != nil {
		return err
	}

	if err := c.Manager().Sync(); err != nil {
		return err
	}

	return c.Manager().Create()
}

// upgrade upgrades an existing cluster.
func (c *Cluster) upgrade() error {
	if err := c.Provisioner().Init(nil); err != nil {
		return err
	}

	if err := c.Provisioner().Apply(); err != nil {
		return err
	}

	if err := c.Sync(); err != nil {
		return err
	}

	if err := c.Manager().Init(); err != nil {
		return err
	}

	if err := c.Manager().Sync(); err != nil {
		return err
	}

	return c.Manager().Upgrade()
}

// scale scales an existing cluster.
func (c *Cluster) scale(events []event.Event) error {
	if err := c.Manager().Init(); err != nil {
		return err
	}

	if err := c.Manager().ScaleDown(events); err != nil {
		return err
	}

	if err := c.Provisioner().Init(events); err != nil {
		return err
	}

	if err := c.Provisioner().Apply(); err != nil {
		return err
	}

	if err := c.Sync(); err != nil {
		return err
	}

	if err := c.Manager().Sync(); err != nil {
		return err
	}

	return c.Manager().ScaleUp(events)
}

// prepare prepares the cluster directory. It ensures all required project
// files are present in the directory and new configuration file is stored in
// the temporary location.
func (c *Cluster) prepare() error {
	if err := os.MkdirAll(c.Path, os.ModePerm); err != nil {
		return fmt.Errorf("create cluster directory: %v", err)
	}

	for _, rf := range env.ProjectRequiredFiles {
		err := embed.MirrorResource(rf, c.Path)
		if err != nil {
			return err
		}
	}

	if err := verifyClusterDir(c.Path); err != nil {
		eb, ok := err.(ui.ErrorBlock)
		if !ok {
			return err
		}

		ui.PrintBlockE(eb)
		return fmt.Errorf("cluster %s is missing some required files", c.Name)
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
	// - Check if the user has provided a custom path to the key pair
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
