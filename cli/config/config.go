package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
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
	configBuf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Failed reading yaml file on path '%s': %w", configPath, err)
	}

	path, err := yaml.PathString("$." + key)
	if err != nil {
		return fmt.Errorf("Invalid yaml key: %w", err)
	}

	// Parse specific value from provided yaml file into provided object.
	err = path.Read(strings.NewReader(string(configBuf)), obj)
	if err != nil {
		return fmt.Errorf("Failed parsing value for key='%s' from file '%s': %w", key, configPath, err)
	}

	return nil
}

// ReplaceValue function reads yaml file on 'configPath', replaces node on
// the 'key' path with the 'newValue' and returns modified yaml as string.
func ReplaceValue(configPath string, key string, newValue interface{}) (string, error) {

	if len(configPath) < 1 {
		return "", fmt.Errorf("Argument 'clusterPath' cannot be an empty string!")
	}

	// Read file on the config path.
	configBuf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("Failed reading yaml file on path '%s': %w", configPath, err)
	}

	path, err := yaml.PathString("$." + key)
	if err != nil {
		return "", fmt.Errorf("Invalid yaml key: %w", err)
	}

	// Parse bytes array into file type.
	configYaml, err := parser.ParseBytes(configBuf, 0)
	if err != nil {
		return "", fmt.Errorf("Failed parsing yaml file: %w", err)
	}

	// Marshal new value into bytes array.
	newValueBuf, err := yaml.Marshal(newValue)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal replacement yaml object: %w", err)
	}

	// Replace yaml node in the config with new value.
	err = path.ReplaceWithReader(configYaml, bytes.NewReader(newValueBuf))
	if err != nil {
		return "", fmt.Errorf("Failed replacing yaml object: %w", err)
	}

	// Return new yaml as string.
	return configYaml.String(), nil
}
