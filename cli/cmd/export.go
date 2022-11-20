package cmd

import (
	"github.com/spf13/cobra"
)

var (
	exportShort = "Export specific configuration file"
	exportLong  = LongDesc(`
		Exports specific configuration file`)
)

func NewExportCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "export",
		GroupID: "support",
		Short:   exportShort,
		Long:    exportLong,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
}
