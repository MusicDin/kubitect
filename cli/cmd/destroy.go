package cmd

import (
	"cli/env"
	"cli/helpers"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroys the cluster",
	Long: `Destroys the cluster. If cluster is not specifed (using --cluster flag) 
the operation is executed on the 'default' cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := destroy()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	destroyCmd.PersistentFlags().StringVar(&env.ClusterName, "cluster", env.DefaultClusterName, "specify the cluster to be used")
	destroyCmd.PersistentFlags().BoolVar(&env.Local, "local", false, "use a current directory as the cluster path")

	// Auto complete cluster names from project clusters directory
	// for flag 'cluster'.
	destroyCmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		clustersPath := filepath.Join(env.ProjectHomePath, env.ConstProjectClustersDir)
		return []string{clustersPath}, cobra.ShellCompDirectiveFilterDirs
	})
}

// destroy function destroys the cluster.
func destroy() error {

	var err error

	fmt.Printf("Destroying cluster '%s'...\n", env.ClusterName)

	// Terraform apply
	err = helpers.TerraformDestroy(env.ClusterPath)
	if err != nil {
		return err
	}

	return nil
}
