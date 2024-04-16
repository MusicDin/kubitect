package ansible

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/MusicDin/kubitect/pkg/ui"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
)

type Playbook struct {
	Inventory  string
	Tags       []string
	User       string
	PrivateKey string
	Become     bool
	Local      bool
	Timeout    int
	ExtraVars  map[string]string

	// Playbook's path.
	Path string

	// Working directory. Defaults to playbook's directory.
	WorkingDir string
}

type (
	Ansible interface {
		Exec(Playbook) error
	}

	ansible struct {
		binPath  string
		cacheDir string
	}
)

func NewAnsible(binDir, cacheDir string) Ansible {
	return &ansible{
		binPath:  path.Join(binDir, "ansible-playbook"),
		cacheDir: cacheDir,
	}
}

// Exec executes the given ansible playbook.
func (a *ansible) Exec(pb Playbook) error {
	if len(pb.Path) < 1 {
		return fmt.Errorf("ansible-playbook: playbook path not set")
	}

	if pb.Local && pb.Inventory == "" {
		pb.Inventory = "localhost,"
	}

	if pb.Inventory == "" {
		return fmt.Errorf("ansible-playbook (%s): inventory not set", pb.Path)
	}

	privilegeEscalationOptions := &options.AnsiblePrivilegeEscalationOptions{
		Become: pb.Become,
	}

	connectionOptions := &options.AnsibleConnectionOptions{
		PrivateKey: pb.PrivateKey,
		Timeout:    pb.Timeout,
		User:       pb.User,
	}

	if pb.Local {
		connectionOptions.Connection = "local"
	}

	playbookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: pb.Inventory,
		Tags:      strings.Join(pb.Tags, ","),
		Forks:     "50",
	}

	for k, v := range pb.ExtraVars {
		playbookOptions.AddExtraVar(k, v)
		ui.Printf(ui.WARN, "%s=%s\n", k, v)
	}

	executor := &execute.DefaultExecute{
		CmdRunDir:   filepath.Dir(pb.Path),
		Write:       ui.Streams().Out().File(),
		WriterError: ui.Streams().Err().File(),
	}

	if pb.WorkingDir != "" {
		executor.CmdRunDir = pb.WorkingDir
	}

	playbook := &playbook.AnsiblePlaybookCmd{
		Binary:                     a.binPath,
		Exec:                       executor,
		Playbooks:                  []string{pb.Path},
		Options:                    playbookOptions,
		ConnectionOptions:          connectionOptions,
		PrivilegeEscalationOptions: privilegeEscalationOptions,
		StdoutCallback:             "yaml",
	}

	if !ui.HasColor() {
		options.AnsibleSetEnv("ANSIBLE_NO_COLOR", "true")
	} else {
		// options.AnsibleForceColor()
		options.AnsibleSetEnv("ANSIBLE_FORCE_COLOR", "true")
	}

	if ui.Debug() {
		options.AnsibleSetEnv("ANSIBLE_VERBOSITY", "2")
	} else {
		options.AnsibleSetEnv("ANSIBLE_LOCALHOST_WARNING", "false")
		options.AnsibleSetEnv("ANSIBLE_INVENTORY_UNPARSED_WARNING", "false")
	}

	options.AnsibleSetEnv("ANSIBLE_CACHE_PLUGIN", "jsonfile")
	options.AnsibleSetEnv("ANSIBLE_CACHE_PLUGIN_CONNECTION", a.cacheDir)
	options.AnsibleSetEnv("ANSIBLE_CACHE_PLUGIN_TIMEOUT", "86400")

	options.AnsibleSetEnv("ANSIBLE_CALLBACKS_ENABLED", "yaml")
	options.AnsibleSetEnv("ANSIBLE_HOST_PATTERN_MISMATCH", "ignore")
	options.AnsibleSetEnv("ANSIBLE_DISPLAY_FAILED_STDERR", "true")
	options.AnsibleSetEnv("ANSIBLE_DISPLAY_SKIPPED_HOSTS", "false")
	options.AnsibleSetEnv("ANSIBLE_STDOUT_CALLBACK", "yaml")
	options.AnsibleSetEnv("ANSIBLE_STDERR_CALLBACK", "yaml")

	err := playbook.Run(context.TODO())
	if err != nil {
		return fmt.Errorf("ansible-playbook (%s): %v", filepath.Base(pb.Path), err)
	}

	return nil
}
