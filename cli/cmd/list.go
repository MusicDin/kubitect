package cmd

import (
	"cli/actions"

	"github.com/spf13/cobra"
)

func init() {
	listCmd := &cobra.Command{
		Aliases:    []string{"ls"},
		SuggestFor: []string{"show"},
		Use:        "list",
		GroupID:    "support",
		Short:      "Lists clusters",
		Long: `
Lists clusters.

Local clusters are also listed if current directory is a Kubitect project.`,

		RunE: func(cmd *cobra.Command, args []string) error {
			return actions.ListClusters()
		},
	}

	rootCmd.AddCommand(listCmd)
}
