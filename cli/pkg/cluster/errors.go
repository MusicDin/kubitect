package cluster

import (
	ui2 "github.com/MusicDin/kubitect/cli/pkg/ui"
)

func NewInvalidClusterDirError(missingFiles []string) error {
	return ui2.NewErrorBlock(ui2.ERROR,
		[]ui2.Content{
			ui2.NewErrorLine("Error type:", "Invalid working directory"),
			ui2.NewErrorSection("Missing files:", missingFiles...),
		},
	)
}

func NewValidationError(msg string, path string) error {
	return ui2.NewErrorBlock(ui2.ERROR,
		[]ui2.Content{
			ui2.NewErrorLine("Error type:", "Validation Error"),
			ui2.NewErrorSection("Config path:", path),
			ui2.NewErrorSection("Error:", msg),
		},
	)
}
