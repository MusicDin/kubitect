package cmp

import "strings"

type ChangeEvent interface {
	Action() ActionType // Affected action
	Paths() []string    // Affected paths
	TriggerPath(string) // Set path of a change that triggered an event
}

// TriggerEvents returns a list of triggered events.
func TriggerEvents[E ChangeEvent](n *DiffNode, events []E) []E {
	triggered := new([]E)
	triggerEvents(n, events, triggered)
	return *triggered
}

// triggerEvents detects triggered events and appends them to the triggered slice.
// Whenever an event is triggered, a TriggerPath method is called with an actual path
// that has triggered an event.
func triggerEvents[E ChangeEvent](n *DiffNode, events []E, triggered *[]E) {
	for _, c := range n.children {
		triggerEvents(c, events, triggered)
	}

	if n.isRoot() || !n.hasChanged() {
		return
	}

	for i, e := range *triggered {
		if n.action == e.Action() && triggers(n, e) {
			(*triggered)[i].TriggerPath(n.exactPath())
			return
		}
	}

	for _, e := range events {
		if triggers(n, e) {
			e.TriggerPath(n.exactPath())
			*triggered = append(*triggered, e)
			return
		}
	}
}

// MatchingChanges returns changes that match (trigger) given events.
func MatchingChanges[E ChangeEvent](n *DiffNode, events []E) Changes {
	m, _ := categorizeChanges(n, events, false)
	return m
}

// NonMatchingChanges returns changes that are either excluded from
// all given events or have a conflicting actions (but their paths match).
func NonMatchingChanges[E ChangeEvent](n *DiffNode, events []E) Changes {
	_, mm := categorizeChanges(n, events, false)
	return mm
}

// categorizeChanges returns two slices of changes. The first slice contains
// changes categorized as matching (those that trigger a specific event) and
// second contains changes that are completely excluded from these events.
func categorizeChanges[E ChangeEvent](n *DiffNode, events []E, mismatch bool) (Changes, Changes) {
	matched := make(Changes, 0)
	mismatched := make(Changes, 0)

	if !n.isRoot() {
		if !n.hasChanged() {
			return matched, mismatched
		}

		if n.isLeaf() {
			if mismatch || excludes(n, events) {
				mismatched = append(mismatched, n.toChange())
			} else {
				matched = append(matched, n.toChange())
			}

			return matched, mismatched
		}

		if conflicts(n, events) {
			mismatch = true
		}
	}

	for _, c := range n.children {
		m, mm := categorizeChanges(c, events, mismatch)
		matched = append(matched, m...)
		mismatched = append(mismatched, mm...)
	}

	return matched, mismatched
}

// triggers returns true if the path and action of the node match the
// path and action of the event.
func triggers[E ChangeEvent](n *DiffNode, e E) bool {
	if n.action == NONE {
		return false
	}

	a := e.Action()
	np := n.genericPath()

	for _, p := range e.Paths() {
		if p == np && (a == UNKNOWN || a == n.action) {
			return true
		}
	}

	return false
}

// conflicts returns true if node and event paths are the same,
// but their actions do not match (conflict).
func conflicts[E ChangeEvent](n *DiffNode, events []E) bool {
	var conflict bool

	for _, e := range events {
		a := e.Action()

		for _, p := range e.Paths() {
			if p == n.genericPath() {
				if a == UNKNOWN || a == n.action {
					return false
				}

				conflict = true
			}
		}
	}

	return conflict
}

// excludes returns true if none of the given event paths is a prefix
// of a node's path.
func excludes[E ChangeEvent](n *DiffNode, events []E) bool {
	for _, e := range events {
		for _, p := range e.Paths() {
			if strings.HasPrefix(n.genericPath(), p) {
				return false
			}
		}
	}

	return true
}
