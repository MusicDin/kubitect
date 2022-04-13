package cmd

import (
	"cli/config"
	"cli/env"
	"cli/helpers"
	"cli/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	tmpDirName            = "temp"
	mainRequirements      = "requirements.txt"
	mainVenvName          = "main-venv"
	kubesprayRequirements = "ansible/kubespray/requirements.txt"
	kubesprayVenvName     = "kubespray-venv"
)

var (
	validNodeTypes = []string{
		"worker",
		"master",
		"loadBalancer",
	}
)

type Node struct {
	Id   int    `yaml:"id"`
	Name string `yaml:"name"`
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
			fmt.Println(err)
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

	// Add completion values for flag 'action'.
	applyCmd.RegisterFlagCompletionFunc("action", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return env.ProjectApplyActions[:], cobra.ShellCompDirectiveDefault
	})

	// Auto complete cluster names from project clusters directory
	// for flag 'cluster'.
	applyCmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		clustersPath := filepath.Join(env.ProjectHomePath, env.ConstProjectClustersDir)
		return []string{clustersPath}, cobra.ShellCompDirectiveFilterDirs
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

	// Fail if cluster path is not pointing on a valid cluster directory.
	err = utils.VerifyClusterDir(env.ClusterPath)
	if err != nil {
		return err
	}

	fmt.Println("Creating main virtual environment...")

	// Prepare main virtual environment.
	err = helpers.PrepareVirtualEnironment(env.ClusterPath, mainVenvName, mainRequirements)
	if err != nil {
		return err
	}

	extravars := []string{
		"tkk_home=" + env.ProjectHomePath,
		"tkk_cluster_action=" + env.ClusterAction,
		"tkk_cluster_name=" + env.ClusterName,
		"tkk_cluster_path=" + env.ClusterPath,
	}

	if env.IsCustomConfig {
		extravars = append(extravars, "config_path="+env.ConfigPath)
	}

	// Execute the project ansible playbook.
	err = helpers.ExecAnsiblePlaybookLocal(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
		VenvName:     mainVenvName,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/tkk/init.yaml"),
		Extravars:    extravars,
	})
	if err != nil {
		return err
	}

	// Nodes should be gracefully removed from the cluster before actual virtual
	// machines are removed. Node removal is triggered on 'scale' action when
	// infrastructure config exists (tf output).
	_, err = os.Stat(infraConfigPath)
	if env.ClusterAction == "scale" && err == nil {

		// Extract required values from tf output
		sshUser, err := config.GetStrValue(infraConfigPath, "cluster.ssh.user")
		if err != nil {
			return err
		}

		sshPKey, err := config.GetStrValue(infraConfigPath, "cluster.ssh.pkey")
		if err != nil {
			return err
		}

		removedWorkers, err := getRemovedNodes(env.ConfigPath, infraConfigPath, "worker")
		if err != nil {
			return err
		}

		if len(removedWorkers) > 0 {

			var removedWorkerNames []string

			utils.PrintWarning("The following nodes will get removed: ")

			for _, worker := range removedWorkers {
				removedWorkerNames = append(removedWorkerNames, worker.Name)
				fmt.Println("- " + worker.Name)
			}

			// Ask user for permission
			confirm := utils.AskUserConfirmation()
			if !confirm {
				return fmt.Errorf("User aborted.")
			}

			extravars = []string{
				"skip_confirmation=yes",
				"delete_nodes_confirmation=yes",
				"node=" + strings.Join(removedWorkerNames, ","),
			}

			// Remove nodes.
			err = helpers.ExecAnsiblePlaybook(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
				VenvName:     kubesprayVenvName,
				PlaybookFile: filepath.Join(env.ClusterPath, "ansible/kubespray/remove-node.yml"),
				Inventory:    filepath.Join(env.ClusterPath, "config/hosts.ini"),
				Become:       true,
				User:         sshUser,
				PrivateKey:   sshPKey,
				Timeout:      3000,
				Extravars:    extravars,
			})
			if err != nil {
				return err
			}
		}
	}

	// Apply terraform if cluster action equals 'create' or 'scale'.
	if utils.StrArrayContains([]string{"create", "scale"}, env.ClusterAction) {

		err = helpers.TerraformApply(env.ClusterPath)
		if err != nil {
			return err
		}
	}

	extravars = []string{
		"tkk_cluster_path=" + env.ClusterPath,
	}

	// Prepare Kubespray configuration files (all.yaml, k8s_cluster.yaml, ...)
	// and clone Kubespray git project.
	err = helpers.ExecAnsiblePlaybookLocal(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
		VenvName:     mainVenvName,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/tkk/kubespray-setup.yaml"),
		Extravars:    extravars,
	})
	if err != nil {
		return err
	}

	// Prepare Kubespray's virtual environment.
	err = helpers.PrepareVirtualEnironment(env.ClusterPath, kubesprayVenvName, kubesprayRequirements)
	if err != nil {
		return err
	}

	// Extract required values from tf output
	sshUser, err := config.GetStrValue(infraConfigPath, "cluster.ssh.user")
	if err != nil {
		return err
	}

	sshPKey, err := config.GetStrValue(infraConfigPath, "cluster.ssh.pkey")
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

		err = helpers.ExecAnsiblePlaybook(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
			VenvName:     mainVenvName,
			PlaybookFile: filepath.Join(env.ClusterPath, "ansible/haproxy/haproxy.yaml"),
			Inventory:    filepath.Join(env.ClusterPath, "config/hosts.ini"),
			Become:       true,
			User:         sshUser,
			PrivateKey:   sshPKey,
			Timeout:      3000,
		})
		if err != nil {
			return err
		}

		extravars = []string{
			"kube_version=" + k8sVersion,
		}

		err = helpers.ExecAnsiblePlaybook(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
			VenvName:     kubesprayVenvName,
			PlaybookFile: filepath.Join(env.ClusterPath, "ansible/kubespray/cluster.yml"),
			Inventory:    filepath.Join(env.ClusterPath, "config/hosts.ini"),
			Become:       true,
			User:         sshUser,
			PrivateKey:   sshPKey,
			Timeout:      3000,
			Extravars:    extravars,
		})
		if err != nil {
			return err
		}

	case "upgrade":

		extravars = []string{
			"kube_version=" + k8sVersion,
		}

		err = helpers.ExecAnsiblePlaybook(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
			VenvName:     kubesprayVenvName,
			PlaybookFile: filepath.Join(env.ClusterPath, "ansible/kubespray/upgrade-cluster.yml"),
			Inventory:    filepath.Join(env.ClusterPath, "config/hosts.ini"),
			Become:       true,
			User:         sshUser,
			PrivateKey:   sshPKey,
			Timeout:      3000,
			Extravars:    extravars,
		})
		if err != nil {
			return err
		}

	case "scale":

		extravars = []string{
			"kube_version=" + k8sVersion,
		}

		err = helpers.ExecAnsiblePlaybook(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
			VenvName:     kubesprayVenvName,
			PlaybookFile: filepath.Join(env.ClusterPath, "ansible/kubespray/scale.yml"),
			Inventory:    filepath.Join(env.ClusterPath, "config/hosts.ini"),
			Become:       true,
			User:         sshUser,
			PrivateKey:   sshPKey,
			Timeout:      3000,
			Extravars:    extravars,
		})
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("Unknown cluster action: %s", env.ClusterAction)
	}

	extravars = []string{
		"tkk_cluster_path=" + env.ClusterPath,
	}

	// Finalize Kubernetes cluster installation.
	err = helpers.ExecAnsiblePlaybook(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
		VenvName:     mainVenvName,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/tkk/finalize.yaml"),
		Inventory:    filepath.Join(env.ClusterPath, "config/hosts.ini"),
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
		Extravars:    extravars,
	})
	if err != nil {
		return err
	}

	return nil
}

