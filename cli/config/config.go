package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/goccy/go-yaml"
)

// GetStrValue extracts value for the given key from provided config file.
func GetStrValue(configPath string, key string) (string, error) {

	var value string

	if len(key) < 1 {
		return "", fmt.Errorf("Argument 'key' cannot be an empty string!")
	}

	if len(configPath) < 1 {
		return "", fmt.Errorf("Argument 'clusterPath' cannot be an empty string!")
	}

	// Read file on provided path
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("Failed reading yaml file on path '%s': %w", configPath, err)
	}

	// Construct PathString for yaml parser
	path, err := yaml.PathString("$." + key)
	if err != nil {
		return "", fmt.Errorf("Failed to create PathString for key='%s': %w", key, err)
	}

	// Parse specific value from provided yaml file.
	err = path.Read(strings.NewReader(string(buf)), &value)
	if err != nil {
		return "", fmt.Errorf("Failed parsing value for key='%s' from file '%s': %w", key, configPath, err)
	}

	return value, nil
}
