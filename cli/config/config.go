package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/goccy/go-yaml"
)

// GetStrValue function extracts string value for the given key from the
// provided config file.
func GetStrValue(configPath string, key string) (string, error) {

	var value string

	err := GetValue(configPath, key, &value)

	return value, err
}

// GetStrArrValue function extracts array of strings for the given key from
// the provided config file.
func GetStrArrValue(configPath string, key string) ([]string, error) {

	var value []string

	err := GetValue(configPath, key, &value)

	return value, err
}

// GetStrValue function extracts value for the given key from the provided
// config.
func GetValue(configPath string, key string, obj interface{}) error {

	if obj == nil {
		return fmt.Errorf("Argument 'obj' cannot be nil!")
	}

	if len(configPath) < 1 {
		return fmt.Errorf("Argument 'clusterPath' cannot be an empty string!")
	}

	// Read file on provided path
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Failed reading yaml file on path '%s': %w", configPath, err)
	}

	path, err := yaml.PathString("$." + key)
	if err != nil {
		return fmt.Errorf("Invalid yaml key: %w", err)
	}

	// Parse specific value from provided yaml file.
	err = path.Read(strings.NewReader(string(buf)), obj)
	if err != nil {
		return fmt.Errorf("Failed parsing value for key='%s' from file '%s': %w", key, configPath, err)
	}

	return nil

}
