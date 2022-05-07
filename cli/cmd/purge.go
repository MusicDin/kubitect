package cmd

import (
	"cli/env"
	"cli/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// purgeCmd represents the purge command
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge the cluster directory",
	Long: `Purges the cluster directory. If cluster is not specifed (using --cluster flag) 
the operation is executed on the 'default' cluster.`,

	Run: func(cmd *cobra.Command, args []string) {
		err := purge()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(purgeCmd)

	purgeCmd.PersistentFlags().StringVar(&env.ClusterName, "cluster", env.DefaultClusterName, "specify the cluster to be used")
	purgeCmd.PersistentFlags().BoolVar(&env.AutoApprove, "auto-approve", false, "automatically approve any user permission requests")

	// Auto complete cluster names of inactive clusters for flag 'cluster'.
	purgeCmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		clusterNames, err := GetClusters([]ClusterFilter{IsInactive})
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return clusterNames, cobra.ShellCompDirectiveNoFileComp
	})
}

// purge function removes cluster's directory along with all of its content.
// Active clusters (clusters that contain terraform state file) cannot be
// purged and therfore the action will be aborted.
func purge() error {

	// Fail if cluster path is not pointing on a valid cluster directory.
	err := utils.VerifyClusterDir(env.ClusterPath)
	if err != nil {
		return err
	}

	tfStatePath := filepath.Join(env.ClusterPath, env.ConstTerraformStatePath)

	// Abort purge if terraform state file exists.
	if utils.Exists(tfStatePath) {
		return fmt.Errorf("Only destroyed clusters can be purged! Cluster '%s' is still active and therefore cannot be purged.", env.ClusterName)
	}

	// Ask user for permission.
	confirm := utils.AskUserConfirmation("The '%s' cluster directory will be removed.", env.ClusterName)
	if !confirm {
		return fmt.Errorf("User aborted.")
	}

	fmt.Printf("Purging '%s' cluster...\n", env.ClusterName)

	// Remove terraform state file
	err = os.RemoveAll(env.ClusterPath)
	if err != nil {
		return fmt.Errorf("Failed removing cluster's directory: %v", err)
	}

	fmt.Printf("Cluster '%s' has been successfully purged.\n", env.ClusterName)

	return nil
}
