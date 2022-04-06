package cmd

import (
	"cli/config"
	"cli/env"
	"cli/helpers"
	"cli/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	tmpDirName          = "temp"
	ansiblePlaybookTags = "apply"
)

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

	fmt.Printf("Preparing cluster '%s'...\n", env.ClusterName)

	if !env.Local {
		err = initCluster(env.ClusterPath)
		if err != nil {
			return err
		}
	}

	// Activate virtual environment and install Ansible.
	err = helpers.PrepareVirtualEnironment(env.ClusterPath)
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
		PlaybookFile: filepath.Join(env.ClusterPath, env.ConstAnsiblePlaybookPath),
		Tags:         ansiblePlaybookTags,
		Extravars:    extravars,
	})
	if err != nil {
		return err
	}

	// Terraform apply
	err = helpers.TerraformApply(env.ClusterPath)
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

	// Replace relevant files from temporary git project.
	for _, path := range env.ProjectRequiredFiles {

		srcPath := filepath.Join(gitTmpDir, path)
		dstPath := filepath.Join(clusterPath, path)

		utils.ForceMove(srcPath, dstPath)
	}

	// Remove temp project.
	err = os.RemoveAll(gitTmpDir)

	if err != nil {
		return fmt.Errorf("Failed removing temporary git project: %w", err)
	}

	return nil
}
