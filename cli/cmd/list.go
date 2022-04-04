package cmd

import (
	"github.com/spf13/cobra"
)

var (
	listCmdAliases = []string{
		"ls",
	}

	// listCmd represents the list command
	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: listCmdAliases,
		Short:   "Lists initialized clusters",
		Long: `Lists clusters that have been initialized in $TKK_HOME.
Local clusters (applied with --local flag) are not listed.`,
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
}
