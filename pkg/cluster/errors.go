package cluster

import "github.com/MusicDin/kubitect/pkg/ui"

func NewInvalidClusterDirError(missingFiles []string) error {
	return ui.NewErrorBlock(ui.ERROR,
		[]ui.Content{
			ui.NewErrorLine("Error type:", "Invalid working directory"),
			ui.NewErrorSection("Missing files:", missingFiles...),
		},
	)
}

func NewValidationError(msg string, path string) error {
	return ui.NewErrorBlock(ui.ERROR,
		[]ui.Content{
			ui.NewErrorLine("Error type:", "Validation Error"),
			ui.NewErrorSection("Config path:", path),
			ui.NewErrorSection("Error:", msg),
		},
	)
}

func NewConfigChangeError(msg string, paths ...string) error {
	return ui.NewErrorBlock(ui.ERROR,
		[]ui.Content{
			ui.NewErrorLine("Error type:", "Invalid Configuration Change"),
			ui.NewErrorSection("Config path:", paths...),
			ui.NewErrorSection("Error:", msg),
		},
	)
}

func NewConfigChangeWarning(msg string, paths ...string) error {
	return ui.NewErrorBlock(ui.WARN,
		[]ui.Content{
			ui.NewErrorLine("Warning type:", "Dangerous Configuration Change"),
			ui.NewErrorSection("Config path:", paths...),
			ui.NewErrorSection("Warning:", msg),
		},
	)
}
