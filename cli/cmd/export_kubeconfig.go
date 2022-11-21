package cmd

import (
	"cli/actions"
	"cli/env"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	exportKcShort = "Export cluster kubeconfig file"

	exportKcLong = LongDesc(`
		Command export kubeconfig outputs cluster's kubeconfig file to standard output.`)

	exportKcExample = Example(`
		To save kubeconfig to a specific file, redirect command output to a file:
		> kubitect export kubeconfig --cluster lake > lake.yaml
					
		Use kubeconfig with kubectl to access cluster:
		> kubectl --kubeconfig lake.yaml get nodes`)
)

type ExportKcOptions struct {
	ClusterName string

	env.ContextOptions
}

func NewExportKcCmd() *cobra.Command {
	var opts ExportKcOptions

	cmd := &cobra.Command{
		SuggestFor: []string{"kubecfg", "kube", "kc"},
		Use:        "kubeconfig",
		GroupID:    "main",
		Short:      exportKcShort,
		Long:       exportKcLong,
		Example:    exportKcExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run()
		},
	}

	// cmd.PersistentFlags().BoolVarP(&opts.Local, "local", "l", false, "use a current directory as the cluster path")
	cmd.PersistentFlags().StringVar(&opts.ClusterName, "cluster", "", "specify the cluster to be used")
	cmd.MarkPersistentFlagRequired("cluster")

	cmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// clusterNames, err := GetClusters([]ClusterFilter{IsActive, ContainsKubeconfig})
		clusters, err := actions.Clusters(opts.Context())

		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return clusters.Names(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func (o *ExportKcOptions) Run() error {
	cs, err := actions.Clusters(o.Context())

	c := cs.FindByName(o.ClusterName)

	if c == nil {
		return fmt.Errorf("cluster '%s' does not exist: %v", o.ClusterName, err)
	}

	count := cs.CountByName(o.ClusterName)

	if count > 1 {
		return fmt.Errorf("multiple clusters (%d) have been found with the name '%s'", count, o.ClusterName)
	}

	if !c.ContainsKubeconfig() {
		return fmt.Errorf("cluster '%s' does not have a Kubeconfig file", o.ClusterName)
	}

	kc, err := c.Kubeconfig()

	if err != nil {
		return err
	}

	fmt.Fprint(os.Stdout, kc)

	return nil
}
