package cmd

import (
	"cli/actions"
	"cli/env"

	"github.com/spf13/cobra"
)

var (
	destroyShort = "Destroy the cluster"
	destroyLong  = LongDesc(`
		Destroy the cluster with a given name.`)

	destroyExample = Example(`
		To destroy a cluster named 'cls':
		> kubitect destroy --cluster cls`)
)

type DestroyOptions struct {
	ClusterName string

	env.ContextOptions
}

func NewDestroyCmd() *cobra.Command {
	var opts DestroyOptions

	cmd := &cobra.Command{
		SuggestFor: []string{"remove", "delete", "del"},
		Use:        "destroy",
		GroupID:    "mgmt",
		Short:      destroyShort,
		Long:       destroyLong,
		Example:    destroyExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run()
		},
	}

	cmd.PersistentFlags().StringVar(&opts.ClusterName, "cluster", "", "specify the cluster to be used")
	cmd.PersistentFlags().BoolVarP(&opts.Local, "local", "l", false, "use a current directory as the cluster path")
	cmd.PersistentFlags().BoolVar(&env.AutoApprove, "auto-approve", false, "automatically approve any user permission requests")
	cmd.PersistentFlags().BoolVar(&env.Debug, "debug", false, "enable debug messages")

	cmd.MarkPersistentFlagRequired("cluster")

	cmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		clusters, err := actions.Clusters(opts.Context())

		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return clusters.Names(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func (o *DestroyOptions) Run() error {
	return actions.Destroy(o.Context(), o.ClusterName)
}
