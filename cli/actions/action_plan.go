package actions

import (
	"cli/cmp"
	"cli/ui"
	"fmt"
)

func (c *Cluster) plan(action ApplyAction) (Events, error) {
	if c.AppliedConfig == nil {
		return nil, nil
	}

	comp := cmp.NewComparator()
	comp.Tag = "opt"
	comp.ExtraNameTags = []string{"yaml"}
	comp.IgnoreEmptyChanges = true
	comp.PopulateStructNodes = true

	diff, err := comp.Compare(c.AppliedConfig, c.NewConfig)

	if err != nil || len(diff.Changes()) == 0 {
		return nil, err
	}

	fmt.Printf("Following changes have been detected:\n\n")
	fmt.Println(diff.ToYamlDiff())

	events := triggerEvents(diff, action)
	blocking := events.OfType(BLOCK)

	if len(blocking) > 0 {
		ui.PrintBlock(blocking.Errors()...)
		return nil, fmt.Errorf("Aborted. Configuration file contains errors.")
	}

	warnings := events.OfType(WARN)

	if len(warnings) > 0 {
		ui.PrintBlock(warnings.Errors()...)
		fmt.Println("Above warnings indicate potentially destructive actions.")
	}

	return events, ui.Ask()
}
