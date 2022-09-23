package cmp

import (
	"strings"
)

type ChangeEvent interface {
	GetPath() string       // Affected path
	GetAction() ActionType // Affected action
}

type ChangeEvents[C ChangeEvent] []C

// TriggerEvents returns a list of events triggered by changes.
func TriggerEvents[C ChangeEvent](n *DiffNode, events []C) []C {
	triggered := make([]C, 0)

	if !n.isRoot() {
		if !n.hasChanged() {
			return triggered
		}

		for _, e := range events {
			if n.triggers(e) {
				triggered = append(triggered, e)
			}
		}
	}

	for _, c := range n.children {
		triggered = append(triggered, TriggerEvents(c, events)...)
	}

	return triggered
}

// Triggers returns true if the path and action of the node match the
// path and action of the event.
func (n *DiffNode) triggers(e ChangeEvent) bool {
	na := n.action
	np := n.genericPath()
	ea := e.GetAction()
	ep := e.GetPath()

	return np == ep && (ea == na || ea == UNKNOWN)
}

// GenericPath returns the path as a string with all slice keys
// replaced by an asterisk (*).
func (n *DiffNode) genericPath() string {
	path := make([]string, 0)

	for _, s := range n.path {
		if isSliceKey(s) {
			path = append(path, "*")
		} else {
			path = append(path, s)
		}
	}

	return strings.Join(path, ".")
}
