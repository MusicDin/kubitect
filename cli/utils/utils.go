package utils

import (
	"cli/env"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetEnv check if environment variable with name 'envKey' exists.
// Environment variable is returned if it exists, otherwise 'defaultValue'
// is returned.
func GetEnv(envKey, defaultValue string) string {

	envValue, envExists := os.LookupEnv(envKey)
	if envExists {
		return envValue
	}
	return defaultValue
}

func Exists(path string) bool {

	_, err := os.Stat(path)
	if err != nil {

		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	}
	return true
}

// ForceMove forcibly moves a file or directory to a specified location.
// First the destination file or directory is removed, and then the contents
// are moved there.
func ForceMove(srcPath string, dstPath string) error {

	err := os.RemoveAll(dstPath)
	if err != nil {
		return fmt.Errorf("Failed to force remove destination file: %w", err)
	}

	dstDir := filepath.Dir(dstPath)

	// Create all destination subdirectories if missing.
	err = os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to create dst directory (%s) while moving files: %v", dstDir, err)
	}

	err = os.Rename(srcPath, dstPath)
	if err != nil {
		return fmt.Errorf("Failed to move file from src (%s) to dst (%s) path: %w", srcPath, dstPath, err)
	}

	return nil
}

// StrArrayContains returns true if string array 'arr' contains string 'value'.
func StrArrayContains(arr []string, value string) bool {
	for _, s := range arr {
		if value == s {
			return true
		}
	}
	return false
}

// AskUserConfirmation ask user for confirmation. Function returns true if user
// types any variant of "y" or "yes" and false if user types any variant of "n"
// or "no". Otherwise user is asked again.
func AskUserConfirmation(warning ...any) bool {

	var response string

	// Automatically approve if '--auto-approve' flag is used
	if env.AutoApprove {
		return true
	}

	if len(warning) > 0 {
		PrintWarning(warning...)
	}

	fmt.Println("\nAre you sure you want to continue? (yes/no)")

	_, err := fmt.Scan(&response)
	if err != nil {
		panic(fmt.Errorf("Error asking user for confirmation: %w", err))
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		return AskUserConfirmation(warning...)
	}
}

// VerifyClusterDir returns an error if the provided path is pointing to
// an invalid cluster directory.
func VerifyClusterDir(clusterPath string) error {

	err := verifyClusterDir(clusterPath)

	if err != nil {
		PrintError("Cluster path points to an invalid cluster directory!")

		if env.Local {
			PrintError("Are you sure you are in the right directory?")
		}
		return err
	}

	return nil
}

// verifyClusterDir verifies if the provided cluster directory exists and if it
// contains necessary directories/files that represent a cluster directory. It
// returnes true if specific cluster directories are present. Otherwise it
// returns false.
func verifyClusterDir(clusterPath string) error {

	// Check if cluster directory exists
	_, err := os.Stat(clusterPath)
	if err != nil {
		return err
	}

	// Check if ansible directory exists
	for _, path := range env.ProjectRequiredFiles {

		// Check if all required directories are present
		if strings.HasSuffix(path, "/") {

			_, err = os.Stat(filepath.Join(clusterPath, path))

			if err != nil {
				return err
			}
		}
	}

	return nil
}
