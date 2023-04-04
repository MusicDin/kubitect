package cmd

import (
	"fmt"
	"os"

	"github.com/MusicDin/kubitect/embed"

	"github.com/spf13/cobra"
)

var (
	exportPresetShort = "Export cluster configuration preset"
	exportPresetLong  = LongDesc(`
		Command export preset outputs cluster configuration preset with a given name to
		the standard output.`)

	exportPresetExample = Example(`
		To list available presets run:
		> kubitect list presets
		
		To export a preset to the specific file run:
		> kubitect export preset --name minimal > cluster.yaml`)
)

type ExportPresetOptions struct {
	PresetName string
}

func NewExportPresetCmd() *cobra.Command {
	var o ExportPresetOptions

	cmd := &cobra.Command{
		Use:     "preset",
		GroupID: "main",
		Short:   exportPresetShort,
		Long:    exportPresetLong,
		Example: exportPresetExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run()
		},
	}

	cmd.PersistentFlags().StringVar(&o.PresetName, "name", "", "preset name")
	cmd.MarkPersistentFlagRequired("name")

	cmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var names []string

		presets, err := embed.Presets()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		for _, p := range presets {
			names = append(names, presetName(p.Name))
		}

		return names, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func (o *ExportPresetOptions) Run() error {
	p, err := embed.GetPreset(o.PresetName + ".yaml")
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, string(p.Content))
	return nil
}
