package cmd

import (
	"cli/config"
	"cli/env"
	"cli/helpers"
	"cli/playbook"
	"cli/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	tmpDirName = "temp"
)

var (
	validNodeTypes = []string{
		"worker",
		"master",
		"loadBalancer",
	}
)

type Node struct {
	Id        int    `yaml:"id"`
	Ip        string `yaml:"ip"`
	Name      string `yaml:"name"`
	IsRemoved bool   `yaml:"removed"` // Tag nodes as removed after removal
}

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply configuration and create the cluster",
	Long: `Apply command generates Terraform main.tf file based on the provided
configuration. Generated configuration is then passed to the Terraform 
to install the cluster.This way multiple hosts can be used to deploy a 
single cluster.`,

	Run: func(cmd *cobra.Command, args []string) {
		err := apply()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.PersistentFlags().StringVarP(&env.ConfigPath, "config", "c", "", "specify path to the cluster config file")
	applyCmd.PersistentFlags().StringVarP(&env.ClusterAction, "action", "a", env.DefaultClusterAction, "specify cluster action")
	applyCmd.PersistentFlags().StringVar(&env.ClusterName, "cluster", env.DefaultClusterName, "specify the cluster to be used")
	applyCmd.PersistentFlags().BoolVarP(&env.Local, "local", "l", false, "use a current directory as the cluster path")
	applyCmd.PersistentFlags().BoolVar(&env.AutoApprove, "auto-approve", false, "automatically approve any user permission requests")

	// Add completion values for flag 'action'.
	applyCmd.RegisterFlagCompletionFunc("action", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return env.ProjectApplyActions[:], cobra.ShellCompDirectiveDefault
	})

	// Auto complete cluster names for flag 'cluster'.
	applyCmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		clusterNames, err := GetClusters([]ClusterFilter{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return clusterNames, cobra.ShellCompDirectiveNoFileComp
	})
}

// apply function triggers initialization and installation of the cluster.
func apply() error {

	var err error

	infraConfigPath := filepath.Join(env.ClusterPath, "config/infrastructure.yaml")

	fmt.Printf("Preparing cluster '%s'...\n", env.ClusterName)

	if !env.Local {
		err = initCluster(env.ClusterPath)
		if err != nil {
			return err
		}
	}

	// Fail if the cluster path is pointing on an invalid cluster directory.
	err = utils.VerifyClusterDir(env.ClusterPath)
	if err != nil {
		return err
	}

	// Prepare main virtual environment.
	err = helpers.SetupVirtualEnironment(env.ClusterPath, helpers.Venvs.Main)
	if err != nil {
		return err
	}

	// Execute ansible playbook that verifies configuration file and generates
	// hosts inventory file.
	err = playbook.KubitectInit()
	if err != nil {
		return err
	}

	// Ensure target hosts meet all the requirements.
	err = playbook.KubitectHostsSetup()
	if err != nil {
		return err
	}

	// Remove nodes (if any nodes are removed from the configuration).
	err = removeNodes(env.ConfigPath, infraConfigPath, "worker")
	if err != nil {
		return err
	}

	// Apply terraform if the cluster action equals to 'create' or 'scale'.
	if utils.StrArrayContains([]string{"create", "scale"}, env.ClusterAction) {

		err = helpers.TerraformApply(env.ClusterPath)
		if err != nil {
			return err
		}
	}

	// Prepare Kubespray configuration files.
	err = playbook.KubitectKubespraySetup()
	if err != nil {
		return err
	}

	// Prepare Kubespray's virtual environment.
	err = helpers.SetupVirtualEnironment(env.ClusterPath, helpers.Venvs.Kubespray)
	if err != nil {
		return err
	}

	sshUser, err := config.GetStrValue(infraConfigPath, "cluster.nodeTemplate.user")
	if err != nil {
		return err
	}

	sshPKey, err := config.GetStrValue(infraConfigPath, "cluster.nodeTemplate.ssh.privateKeyPath")
	if err != nil {
		return err
	}

	k8sVersion, err := config.GetStrValue(env.ConfigPath, "kubernetes.version")
	if err != nil {
		return err
	}

	// Run Kubespray role based on the provided cluster action.
	switch env.ClusterAction {
	case "create":

		err = playbook.HAProxySetup(sshUser, sshPKey)
		if err != nil {
			return err
		}

		err = playbook.KubesprayCreate(sshUser, sshPKey, k8sVersion)
		if err != nil {
			return err
		}

	case "upgrade":

		err = playbook.KubesprayUpgrade(sshUser, sshPKey, k8sVersion)
		if err != nil {
			return err
		}

	case "scale":

		err = playbook.HAProxySetup(sshUser, sshPKey)
		if err != nil {
			return err
		}

		playbook.KubesprayScale(sshUser, sshPKey, k8sVersion)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("Unknown cluster action: %s", env.ClusterAction)
	}

	// Finalize Kubernets cluster installation.
	playbook.KubitectFinalize(sshUser, sshPKey)

	return nil
}

