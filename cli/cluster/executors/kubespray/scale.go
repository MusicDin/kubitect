package kubespray

import (
	"cli/cluster/event"
	"cli/config/modelconfig"
	"fmt"
)

// ScaleUp adds new nodes to the cluster.
func (e *KubesprayExecutor) ScaleUp(events event.Events) error {
	events = events.OfType(event.SCALE_UP)

	if len(events) == 0 {
		return nil
	}

	if err := e.KubitectInit(TAG_KUBESPRAY, TAG_GEN_NODES); err != nil {
		return err
	}

	if err := e.HAProxy(); err != nil {
		return err
	}

	return e.KubesprayScale()
}

// scaleDown gracefully removes nodes from the cluster.
func (e *KubesprayExecutor) ScaleDown(events event.Events) error {
	events = events.OfType(event.SCALE_DOWN)

	if len(events) == 0 {
		return nil
	}

	rmNodes, err := extractRemovedNodes(events)

	if err != nil {
		return err
	}

	if len(rmNodes) == 0 {
		return nil
	}

	var names []string

	for _, n := range rmNodes {
		name := fmt.Sprintf("%s-%s-%s", e.ClusterName, n.GetTypeName(), *n.GetID())
		names = append(names, name)
	}

	if err := e.KubitectInit(TAG_KUBESPRAY); err != nil {
		return err
	}

	return e.KubesprayRemoveNodes(names)
}

// extractRemovedNodes returns node instances from the event changes.
func extractRemovedNodes(events event.Events) ([]modelconfig.Instance, error) {
	var nodes []modelconfig.Instance

	for _, e := range events {
		for _, ch := range e.Changes() {
			if i, ok := ch.Before.(modelconfig.Instance); ok {
				nodes = append(nodes, i)
				continue
			}

			return nil, fmt.Errorf("%v cannot be scaled", ch.Type.Name())
		}
	}

	return nodes, nil
}
