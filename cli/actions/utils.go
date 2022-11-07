package actions

import (
	"cli/env"
	"cli/utils"
	"cli/validation"
	"fmt"
	"os"
	"path/filepath"
)

// readConfig reads configuration file on the given path and converts it into
// the provided model.
func readConfig[T validation.Validatable](path string, model T) (*T, error) {
	if !utils.Exists(path) {
		return nil, fmt.Errorf("file '%s' does not exist", path)
	}

	return utils.ReadYaml(path, model)
}

// readConfig reads configuration file on the given path and converts it into
// the provided model. If file on the provided path does not exist, neither error
// nor model is returned.
func readConfigIfExists[T validation.Validatable](path string, model T) (*T, error) {
	if !utils.Exists(path) {
		return nil, nil
	}

	return utils.ReadYaml(path, model)
}

// validateConfig validates provided configuration file.
func validateConfig[T validation.Validatable](config T) error {
	var errs utils.Errors

	err := config.Validate()

	if err == nil {
		return nil
	}

	for _, e := range err.(validation.ValidationErrors) {
		errs = append(errs, NewValidationError(e.Error(), e.Namespace))
	}

	return errs
}

// copyReqFiles copies project required files from source directory
// to the destination directory.
func copyReqFiles(srcDir, dstDir string) error {
	if err := verifyClusterDir(srcDir); err != nil {
		return err
	}

	for _, path := range env.ProjectRequiredFiles {
		src := filepath.Join(srcDir, path)
		dst := filepath.Join(dstDir, path)

		if err := utils.ForceCopy(src, dst); err != nil {
			return err
		}
	}

	return verifyClusterDir(dstDir)
}

// verifyClusterDir verifies if the provided cluster directory
// exists and if it contains all necessary directories.
func verifyClusterDir(clusterPath string) error {
	if !utils.Exists(clusterPath) {
		return fmt.Errorf("cluster path '%s' does not exist", clusterPath)
	}

	var missing []string

	for _, path := range env.ProjectRequiredFiles {
		p := filepath.Join(clusterPath, path)

		if !utils.Exists(p) {
			missing = append(missing, path)
		}
	}

	if len(missing) == 0 {
		return nil
	}

	wd, err := os.Getwd()

	if err != nil {
		return err
	}

	if wd == clusterPath {
		return NewInvalidWorkingDirError(missing)
	}

	return NewInvalidProjectDirError(clusterPath, missing...)
}
