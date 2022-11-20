package actions

import (
	"cli/ui"
	"fmt"
	"strings"
)

func NewValidationError(msg string, path string) error {
	return ui.ErrorBlock{
		Level: ui.ERROR,
		Content: []ui.ErrorContent{
			ui.NewErrorLine("Error type:", "Validation Error"),
			ui.NewErrorSection("Config path:", path),
			ui.NewErrorSection("Error:", strings.Split(msg, "\n")...),
		},
	}
}

func NewConfigChangeError(msg string, paths ...string) error {
	return ui.ErrorBlock{
		Level: ui.ERROR,
		Content: []ui.ErrorContent{
			ui.NewErrorLine("Error type:", "Config Change"),
			ui.NewErrorSection("Config paths:", paths...),
			ui.NewErrorSection("Error:", strings.Split(msg, "\n")...),
		},
	}
}

func NewConfigChangeWarning(msg string, paths ...string) error {
	return ui.ErrorBlock{
		Level: ui.WARN,
		Content: []ui.ErrorContent{
			ui.NewErrorLine("Warning type:", "Config Change"),
			ui.NewErrorSection("Config path:", paths...),
			ui.NewErrorSection("Warning:", strings.Split(msg, "\n")...),
		},
	}
}

func NewInvalidWorkingDirError(missingFiles []string) error {
	e := []string{
		"Current directory is missing some required files.\n",
		"Are you sure you are in the right directory?",
	}

	return ui.ErrorBlock{
		Level: ui.ERROR,
		Content: []ui.ErrorContent{
			ui.NewErrorLine("Error type:", "Invalid working directory"),
			ui.NewErrorSection("Missing files:", missingFiles...),
			ui.NewErrorSection("Error:", e...),
		},
	}
}

func NewInvalidProjectDirError(clusterPath string, missingFiles ...string) error {
	e := fmt.Sprintf("Cluster directory (%s) is missing some required files.", clusterPath)

	return ui.ErrorBlock{
		Level: ui.ERROR,
		Content: []ui.ErrorContent{
			ui.NewErrorLine("Error type:", "Invalid project directory"),
			ui.NewErrorSection("Missing files:", missingFiles...),
			ui.NewErrorSection("Error:", e),
		},
	}
}
