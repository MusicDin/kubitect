package cmd

import "github.com/spf13/cobra"

var (
	exportShort = "Export specific configuration file"
	exportLong  = LongDesc(`
		Exports specific configuration file`)

	exportExample = Example(`
		Export kubeconfig for cluster 'cls-name':
		> kubitect export kubeconfig --cluster cls-name`)
)

func NewExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "export",
		GroupID: "support",
		Short:   exportShort,
		Long:    exportLong,
		Example: exportExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddGroup(
		&cobra.Group{
			ID:    "main",
			Title: "Commands:",
		},
	)

	cmd.AddCommand(NewExportKcCmd())
	cmd.AddCommand(NewExportConfigCmd())

	return cmd
}