// initCluster makes sure cluster directory exists and that all required
// files are copied from the git project to the cluster directory.
func initCluster(clusterPath string) error {

	var err error

	gitTmpDir := filepath.Join(clusterPath, tmpDirName)

	url, err := config.GetStrValue(env.ConfigPath, "tkk.url")
	if err != nil {
		url = env.DefaultGitProjectUrl
	}

	version, err := config.GetStrValue(env.ConfigPath, "tkk.version")
	if err != nil {
		version = env.DefaultGitProjectVersion
	}

	if env.DebugMode {
		fmt.Println("tkk.url: " + url)
		fmt.Println("tkk.version: " + version)
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

// getRemovedNodes function returns a list of nodes that are present
// in the infrastructure config (config created by terraform) and are not
// present in the cluster config (currently applied config).
func getRemovedNodes(configPath string, infraConfigPath string, nodeType string) ([]Node, error) {

	if !utils.StrArrayContains(validNodeTypes, nodeType) {
		return nil, fmt.Errorf("Invalid node type '%s'. Valid node types are [%s].", nodeType, strings.Join(validNodeTypes, ", "))
	}

	var configNodes []Node
	var infraNodes []Node

	configKey := fmt.Sprintf("cluster.nodes.%s.instances[*]", nodeType)

	// Ignore errors, because it is possible that nodes of provided type
	config.GetValue(configPath, configKey, &configNodes)
	config.GetValue(infraConfigPath, configKey, &infraNodes)

	removedNodes := []Node{}

	for _, infNode := range infraNodes {

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
