package actions

import (
	"cli/cmp"
	"cli/env"
	"cli/utils"
	"fmt"
)

func plan(c Cluster, action env.ApplyAction) ([]*OnChangeEvent, error) {
	if c.OldCfg == nil {
		return nil, nil
	}

	comp := cmp.NewComparator()
	comp.TagName = "opt"

	diff, err := comp.Compare(c.OldCfg, c.NewCfg)

	if err != nil || len(diff.Changes()) == 0 {
		return nil, err
	}

	fmt.Printf("Following changes have been detected:\n\n")
	fmt.Println(diff.ToYamlDiff())

	events := triggerEvents(diff, action)

	var warns utils.Errors
	var errs utils.Errors

	for _, t := range events {
		switch t.cType {
		case WARN:
			warns = append(warns, NewConfigChangeWarning(t.msg, t.triggerPaths...))
		case BLOCK:
			errs = append(errs, NewConfigChangeError(t.msg, t.triggerPaths...))
		}
	}

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
