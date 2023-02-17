package cmd

import (
	"fmt"
	"os"

	"github.com/MusicDin/kubitect/cli/app"
	"github.com/MusicDin/kubitect/cli/utils/file"

	"github.com/spf13/cobra"
)

var (
	exportConfigShort = "Export cluster config file"
	exportConfigLong  = LongDesc(`
		Command export config outputs cluster's configuration file to standard output.`)

	exportConfigExample = Example(`
		To save a config to the specific file, redirect command output to that file:
		> kubitect export config --cluster lake > cls.yaml`)
)

type ExportConfigOptions struct {
	ClusterName string

	app.AppContextOptions
}

func NewExportConfigCmd() *cobra.Command {
	var o ExportConfigOptions

	cmd := &cobra.Command{
		SuggestFor: []string{"cfg"},
		Use:        "config",
		GroupID:    "main",
		Short:      exportConfigShort,
		Long:       exportConfigLong,
		Example:    exportConfigExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run()
		},
	}

	cmd.PersistentFlags().StringVar(&o.ClusterName, "cluster", "", "specify the cluster to be used")
	cmd.MarkPersistentFlagRequired("cluster")

	cmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var names []string

		clusters, err := AllClusters(o.AppContext())

		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		for _, c := range clusters {
			if c.ContainsAppliedConfig() {
				names = append(names, c.Name)
			}
		}

		return names, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func (o *ExportConfigOptions) Run() error {
	cs, err := AllClusters(o.AppContext())

	c := cs.FindByName(o.ClusterName)

	if c == nil {
		return fmt.Errorf("cluster '%s' does not exist", o.ClusterName)
	}

	count := cs.CountByName(o.ClusterName)

	if count > 1 {
		return fmt.Errorf("multiple clusters (%d) have been found with the name '%s'", count, o.ClusterName)
	}

	if !c.ContainsAppliedConfig() {
		return fmt.Errorf("cluster '%s' does not contain a config file", o.ClusterName)
	}

	config, err := file.Read(c.AppliedConfigPath())

	if err != nil {
		return err
	}

	fmt.Fprint(os.Stdout, config)

	return nil
}
