package helpers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
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
	ClusterPathMissing              = errors.New("Cluster path is missing.")
	AnsiblePlaybookCmdMissing       = errors.New("AnsiblePlaybookCmd is null.")
	AnsiblePlaybookFilePathMissing  = errors.New("To run ansible-playbook playbook file path must be specified.")
	AnsiblePlaybookInventoryMissing = errors.New("To run ansible-playbook an inventory must be specified.")
)

type AnsiblePlaybookCmd struct {
	PlaybookFile    string // "clusterPath/PlaybookFile"
	Tags            string
	Inventory       string
	ConnectionLocal bool
	Extravars       []string
}

// Sets inventory and connection type to localhost before executing
// ansible playbook.
func ExecAnsiblePlaybookLocal(clusterPath string, ansibleCmd *AnsiblePlaybookCmd) error {

	if ansibleCmd == nil {
		return AnsiblePlaybookCmdMissing
	}

	ansibleCmd.Inventory = "127.0.0.1,"
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
		return AnsiblePlaybookFilePathMissing
	}

	if len(ansibleCmd.Inventory) < 1 {
		return AnsiblePlaybookInventoryMissing
	}

	if len(clusterPath) < 1 {
		return ClusterPathMissing
	}

	vars, err := extravarsListToMap(ansibleCmd.Extravars)
	if err != nil {
		return err
	}

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{}
	if ansibleCmd.ConnectionLocal {
		ansiblePlaybookConnectionOptions.Connection = "local"
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: ansibleCmd.Inventory,
		Tags:      ansibleCmd.Tags,
	}

	for keyVar, valueVar := range vars {
		ansiblePlaybookOptions.AddExtraVar(keyVar, valueVar)
	}

	executor := execute.NewDefaultExecute(
		execute.WithWriteError(io.Writer(os.Stdout)),
		// execute.WithTransformers(
		// 	results.Prepend("[ - ]"),
		// ),
	)

	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:         []string{ansibleCmd.PlaybookFile},
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
		StdoutCallback:    "yaml",
		Binary:            filepath.Join(clusterPath, venvName, "bin", "ansible-playbook"),
		Exec:              executor,
	}

	options.AnsibleForceColor()

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
