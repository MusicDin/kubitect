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

// ForceMove forcfully moves a file or a directory to the specific location.
func ForceMove(srcPath string, dstPath string) error {

	err := os.RemoveAll(dstPath)
	if err != nil {
		return fmt.Errorf("Failed to force remove destination file: %w", err)
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
		format := fmt.Sprint(warning[0])
		args := warning[1:]
		PrintWarning(fmt.Sprintf(format, args...))
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

	if !isClusterDir(clusterPath) {
		PrintError("Cluster path points to the invalid cluster directory!")

		if env.Local {
			PrintError("Are you sure you are in the right directory?")
		}

		return fmt.Errorf("Invalid cluster directory.")
	}

	return nil
}

// IsClusterDir verifies if the provided cluster directory exists and if it
// contains necessary directories/files that represent a cluster directory. It
// returnes true if specific cluster directories are present. Otherwise it
// returns false.
func isClusterDir(clusterPath string) bool {

	// Check if cluster directory exists
	_, err := os.Stat(clusterPath)
	if err != nil {
		return false
	}

	// Check if ansible directory exists
	for _, path := range env.ProjectRequiredFiles {

		// Check if all required directories are present
		if strings.HasSuffix(path, "/") {
			_, err = os.Stat(filepath.Join(clusterPath, path))
			if err != nil {
				return false
			}
		}
	}

	return true
}
