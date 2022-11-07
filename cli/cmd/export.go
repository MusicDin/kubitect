package cmd

import (
	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:     "export",
	GroupID: "support",
	Short:   "Export specific configuration file",
	Long: `
Exports specific configuration file`,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
