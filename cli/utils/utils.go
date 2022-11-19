package utils

import (
	"cli/env"
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
func AskUserConfirmation() error {
	// Automatically approve if '--auto-approve' flag is used
	if env.AutoApprove {
		return nil
	}

	fmt.Println("\nWould you like to continue? (yes/no)")

	var response string

	if _, err := fmt.Scan(&response); err != nil {
		return fmt.Errorf("Error occurred while asking user for the confirmation: %v", err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return nil
	case "n", "no":
		return fmt.Errorf("User aborted...")
	default:
		return AskUserConfirmation()
	}
}
