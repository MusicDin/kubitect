package cmd

import (
	"github.com/spf13/cobra"
)

var (
	listShort = "List Kubitect resources"
	listLong  = LongDesc(`
		List Kubitect resources.`)
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Aliases:    []string{"ls"},
		SuggestFor: []string{"show"},
		Use:        "list",
		GroupID:    "support",
		Short:      listShort,
		Long:       listLong,
	}

	cmd.AddGroup(
		&cobra.Group{
			ID:    "main",
			Title: "Commands:",
		},
	)

	cmd.AddCommand(NewListClustersCmd())
	cmd.AddCommand(NewListPresetsCmd())

	return cmd
}
