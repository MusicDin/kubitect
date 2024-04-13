package main

import (
	"github.com/MusicDin/kubitect/embed"
	"github.com/MusicDin/kubitect/pkg/app"
	"github.com/MusicDin/kubitect/pkg/ui"

	"github.com/spf13/cobra"
)

var (
	listPresetsShort = "List presets"
	listPresetsLong  = LongDesc(`
		Command list presets lists all available cluster configuration presets.`)

	listPresetsExample = Example(`
		List all presets:
		> kubitect list presets`)
)

type ListPresetsOptions struct {
	app.AppContextOptions
}

func NewListPresetsCmd() *cobra.Command {
	var o ListPresetsOptions

	cmd := &cobra.Command{
		Use:     "presets",
		Aliases: []string{"preset"},
		GroupID: "main",
		Short:   listPresetsShort,
		Long:    listPresetsLong,
		Example: listPresetsExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run()
		},
	}

	return cmd
}

func (o *ListPresetsOptions) Run() error {
	presets, err := embed.Presets()
	if err != nil {
		return err
	}

	ui.Println(ui.INFO, "Available presets:")
	for _, p := range presets {
		ui.Printf(ui.INFO, "- %s\n", presetName(p.Name))
	}

	return nil
}
