package executors

import (
	"fmt"

	"github.com/MusicDin/kubitect/pkg/cluster/event"
	"github.com/MusicDin/kubitect/pkg/models/config"
)

// extractRemovedNodes returns removed node instances extracted from the event changes.
func extractRemovedNodes(events []event.Event) ([]config.Instance, error) {
	var nodes []config.Instance
	for _, e := range events {
		if e.Rule.ActionType != event.Action_ScaleDown {
			continue
		}

		node, ok := e.Change.ValueBefore.(config.Instance)
		if ok {
			nodes = append(nodes, node)
			continue
		}

		return nil, fmt.Errorf("%v cannot be scaled", e.Change.ValueType.Name())
	}

	return nodes, nil
}

// extractNewNodes returns new node instances extracted from the event changes.
func extractNewNodes(events []event.Event) ([]config.Instance, error) {
	var nodes []config.Instance
	for _, e := range events {
		if e.Rule.ActionType != event.Action_ScaleUp {
			continue
		}

		node, ok := e.Change.ValueAfter.(config.Instance)
		if ok {
			nodes = append(nodes, node)
			continue
		}

		return nil, fmt.Errorf("%v cannot be scaled", e.Change.ValueType.Name())
	}

	return nodes, nil
}
