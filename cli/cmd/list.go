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
		Short:   "Lists objects",
		Long:    `Command that lists specified objects.`,
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
}
