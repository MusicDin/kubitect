package cluster

import (
	"cli/cluster/event"
	"cli/lib/cmp"
	"cli/ui"
	"fmt"
)

func (c *Cluster) plan(action ApplyAction) (event.Events, error) {
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

	events := event.TriggerEvents(diff, action.events())
	blocking := events.Blocking()

	if len(blocking) > 0 {
		ui.PrintBlockE(blocking.Errors()...)
		return nil, fmt.Errorf("Aborted. Configuration file contains errors.")
	}

	warnings := events.Warns()

	if len(warnings) > 0 {
		ui.PrintBlockE(warnings.Errors()...)
		fmt.Println("Above warnings indicate potentially destructive actions.")
	}

	return events, ui.Ask()
}

// events returns events of the corresponding action.
func (a ApplyAction) events() event.Events {
	switch a {
	case CREATE:
		return event.ModifyEvents
	case SCALE:
		return event.ScaleEvents
	case UPGRADE:
		return event.UpgradeEvents
	default:
		return nil
	}
}
