package cmd

import (
	"cli/actions"
	"cli/env"

	"github.com/spf13/cobra"
)

const DefaultAction = "create"

var (
	applyShort = "Create, scale or upgrade the cluster"
	applyLong  = LongDesc(`
		Apply new configuration file to create a cluster, or scale or upgrade the existing one.`)
)

type ApplyOptions struct {
	Config string
	Action string

	env.ContextOptions
}

func NewApplyCmd() *cobra.Command {
	var opts ApplyOptions

	cmd := &cobra.Command{
		SuggestFor: []string{"create", "scale", "upgrade"},
		Use:        "apply",
		GroupID:    "mgmt",
		Short:      applyShort,
		Long:       applyLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run()
		},
	}

	cmd.PersistentFlags().StringVarP(&opts.Config, "config", "c", "", "specify path to the cluster config file")
	cmd.PersistentFlags().StringVarP(&opts.Action, "action", "a", DefaultAction, "specify cluster action [create, upgrade, scale]")
	cmd.PersistentFlags().BoolVarP(&opts.Local, "local", "l", false, "use a current directory as the cluster path")
	cmd.PersistentFlags().BoolVar(&env.AutoApprove, "auto-approve", false, "automatically approve any user permission requests")
	cmd.PersistentFlags().BoolVar(&env.Debug, "debug", false, "enable debug messages")

	cmd.MarkPersistentFlagRequired("config")

	cmd.RegisterFlagCompletionFunc("action", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return env.ProjectApplyActions[:], cobra.ShellCompDirectiveDefault
	})

	return cmd
}

func (o *ApplyOptions) Run() error {
	c, err := actions.NewCluster(o.Context(), o.Config)

	if err != nil {
		return err
	}

	return c.Apply(o.Action)
}
