package cluster

import (
	"fmt"
	"github.com/MusicDin/kubitect/cli/pkg/utils/file"
	validation2 "github.com/MusicDin/kubitect/cli/pkg/utils/validation"
)

// readConfig reads configuration file on the given path and converts it into
// the provided model.
func readConfig[T validation2.Validatable](path string, model T) (*T, error) {
	if !file.Exists(path) {
		return nil, fmt.Errorf("file '%s' does not exist", path)
	}

	return file.ReadYaml(path, model)
}

// readConfig reads configuration file on the given path and converts it into
// the provided model. If file on the provided path does not exist, neither error
// nor model is returned.
func readConfigIfExists[T validation2.Validatable](path string, model T) (*T, error) {
	if !file.Exists(path) {
		return nil, nil
	}

	return file.ReadYaml(path, model)
}

// validateConfig validates provided configuration file.
func validateConfig[T validation2.Validatable](config T) []error {
	var errs []error

	err := config.Validate()

	if err == nil {
		return nil
	}

	for _, e := range err.(validation2.ValidationErrors) {
		errs = append(errs, NewValidationError(e.Error(), e.Namespace))
	}

	return errs
}
