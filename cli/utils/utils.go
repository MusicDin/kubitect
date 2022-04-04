package utils

import (
	"fmt"
	"os"
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
