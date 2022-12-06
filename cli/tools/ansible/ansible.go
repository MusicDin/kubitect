package ansible

import (
	"cli/ui"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
)

type Playbook struct {
	PlaybookFile string
	Inventory    string
	Tags         []string
	User         string
	PrivateKey   string
	Become       bool
	Local        bool
	Timeout      int
	ExtraVars    []string
}

type Ansible struct {
	BinPath    string
	WorkingDir string

	Ui *ui.Ui
}

// Exec executes given ansible playbook.
func (a *Ansible) Exec(pb Playbook) error {
	if pb.Local {
		pb.Inventory = "localhost,"
	}

	if len(pb.PlaybookFile) < 1 {
		return fmt.Errorf("ansible-playbook: file path not set")
	}

	if pb.Inventory == "" {
		return fmt.Errorf("ansible-playbook (%s): inventory not set", pb.PlaybookFile)
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

	if a.Ui.Debug {
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
		CmdRunDir:   filepath.Dir(pb.PlaybookFile),
		Write:       a.Ui.Streams.Out.File,
		WriterError: a.Ui.Streams.Err.File,
	}

	playbook := &playbook.AnsiblePlaybookCmd{
		Binary:                     a.BinPath,
		Exec:                       executor,
		Playbooks:                  []string{pb.PlaybookFile},
		Options:                    playbookOptions,
		ConnectionOptions:          connectionOptions,
		PrivilegeEscalationOptions: privilegeEscalationOptions,
		StdoutCallback:             "yaml",
	}

	// options.AnsibleSetEnv("ANSIBLE_NO_COLOR", "true")    // disable color
	// options.AnsibleSetEnv("ANSIBLE_FORCE_COLOR", "true") // force color

	options.AnsibleForceColor()
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
		pb := filepath.Base(pb.PlaybookFile)
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
