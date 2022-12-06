package cmd

import (
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

	GenericOptions
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
	cmd.PersistentFlags().BoolVar(&opts.AutoApprove, "auto-approve", false, "automatically approve any user permission requests")
	cmd.PersistentFlags().BoolVar(&opts.Debug, "debug", false, "enable debug messages")

	cmd.MarkPersistentFlagRequired("cluster")

	cmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		clusters, err := AllClusters(opts.GlobalContext())

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

	clusters, err := AllClusters(o.GlobalContext())

	if err != nil {
		return err
	}

	c := clusters.FindByName(o.ClusterName)

	if c == nil {
		return fmt.Errorf("cluster '%s' does not exist", o.ClusterName)
	}

	count := clusters.CountByName(o.ClusterName)

	if count > 1 {
		return fmt.Errorf("multiple clusters (%d) have been found with the name '%s'", count, o.ClusterName)
	}

	return c.Destroy()
}
