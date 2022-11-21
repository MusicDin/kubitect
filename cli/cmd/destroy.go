package cmd

import (
	"cli/actions"
	"cli/env"
	"fmt"

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

	if o.ClusterName == "" {
		return fmt.Errorf("a valid (non-empty) cluster name must be provided")
	}

	clusters, err := actions.Clusters(o.Context())

	if err != nil {
		return err
	}

	c := clusters.FindByName(o.ClusterName)

	if c == nil {
		return fmt.Errorf("cluster '%s' not found.", c.Name)
	}

	count := clusters.CountByName(c.Name)

	if count > 1 {
		return fmt.Errorf("cannot destroy the cluster: multiple clusters (%d) have been found with the name '%s'", count, c.Name)
	}

	return c.Cluster(o.Context()).Destroy()
}
