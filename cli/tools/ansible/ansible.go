package ansible

import (
	"cli/env"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
)

type Playbook struct {
	PlaybookFile string // "clusterPath/PlaybookFile"
	Inventory    string
	Tags         []string
	User         string
	PrivateKey   string
	Become       bool
	Local        bool
	Timeout      int
	ExtraVars    []string
	VenvPath     string
}

func (pb Playbook) Exec() error {
	if pb.Local {
		pb.Inventory = "localhost,"
	}

	return pb.exec()
}

// exec executes ansible playbook with working directory
// set to the cluster path directory.
func (pb Playbook) exec() error {
	if len(pb.PlaybookFile) < 1 {
		return fmt.Errorf("ansible-playbook: file path not set")
	}

	if len(pb.Inventory) < 1 {
		return fmt.Errorf("ansible-playbook (%s): inventory not set", pb.PlaybookFile)
	}

	if len(pb.VenvPath) < 1 {
		return fmt.Errorf("ansible-playbook (%s): virtual environment path not set", pb.PlaybookFile)
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

	if env.Debug {
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
		CmdRunDir: filepath.Dir(pb.PlaybookFile),
	}

	playbook := &playbook.AnsiblePlaybookCmd{
		Binary:                     filepath.Join(pb.VenvPath, "bin", "ansible-playbook"),
		Exec:                       executor,
		Playbooks:                  []string{pb.PlaybookFile},
		Options:                    playbookOptions,
		ConnectionOptions:          connectionOptions,
		PrivilegeEscalationOptions: privilegeEscalationOptions,
		StdoutCallback:             "yaml",
	}

	options.AnsibleForceColor()
	options.AnsibleSetEnv("ANSIBLE_DISPLAY_FAILED_STDERR", "true")
	options.AnsibleSetEnv("ANSIBLE_DISPLAY_SKIPPED_HOSTS", "false")

	err = playbook.Run(context.TODO())
	if err != nil {
		pb := filepath.Base(pb.PlaybookFile)
		return fmt.Errorf("ansible-playbook (%s): %v", pb, err)
	}

	return nil
}

// extraVarsToMap converts slice of "key=value" strings into map.
func extraVarsToMap(extraVars []string) (map[string]string, error) {
	var evMap map[string]string

	for _, v := range extraVars {
		tokens := strings.Split(v, "=")

		if len(tokens) != 2 {
			return nil, fmt.Errorf("extraVarsToMap: variable (%s) must be in 'key=value' format", v)
		}

		evMap[tokens[0]] = tokens[1]
	}

	return evMap, nil
}
