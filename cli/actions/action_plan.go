package actions

import (
	"cli/cmp"
	"cli/config/modelconfig"
	"cli/env"
	"cli/utils"
	"fmt"
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
