package cmd

import (
	"cli/actions"
	"cli/env"

	"github.com/spf13/cobra"
)

var (
	configPath string
	action     string

	defaultAction = "create"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply configuration",
	Long: `
Apply command creates a desired cluster based on the provided
configuration file. If cluster already exists, it is modified
according to the detected changes.

Cluster scaling or upgrading must be explicitly specified using
action flag.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := env.ToApplyAction(action)

		if err != nil {
			return err
		}

		return actions.Apply(configPath, a)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "specify path to the cluster config file")
	applyCmd.PersistentFlags().StringVarP(&action, "action", "a", defaultAction, "specify cluster action [create, upgrade, scale]")
	applyCmd.PersistentFlags().BoolVarP(&env.Local, "local", "l", false, "use a current directory as the cluster path")
	applyCmd.PersistentFlags().BoolVar(&env.AutoApprove, "auto-approve", false, "automatically approve any user permission requests")

	applyCmd.MarkPersistentFlagRequired("config")

	// Add completion values for flag 'action'.
	applyCmd.RegisterFlagCompletionFunc("action", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return env.ProjectApplyActions[:], cobra.ShellCompDirectiveDefault
	})

	// Auto complete cluster names for flag 'cluster'.
	applyCmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		clusters, err := actions.ReadClustersInfo()

		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return clusters.Names(), cobra.ShellCompDirectiveNoFileComp
	})
}
