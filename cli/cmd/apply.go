package cmd

import (
	"cli/actions"
	"cli/env"

	"github.com/spf13/cobra"
)

const (
	defaultAction = "create"
)

var (
	configPath string
	action     string
)

var applyCmd = &cobra.Command{
	SuggestFor: []string{"create", "scale", "upgrade"},
	Use:        "apply",
	GroupID:    "mgmt",
	Short:      "Create, scale or upgrade the cluster",
	Long: `
Apply new configuration file to create a cluster, or scale or upgrade the existing one.`,

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
}
