package ansible

import (
	"context"
	"fmt"
	"github.com/MusicDin/kubitect/pkg/ui"
	"path"
	"path/filepath"
	"strings"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
)

type Playbook struct {
	Path       string
	Inventory  string
	Tags       []string
	User       string
	PrivateKey string
	Become     bool
	Local      bool
	Timeout    int
	ExtraVars  []string
}

type (
	Ansible interface {
		Exec(Playbook) error
	}

	ansible struct {
		binPath string
	}
)

func NewAnsible(binDir string) Ansible {
	return &ansible{
		binPath: path.Join(binDir, "ansible-playbook"),
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
	}

	if ui.Debug() {
		playbookOptions.Verbose = true
	}

	vars, err := extraVarsToMap(pb.ExtraVars)
	if err != nil {
		return err
	}

	for keyVar, valueVar := range vars {
		playbookOptions.AddExtraVar(keyVar, valueVar)
	}

	executor := &execute.DefaultExecute{
		CmdRunDir:   filepath.Dir(pb.Path),
		Write:       ui.Streams().Out().File(),
		WriterError: ui.Streams().Err().File(),
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

	options.AnsibleSetEnv("ANSIBLE_CALLBACKS_ENABLED", "yaml")
	options.AnsibleSetEnv("ANSIBLE_HOST_PATTERN_MISMATCH", "ignore")
	options.AnsibleSetEnv("ANSIBLE_DISPLAY_FAILED_STDERR", "true")
	options.AnsibleSetEnv("ANSIBLE_DISPLAY_SKIPPED_HOSTS", "false")
	options.AnsibleSetEnv("ANSIBLE_DISPLAY_ARGS_TO_STDOUT", "false")
	options.AnsibleSetEnv("ANSIBLE_FORKS", "10")
	options.AnsibleSetEnv("ANSIBLE_STDOUT_CALLBACK", "yaml")
	options.AnsibleSetEnv("ANSIBLE_STDERR_CALLBACK", "yaml")

	err = playbook.Run(context.TODO())

	if err != nil {
		pb := filepath.Base(pb.Path)
		return fmt.Errorf("ansible-playbook (%s): %v", pb, err)
	}

	return nil
}

// extraVarsToMap converts slice of "key=value" strings into map.
func extraVarsToMap(extraVars []string) (map[string]string, error) {
	evMap := make(map[string]string)

	for _, v := range extraVars {
		tokens := strings.Split(v, "=")

		if len(tokens) != 2 {
			return nil, fmt.Errorf("extraVarsToMap: variable (%s) must be in 'key=value' format", v)
		}

		evMap[tokens[0]] = tokens[1]
	}

	return evMap, nil
}
