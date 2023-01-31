package cmd

import (
	"cli/app"
	"cli/cluster"
	"cli/env"

	"github.com/spf13/cobra"
)

const DefaultAction = "create"

var (
	applyShort = "Create, scale or upgrade the cluster"
	applyLong  = LongDesc(`
		Apply new configuration file to create a cluster, or to modify, scale or upgrade the existing one.`)

	applyExample = Example(`
		Create a new cluster or modify an existing one:
		> kubitect apply --config cluster.yaml
		
		To upgrade an existing cluster, bump the Kubernetes version in current cluster config and run:
		> kubitect apply --config cluster.yaml --action upgrade
		
		To scale an existing cluster, add or remove node instances in current cluster config and run:
		> kubitect apply --config cluster.yaml --action scale`)
)

type ApplyOptions struct {
	Config string
	Action string

	app.AppContextOptions
}

func NewApplyCmd() *cobra.Command {
	var opts ApplyOptions

	cmd := &cobra.Command{
		SuggestFor: []string{"create", "scale", "upgrade"},
		Use:        "apply",
		GroupID:    "mgmt",
		Short:      applyShort,
		Long:       applyLong,
		Example:    applyExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run()
		},
	}

	cmd.PersistentFlags().StringVarP(&opts.Config, "config", "c", "", "specify path to the cluster config file")
	cmd.PersistentFlags().StringVarP(&opts.Action, "action", "a", DefaultAction, "specify cluster action [create, upgrade, scale]")
	cmd.PersistentFlags().BoolVarP(&opts.Local, "local", "l", false, "use a current directory as the cluster path")
	cmd.PersistentFlags().BoolVar(&opts.AutoApprove, "auto-approve", false, "automatically approve any user permission requests")
	cmd.PersistentFlags().BoolVar(&opts.Debug, "debug", false, "enable debug messages")

	cmd.MarkPersistentFlagRequired("config")

	cmd.RegisterFlagCompletionFunc("action", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return env.ProjectApplyActions[:], cobra.ShellCompDirectiveDefault
	})

	return cmd
}

func (o *ApplyOptions) Run() error {
	c, err := cluster.NewCluster(o.AppContext(), o.Config)

	if err != nil {
		return err
	}

	return c.Apply(o.Action)
}
