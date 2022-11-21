package cmd

import (
	"cli/actions"
	"cli/env"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	listShort = "Lists clusters"
	listLong  = LongDesc(`
		Lists all clusters located in the project directory.
		Local clusters are also listed if current (working) directory is a Kubitect project.`)
)

type ListOptions struct {
	env.ContextOptions
}

func NewListCmd() *cobra.Command {
	var opts ListOptions

	return &cobra.Command{
		Aliases:    []string{"ls"},
		SuggestFor: []string{"show"},
		Use:        "list",
		GroupID:    "support",
		Short:      listShort,
		Long:       listLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run()
		},
	}
}

func (o *ListOptions) Run() error {
	clusters, err := actions.Clusters(o.Context())

	if err != nil {
		return err
	}

	if len(clusters) == 0 {
		fmt.Println("No clusters initialized yet. Run 'kubitect apply' to create the cluster.")
		return nil
	}

	fmt.Println("Clusters:")

	for _, c := range clusters {
		var opt []string

		if c.ContainsTfStateConfig() {
			opt = append(opt, "active")
		}

		if c.Local {
			opt = append(opt, "local")
		}

		if len(opt) > 0 {
			fmt.Printf("  - %s (%s)\n", c.Name, strings.Join(opt, ", "))
		} else {
			fmt.Printf("  - %s\n", c.Name)
		}
	}

	return nil
}
