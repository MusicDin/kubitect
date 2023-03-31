package cmd

import (
	"strings"

	"github.com/MusicDin/kubitect/pkg/app"
	"github.com/MusicDin/kubitect/pkg/ui"

	"github.com/spf13/cobra"
)

var (
	listClustersShort = "List clusters"
	listClustersLong  = LongDesc(`
		Command list clusters lists all clusters including local clusters if 
		a current (working) directory is Kubitect project.`)

	listClusterExample = Example(`
		List all clusters:
		> kubitect list clusters`)
)

type ListClustersOptions struct {
	app.AppContextOptions
}

func NewListClustersCmd() *cobra.Command {
	var o ListClustersOptions

	cmd := &cobra.Command{
		Use:     "clusters",
		Aliases: []string{"cluster"},
		GroupID: "main",
		Short:   listClustersShort,
		Long:    listClustersLong,
		Example: listClusterExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run()
		},
	}

	return cmd
}

func (o *ListClustersOptions) Run() error {
	ac := o.AppContext()

	clusters, err := AllClusters(ac)

	if err != nil {
		return err
	}

	if len(clusters) == 0 {
		ui.Println(ui.INFO, "No clusters initialized yet. Run 'kubitect apply' to create the cluster.")
		return nil
	}

	ui.Println(ui.INFO, "Clusters:")

	for _, c := range clusters {
		var opt []string

		if c.ContainsTfStateConfig() {
			opt = append(opt, "active")
		}

		if c.Local {
			opt = append(opt, "local")
		}

		if len(opt) > 0 {
			ui.Printf(ui.INFO, "  - %s (%s)\n", c.Name, strings.Join(opt, ", "))
		} else {
			ui.Printf(ui.INFO, "  - %s\n", c.Name)
		}
	}

	return nil
}
