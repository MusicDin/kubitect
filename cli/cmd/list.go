package cmd

import (
	"cli/actions"

	"github.com/spf13/cobra"
)

var desc = `
Lists clusters. 
	
Local clusters are also listed if current directory is a Kubitect project.
`

func init() {
	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists clusters",
		Long:    desc,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return actions.ListClusters()
		},
	}

	rootCmd.AddCommand(listCmd)
}
