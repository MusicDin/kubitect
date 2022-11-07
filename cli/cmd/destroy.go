package cmd

import (
	"cli/actions"
	"cli/env"

	"github.com/spf13/cobra"
)

func init() {
	var clusterName string

	destroyCmd := &cobra.Command{
		SuggestFor: []string{"remove", "delete", "del"},
		Use:        "destroy",
		GroupID:    "mgmt",
		Short:      "Destroy the cluster",
		Long: `
Destroys the cluster with a given name.`,

		RunE: func(cmd *cobra.Command, args []string) error {
			return actions.Destroy(clusterName)
		},
	}

	rootCmd.AddCommand(destroyCmd)

	destroyCmd.PersistentFlags().StringVar(&clusterName, "cluster", "", "specify the cluster to be used")
	destroyCmd.PersistentFlags().BoolVarP(&env.Local, "local", "l", false, "use a current directory as the cluster path")
	destroyCmd.PersistentFlags().BoolVar(&env.AutoApprove, "auto-approve", false, "automatically approve any user permission requests")

	destroyCmd.MarkPersistentFlagRequired("cluster")

	// Auto complete cluster names of active clusters for flag 'cluster'.
	destroyCmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		clusters, err := actions.GetClusters()

		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return clusters.Names(), cobra.ShellCompDirectiveNoFileComp
	})
}
