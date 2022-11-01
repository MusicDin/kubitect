package actions

import (
	"cli/env"
	"cli/utils"
	"fmt"
	"os"
	"path/filepath"
)

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
		return fmt.Errorf("Cluster path (%s) does not exist!", clusterPath)
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
