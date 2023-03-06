package event

import "github.com/MusicDin/kubitect/pkg/ui"

func NewConfigChangeError(msg string, paths ...string) error {
	return ui.NewErrorBlock(ui.ERROR,
		[]ui.Content{
			ui.NewErrorLine("Error type:", "Config Change"),
			ui.NewErrorSection("Config path:", paths...),
			ui.NewErrorSection("Error:", msg),
		},
	)
}

func NewConfigChangeWarning(msg string, paths ...string) error {
	return ui.NewErrorBlock(ui.WARN,
		[]ui.Content{
			ui.NewErrorLine("Warning type:", "Config Change"),
			ui.NewErrorSection("Config path:", paths...),
			ui.NewErrorSection("Warning:", msg),
		},
	)
}
