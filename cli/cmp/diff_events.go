package cmp

import (
	"strings"
)

type ChangeEvent interface {
	GetPaths() []string    // Affected paths
	GetAction() ActionType // Affected action
}

// TriggerEvents returns a list of events triggered by changes.
func TriggerEvents[E ChangeEvent](n *DiffNode, events []E) []E {
	triggered := make([]E, 0)

	if !n.isRoot() {
		if !n.hasChanged() {
			return triggered
		}

		for _, e := range events {
			if triggers(n, e) {
				triggered = append(triggered, e)
			}
		}
	}

	for _, c := range n.children {
		triggered = append(triggered, TriggerEvents(c, events)...)
	}

	return triggered
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

	a := n.action

	if n.action == NONE {
		return false
	}

	p := n.genericPath()
	ea := e.GetAction()

	for _, ep := range e.GetPaths() {
		if ep == p && (ea == a || ea == UNKNOWN) {
			return true
		}
	}

	return false
}

// excludes returns true if the path of the node excludes paths
// from all events.
func excludes[E ChangeEvent](n *DiffNode, events []E) bool {
	for _, e := range events {
		for _, p := range e.GetPaths() {
			if strings.Contains(n.genericPath(), p) {
				return false
			}
		}
	}

	return true
}

// conflicts returns true if node and event paths are the same,
// but their actions do not match (conflict).
func conflicts[E ChangeEvent](n *DiffNode, events []E) bool {
	var matched bool

	for _, e := range events {
		for _, p := range e.GetPaths() {
			if p == n.genericPath() {
				a := e.GetAction()

				if a == UNKNOWN || a == n.action {
					return false
				}

				matched = true
			}
		}
	}

	return matched
}

// GenericPath returns the path as a string with all slice keys
// replaced by an asterisk (*).
func (n *DiffNode) genericPath() string {
	path := make([]string, 0)

	for _, s := range n.path {
		if isSliceKey(s) {
			path = append(path, "[*]")
		} else {
			path = append(path, s)
		}
	}

	return strings.Join(path, ".")
}
