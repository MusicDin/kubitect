package cmd

import (
	"cli/env"

	"github.com/spf13/cobra"
)

var (
	rootLong = LongDesc(`Kubitect is a CLI tool that helps you manage multiple Kubernetes clusters.`)
)

// Root command (cli name)
var rootCmd = &cobra.Command{
	Use:     "kubitect",
	Short:   "Kubitect",
	Long:    rootLong,
	Version: env.ConstProjectVersion,
}

// Execute adds all child commands to the root command and sets the flags
// accordingly.
func Execute() error {
	return rootCmd.Execute()
	// return utils.FormatError(err)
}

func init() {
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.SuggestionsMinimumDistance = 3

	rootCmd.AddGroup(
		&cobra.Group{
			Title: "Cluster Management Commands:",
			ID:    "mgmt",
		},
		&cobra.Group{
			Title: "Support Commands:",
			ID:    "support",
		},
		&cobra.Group{
			Title: "Other Commands:",
			ID:    "other",
		},
	)

	rootCmd.AddCommand(
		NewApplyCmd(),
		NewDestroyCmd(),
		NewExportCmd(),
		NewListCmd(),
	)

	rootCmd.SetCompletionCommandGroupID("other")
	rootCmd.SetHelpCommandGroupID("other")
}