// initCluster makes sure cluster directory exists and that all required
// files are copied from the git project to the cluster directory.
func initCluster(clusterPath string) error {

	var err error

	gitTmpDir := filepath.Join(clusterPath, tmpDirName)

	url, err := config.GetStrValue(env.ConfigPath, "kubitect.url")
	if err != nil {
		url = env.ConstProjectUrl
	}

	version, err := config.GetStrValue(env.ConfigPath, "kubitect.version")
	if err != nil {
		version = env.ConstProjectVersion
	}

	if env.DebugMode {
		utils.PrintDebug("kubitect.url: %s", url)
		utils.PrintDebug("kubitect.version: %s", version)
	}

	// Make sure that the cluster folder exists
	os.MkdirAll(clusterPath, os.ModePerm)

	// Remove git project temporary directory if it exists.
	os.RemoveAll(gitTmpDir)

	// Clone git project into the temporary directory.
	err = helpers.GitClone(gitTmpDir, url, version)
	if err != nil {
		return err
	}

	// Move relevant project files into cluster directory.
	for _, path := range env.ProjectRequiredFiles {

		srcPath := filepath.Join(gitTmpDir, path)
		dstPath := filepath.Join(clusterPath, path)

		utils.ForceMove(srcPath, dstPath)
	}

	// Remove temporary directory.
	err = os.RemoveAll(gitTmpDir)
	if err != nil {
		return fmt.Errorf("Failed removing temporary git project: %w", err)
	}

	return nil
}

// removeNodes function identifies which nodes are deleted from the cluster config
// and gracefully removes them from the Kubernetes cluster. Removed nodes are then
// marked as removed in the infrastructure config to prevent them from being removed
// again. Node removal is triggered only when cluster action is set to scale and
// when infrastructure config already exists.
func removeNodes(configPath string, infraConfigPath string, nodeType string) error {

	// Verify that nodeType is valid.
	if !utils.StrArrayContains(validNodeTypes, nodeType) {
		return fmt.Errorf("Invalid node type '%s'. Valid node types are [%s].", nodeType, strings.Join(validNodeTypes, ", "))
	}

	// Check if infrastructure config exists.
	_, err := os.Stat(infraConfigPath)

	// Trigger removal if the cluster action is set to 'scale' and the infrastrcutre
	// config already exists.
	if env.ClusterAction == "scale" && err == nil {

		// Extract required values from tf output.
		sshUser, err := config.GetStrValue(infraConfigPath, "cluster.nodeTemplate.user")
		if err != nil {
			return err
		}

		sshPKey, err := config.GetStrValue(infraConfigPath, "cluster.nodeTemplate.ssh.privateKeyPath")
		if err != nil {
			return err
		}

		// Get list of all nodes.
		nodes, err := getNodes(infraConfigPath, nodeType)
		if err != nil {
			return err
		}

		// Get list of removed nodes.
		removedNodes, err := getRemovedNodes(configPath, infraConfigPath, nodeType)
		if err != nil {
			return err
		}

		if len(removedNodes) > 0 {

			var removedNodeNames []string

			warning := "The following nodes will get removed:\n"

			for _, node := range removedNodes {
				removedNodeNames = append(removedNodeNames, node.Name)
				warning = warning + "- " + node.Name + "\n"
			}

			// Ask user for permission.
			confirm := utils.AskUserConfirmation(warning)
			if !confirm {
				return fmt.Errorf("User aborted.")
			}

			// Remove Kubespray nodes
			err = playbook.KubesprayRemoveNodes(sshUser, sshPKey, removedNodeNames)
			if err != nil {
				return err
			}

			// Tag nodes in infrastructure config as removed. This prevents
			// nodes from being removed again on the next run if Terraform
			// fails for some reason on the first run.
			saveTaggedNodes(infraConfigPath, nodes, removedNodes, nodeType)
		}
	}

	return nil
}

// getNodes returns all a list of nodes from provided config file. Only
// nodes that match provided node type are returned.
func getNodes(configPath string, nodeType string) ([]Node, error) {

	var nodes []Node

	configKey := fmt.Sprintf("cluster.nodes.%s.instances[*]", nodeType)

	// Ignore errors, because it is possible that nodes of provided type
	// does not exist.
	config.GetValue(configPath, configKey, &nodes)

	return nodes, nil
}

// getRemovedNodes function returns a list of nodes that are present in
// the infrastructure config (config created by terraform) and are not
// present in the cluster config (currently applied config).
func getRemovedNodes(configPath string, infraConfigPath string, nodeType string) ([]Node, error) {

	configNodes, err := getNodes(configPath, nodeType)
	if err != nil {
		return nil, err
	}

	infraNodes, err := getNodes(infraConfigPath, nodeType)
	if err != nil {
		return nil, err
	}

	removedNodes := []Node{}

	for _, infNode := range infraNodes {

		// Skip already removed nodes
		if infNode.IsRemoved {
			continue
		}

		isNodeRemoved := true

		for _, cfgNode := range configNodes {
			if cfgNode.Id == infNode.Id {
				isNodeRemoved = false
				break
			}
		}

		if isNodeRemoved {
			removedNodes = append(removedNodes, infNode)
		}
	}

	return removedNodes, nil
}

// saveTaggedNodes function taggs removed nodes and saves them into provided
// config file.
func saveTaggedNodes(configPath string, nodes []Node, removedNodes []Node, nodeType string) error {

	configKey := fmt.Sprintf("cluster.nodes.%s.instances", nodeType)

	// Tag removed nodes.
	for i := range nodes {
		for _, removedNode := range removedNodes {
			if nodes[i].Id == removedNode.Id {
				nodes[i].IsRemoved = true
			}
		}
	}

	// Replace existing nodes with modified nodes in config.
	yaml, err := config.ReplaceValue(configPath, configKey, nodes)
	if err != nil {
		return err
	}

	fmt.Println(yaml)

	// Save config.
	err = os.WriteFile(configPath, []byte(yaml), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
