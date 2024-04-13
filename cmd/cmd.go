package main

import (
	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/spf13/cobra"
)

var (
	rootLong = LongDesc(`
		Kubitect is a CLI tool that helps you manage multiple Kubernetes clusters.`)
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "kubitect",
		Short:   "Kubitect",
		Long:    rootLong,
		Version: env.ConstProjectVersion,
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SuggestionsMinimumDistance = 3

	cmd.AddGroup(
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

	cmd.AddCommand(NewApplyCmd())
	cmd.AddCommand(NewDestroyCmd())
	cmd.AddCommand(NewExportCmd())
	cmd.AddCommand(NewListCmd())

	cmd.SetCompletionCommandGroupID("other")
	cmd.SetHelpCommandGroupID("other")

	return cmd
}
