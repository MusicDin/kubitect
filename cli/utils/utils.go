package utils

import (
	"fmt"
	"os"
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
func AskUserConfirmation() bool {

	var response string

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
		return AskUserConfirmation()
	}
}
