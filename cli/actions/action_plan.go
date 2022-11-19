package actions

import (
	"cli/cmp"
	"cli/env"
	"cli/utils"
	"fmt"
)

func plan(c Cluster, action env.ApplyAction) (Events, error) {
	if c.OldCfg == nil {
		return nil, nil
	}

	comp := cmp.NewComparator()
	comp.Tag = "opt"
	comp.ExtraNameTags = []string{"yaml"}
	comp.IgnoreEmptyChanges = true
	comp.PopulateStructNodes = true

	diff, err := comp.Compare(c.OldCfg, c.NewCfg)

	if err != nil || len(diff.Changes()) == 0 {
		return nil, err
	}

	fmt.Printf("Following changes have been detected:\n\n")
	fmt.Println(diff.ToYamlDiff())

	events := triggerEvents(diff, action)
	blocking := events.OfType(BLOCK)

	if len(blocking) > 0 {
		return nil, blocking.Errors()
	}

	warnings := events.OfType(WARN)

	if len(warnings) > 0 {
		fmt.Println(warnings.Errors())
		fmt.Println("Above warnings indicate potentially destructive actions.")
	}

	err = utils.AskUserConfirmation()

	return events, err
}
