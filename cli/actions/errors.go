package actions

import (
	"cli/utils"
	"fmt"
	"strings"
)

func NewValidationError(msg string, path string) error {
	content := []utils.ErrorContent{
		utils.NewErrorLine("Error type:", "Validation Error"),
		utils.NewErrorSection("Config path:", path),
		utils.NewErrorSection("Error:", strings.Split(msg, "\n")...),
	}

	return utils.NewErrorB(content)
}

func NewConfigChangeError(msg string, paths ...string) error {
	content := []utils.ErrorContent{
		utils.NewErrorLine("Error type:", "Config Change"),
		utils.NewErrorSection("Config paths:", paths...),
		utils.NewErrorSection("Error:", strings.Split(msg, "\n")...),
	}

	return utils.NewErrorB(content)
}

func NewConfigChangeWarning(msg string, paths ...string) error {
	content := []utils.ErrorContent{
		utils.NewErrorLine("Warning type:", "Config Change"),
		utils.NewErrorSection("Config path:", paths...),
		utils.NewErrorSection("Warning:", strings.Split(msg, "\n")...),
	}

	return utils.NewWarnB(content)
}

func NewInvalidWorkingDirError(missingFiles []string) error {
	var e []string

	e = append(e, "Current directory is missing some required files.\n")
	e = append(e, "Are you sure you are in the right directory?")

	content := []utils.ErrorContent{
		utils.NewErrorLine("Error type:", "Invalid working directory"),
		utils.NewErrorSection("Missing files:", missingFiles...),
		utils.NewErrorSection("Error:", e...),
	}

	return utils.NewErrorB(content)
}

func NewInvalidProjectDirError(clusterPath string, missingFiles ...string) error {
	e := fmt.Sprintf("Cluster directory (%s) is missing some required files.", clusterPath)

	content := []utils.ErrorContent{
		utils.NewErrorLine("Error type:", "Invalid project directory"),
		utils.NewErrorSection("Missing files:", missingFiles...),
		utils.NewErrorSection("Error:", e),
	}

	return utils.NewErrorB(content)
}
