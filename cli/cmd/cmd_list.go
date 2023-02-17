package cmd

import (
	"strings"

	"github.com/MusicDin/kubitect/cli/app"

	"github.com/MusicDin/kubitect/cli/ui"

	"github.com/spf13/cobra"
)

var (
	listShort = "Lists clusters"
	listLong  = LongDesc(`
		Lists all clusters located in the project directory.
		Local clusters are also listed if current (working) directory is a Kubitect project.`)
)

type ListOptions struct {
	app.AppContextOptions
}

func NewListCmd() *cobra.Command {
	var o ListOptions

	return &cobra.Command{
		Aliases:    []string{"ls"},
		SuggestFor: []string{"show"},
		Use:        "list",
		GroupID:    "support",
		Short:      listShort,
		Long:       listLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run()
		},
	}
}

func (o *ListOptions) Run() error {
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
