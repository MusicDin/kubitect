package actions

import (
	"cli/cmp"
	"cli/config/modelconfig"
	"cli/config/modelinfra"
	"cli/env"
	"cli/utils"
	"cli/validation"
	"fmt"
	"path"
	"path/filepath"
)

var (
	newCfgPath   = "kubitect.yaml"
	oldCfgPath   = "kubitect-applied.yaml"
	infraCfgPath = "infrastructure.yaml"
)

type Context struct {
	action Action
	oldCfg *modelconfig.Config
	newCfg *modelconfig.Config
	events []*OnChangeEvent
}

func plan(c Cluster, action env.ApplyAction) ([]*OnChangeEvent, error) {

	diff, err := compareConfigs(c)

	if err != nil {
		return nil, err
	}

	if len(diff.Changes()) == 0 {
		return nil, nil
	}

	fmt.Printf("Following changes have been detected:\n\n")
	fmt.Println(diff.ToYamlDiff())

	events, warns, errs := triggerEvents(diff, action)

	if len(errs) > 0 {
		return nil, errs
	}

	var msg string

	if len(warns) > 0 {
		fmt.Println(warns)
		msg = "Above warnings indicate potentially destructive actions."
	}

	err = utils.AskUserConfirmation(msg)

	return events, err
}

// CompareConfigs compares two configuration files and returns DiffNode representing
// all changes.
func compareConfigs(c Cluster) (*cmp.DiffNode, error) {
	if c.OldCfg == nil {
		return nil, nil
	}

	comp := cmp.NewComparator()
	comp.TagName = "opt"

	return comp.Compare(c.OldCfg, c.NewCfg)
}

// triggerEvents checks whether any events is triggered based on the provided
// changes and action. If any blocking event is triggered or some changes are
// not covered by any event, an error is thrown.
func triggerEvents(diff *cmp.DiffNode, action env.ApplyAction) ([]*OnChangeEvent, utils.Errors, utils.Errors) {
	var warns utils.Errors
	var errs utils.Errors

	events := events(action)

	triggered := cmp.TriggerEvents(diff, events)
	nmc := cmp.NonMatchingChanges(diff, events)

	// Changes not covered by the events are automatically
	// considered disallowed.
	if len(nmc) > 0 {
		var changes []string

		for _, ch := range nmc {
			changes = append(changes, ch.Path)
		}

		errs = append(errs, NewConfigChangeError("Disallowed changes.", changes...))
	}

	for _, t := range triggered {
		switch t.cType {
		case WARN:
			warns = append(warns, NewConfigChangeWarning(t.msg, t.triggerPaths...))
		case BLOCK:
			errs = append(errs, NewConfigChangeError(t.msg, t.triggerPaths...))
		}
	}

	return triggered, warns, errs
}

// readInfraConfig reads the configuration file produced by Terraform.
func readInfraConfig(clusterPath string) (*modelinfra.Config, error) {
	path := path.Join(clusterPath, env.ConstClusterConfigDir, infraCfgPath)

	if !utils.Exists(path) {
		return nil, fmt.Errorf("Terraform did not produce an output file. (%s)", path)
	}

	config, err := utils.ReadYaml(path, modelinfra.Config{})
	if err != nil {
		return nil, err
	}

	return config, validateInfraConfig(config)
}

// validateInfraConfig validates infrastructure configuration file.
// If validation fails, a formatted error is returned.
func validateInfraConfig(cfg *modelinfra.Config) error {
	err := cfg.Validate()

	if err == nil {
		return nil
	}

	var errs utils.Errors
	errs = append(errs, fmt.Errorf("Infrastructure configuration file produced by Terraform contains following (%d) errors:\n", len(errs)))

	for _, e := range err.(validation.ValidationErrors) {
		errs = append(errs, NewValidationError(e.Error(), e.Namespace))
	}

	return errs
}

// readOldConfig reads previously applied configuration file if it exists.
// Otherwise it returns nil.
func readOldConfig(clusterPath string) (*modelconfig.Config, error) {
	path := path.Join(clusterPath, env.ConstClusterConfigDir, oldCfgPath)

	if !utils.Exists(path) {
		return nil, nil
	}

	config, err := utils.ReadYaml(path, modelconfig.Config{})

	if err != nil {
		return nil, fmt.Errorf("Failed reading previously applied config: %v", err)
	}

	return config, nil
}

// readNewConfig reads configuration file on user-provided path.
func readNewConfig(path string) (*modelconfig.Config, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("Filepath '%s' cannot be converted to absolute path: %v", path, err)
	}

	if !utils.Exists(path) {
		return nil, fmt.Errorf("File on path '%s' does not exist.", path)
	}

	config, err := utils.ReadYaml(path, modelconfig.Config{})
	if err != nil {
		return nil, err
	}

	return config, nil
}

// storeNewConfig makes a copy of the provided (new) configuration file in
// cluster directory.
func storeNewConfig(clusterPath, configPath string) error {
	src := configPath
	dst := path.Join(clusterPath, env.ConstClusterConfigDir, newCfgPath)

	return utils.CopyFile(src, dst)
}

// applyNewConfig moves new config to the location of the applied.
func applyNewConfig(clusterPath string) error {
	src := path.Join(clusterPath, env.ConstClusterConfigDir, newCfgPath)
	dst := path.Join(clusterPath, env.ConstClusterConfigDir, oldCfgPath)

	return utils.ForceMove(src, dst)
}

// validateNewConfig validates configuration file. If validation fails,
// a formatted error is returned.
func validateNewConfig(cfg *modelconfig.Config) error {
	err := cfg.Validate()

	if err == nil {
		utils.PrintSuccess("Configuration is valid.")
		return nil
	}

	var errs utils.Errors
	errs = append(errs, fmt.Errorf("Validation of the configuration file has failed.\n"))
	errs = append(errs, fmt.Errorf("Configuration file contains the following (%d) errors:\n", len(errs)))

	for _, e := range err.(validation.ValidationErrors) {
		errs = append(errs, NewValidationError(e.Error(), e.Namespace))
	}

	return errs
}
