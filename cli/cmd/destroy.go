package cmd

import (
	"cli/env"
	"cli/helpers"
	"cli/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the cluster",
	Long: `Destroys the cluster. If cluster is not specified (using --cluster flag)
the operation is executed on the 'default' cluster.`,

	Run: func(cmd *cobra.Command, args []string) {
		err := destroy()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	destroyCmd.PersistentFlags().StringVar(&env.ClusterName, "cluster", env.DefaultClusterName, "specify the cluster to be used")
	destroyCmd.PersistentFlags().BoolVarP(&env.Local, "local", "l", false, "use a current directory as the cluster path")
	destroyCmd.PersistentFlags().BoolVar(&env.AutoApprove, "auto-approve", false, "automatically approve any user permission requests")

	// Auto complete cluster names of active clusters for flag 'cluster'.
	destroyCmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		clusterNames, err := GetClusters([]ClusterFilter{IsActive})
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return clusterNames, cobra.ShellCompDirectiveNoFileComp
	})
}

// destroy function destroys the cluster.
func destroy() error {

	var err error

	tfStatePath := filepath.Join(env.ClusterPath, env.ConstTerraformStatePath)

	// Skip destruction if terraform state file does not exist.
	if !utils.Exists(tfStatePath) {
		return fmt.Errorf("Cluster '%s' is already destroyed (or not yet initialized).", env.ClusterName)
	}

	// Fail if cluster path is not pointing on a valid cluster directory.
	err = utils.VerifyClusterDir(env.ClusterPath)
	if err != nil {
		return err
	}

	// Ask user for permission.
	confirm := utils.AskUserConfirmation("The '%s' cluster will be destroyed.", env.ClusterName)
	if !confirm {
		return fmt.Errorf("User aborted.")
	}

	fmt.Printf("Destroying '%s' cluster...\n", env.ClusterName)

	// Terraform destroy
	err = helpers.TerraformDestroy(env.ClusterPath)
	if err != nil {
		return err
	}

	// Remove terraform state file
	err = os.Remove(tfStatePath)
	if err != nil {
		return fmt.Errorf("Failed removing cluster's terraform state file: %v", err)
	}

	return nil
}
