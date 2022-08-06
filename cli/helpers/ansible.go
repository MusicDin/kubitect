package helpers

import (
	"cli/env"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
)

const (
	extraVarsSplitToken = "="
)

var (
	VenvNameMissing           = errors.New("VenvName must be provided!")
	ClusterPathMissing        = errors.New("Cluster path is missing.")
	AnsiblePlaybookCmdMissing = errors.New("AnsiblePlaybookCmd is null.")
	PlaybookFilePathMissing   = errors.New("To run ansible-playbook playbook file path must be specified.")
	InventoryMissing          = errors.New("To run ansible-playbook an inventory must be specified.")
)

type AnsiblePlaybookCmd struct {
	PlaybookFile    string // "clusterPath/PlaybookFile"
	Inventory       string
	Tags            string
	User            string
	PrivateKey      string
	Become          bool
	ConnectionLocal bool
	Timeout         int
	Extravars       []string
	Venv            VirtualEnvironment
}

// ExecAnsiblePlaybookLocal sets inventory and connection type to localhost
// before executing the ansible playbook.
func ExecAnsiblePlaybookLocal(clusterPath string, ansibleCmd *AnsiblePlaybookCmd) error {

	if ansibleCmd == nil {
		return AnsiblePlaybookCmdMissing
	}

	ansibleCmd.Inventory = "localhost,"
	ansibleCmd.ConnectionLocal = true

	return ExecAnsiblePlaybook(clusterPath, ansibleCmd)
}

// ExecAnsibleCmd executes ansible playbook with working directory
// set to the cluster path directory.
func ExecAnsiblePlaybook(clusterPath string, ansibleCmd *AnsiblePlaybookCmd) error {

	if ansibleCmd == nil {
		return AnsiblePlaybookCmdMissing
	}

	if len(ansibleCmd.PlaybookFile) < 1 {
		return PlaybookFilePathMissing
	}

	if len(ansibleCmd.Inventory) < 1 {
		return InventoryMissing
	}

	if len(ansibleCmd.Venv.Name) < 1 {
		return VenvNameMissing
	}

	if len(clusterPath) < 1 {
		return ClusterPathMissing
	}

	privilegeEscalationOptions := &options.AnsiblePrivilegeEscalationOptions{
		Become: ansibleCmd.Become,
	}

	connectionOptions := &options.AnsibleConnectionOptions{
		PrivateKey: ansibleCmd.PrivateKey,
		Timeout:    ansibleCmd.Timeout,
		User:       ansibleCmd.User,
	}

	if ansibleCmd.ConnectionLocal {
		connectionOptions.Connection = "local"
	}

	playbookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: ansibleCmd.Inventory,
		Tags:      ansibleCmd.Tags,
	}

	if env.DebugMode {
		playbookOptions.Verbose = true
	}

	vars, err := extravarsListToMap(ansibleCmd.Extravars)
	if err != nil {
		return err
	}

	for keyVar, valueVar := range vars {
		playbookOptions.AddExtraVar(keyVar, valueVar)
	}

	executor := &execute.DefaultExecute{
		CmdRunDir: filepath.Dir(ansibleCmd.PlaybookFile),
	}

	playbook := &playbook.AnsiblePlaybookCmd{
		Binary:                     filepath.Join(clusterPath, env.ConstVenvBinDir, ansibleCmd.Venv.Name, "bin", "ansible-playbook"),
		Exec:                       executor,
		Playbooks:                  []string{ansibleCmd.PlaybookFile},
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
		return fmt.Errorf("Error while running ansible-playbook: %w", err)
	}

	return nil
}

// extravarsListToMap converts array of "key=value" strings into map.
func extravarsListToMap(extravarsList []string) (map[string]interface{}, error) {

	extravarsMap := map[string]interface{}{}

	for _, extravar := range extravarsList {

		tokens := strings.Split(extravar, extraVarsSplitToken)
		if len(tokens) != 2 {
			return nil, fmt.Errorf("Invalid extravar format for '%s'. The format should be 'key=value'.", extravar)
		}

		extravarsMap[tokens[0]] = tokens[1]
	}

	return extravarsMap, nil
}
