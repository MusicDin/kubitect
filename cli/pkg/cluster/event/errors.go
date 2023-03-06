package event

import (
	ui2 "github.com/MusicDin/kubitect/cli/pkg/ui"
)

func NewConfigChangeError(msg string, paths ...string) error {
	return ui2.NewErrorBlock(ui2.ERROR,
		[]ui2.Content{
			ui2.NewErrorLine("Error type:", "Config Change"),
			ui2.NewErrorSection("Config path:", paths...),
			ui2.NewErrorSection("Error:", msg),
		},
	)
}

func NewConfigChangeWarning(msg string, paths ...string) error {
	return ui2.NewErrorBlock(ui2.WARN,
		[]ui2.Content{
			ui2.NewErrorLine("Warning type:", "Config Change"),
			ui2.NewErrorSection("Config path:", paths...),
			ui2.NewErrorSection("Warning:", msg),
		},
	)
}

// func NewInvalidClusterDirError(missingFiles []string) error {
// 	return ui.NewErrorBlock(ui.ERROR,
// 		[]ui.Content{
// 			ui.NewErrorLine("Error type:", "Invalid working directory"),
// 			ui.NewErrorSection("Missing files:", missingFiles...),
// 		},
// 	)
// }
