package cmd

import (
	"cli/env"
	"cli/utils"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// exportKubeconfigCmd represents the exportKubeconfig command
var exportKubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "Export cluster kubeconfig file",
	Long:  `Command export kubeconfig prints content of the kubeconfig file.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := exportKubeconfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	},
}

func init() {
	exportCmd.AddCommand(exportKubeconfigCmd)

	exportKubeconfigCmd.PersistentFlags().StringVar(&env.ClusterName, "cluster", env.DefaultClusterName, "specify the cluster to be used")
	exportKubeconfigCmd.PersistentFlags().BoolVar(&env.Local, "local", false, "use a current directory as the cluster path")

	// Auto complete cluster names from project clusters directory
	// for the flag 'cluster'.
	exportKubeconfigCmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		clustersPath := filepath.Join(env.ProjectHomePath, env.ConstProjectClustersDir)
		return []string{clustersPath}, cobra.ShellCompDirectiveFilterDirs
	})
}

// exportKubeconfig exports (prints) content of the cluster
// Kubeconfig file.
func exportKubeconfig() error {

	kubeconfigPath := filepath.Join(env.ClusterPath, "config", "admin.conf")

	err := utils.VerifyClusterDir(env.ClusterName)
	if err != nil {
		return fmt.Errorf("Cluster '%s' does not exist: %w", env.ClusterName, err)
	}

	_, err = os.Stat(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("Kubeconfig for cluster '%s' does not exist: %w", env.ClusterName, err)
	}

	kubeconfig, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("Failed reading Kubeconfig file: %w", err)
	}

	fmt.Print(string(kubeconfig))

	return nil
}
