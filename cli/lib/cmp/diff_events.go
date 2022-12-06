package cmp

import "strings"

type Event interface {
	Action() ActionType // Affected action
	Paths() []string    // Affected paths
}

type TriggerEvent interface {
	Event
	Trigger(Change)
}

type TriggerFunc[E Event] func(E, Change)

// TriggerEvents calls Trigger function of each detected event.
func TriggerEvents[E TriggerEvent](n *DiffNode, events []E) {
	for _, c := range n.children {
		TriggerEvents(c, events)
	}

	if n.isRoot() || !n.hasChanged() {
		return
	}

	for _, e := range events {
		if matches(n, e) {
			e.Trigger(n.toChange())
			return
		}
	}
}

// TriggerEventsF calls trigger function for each detected event.
func TriggerEventsF[E Event](n *DiffNode, events []E, trigger TriggerFunc[E]) {
	for _, c := range n.children {
		TriggerEventsF(c, events, trigger)
	}

	if n.isRoot() || !n.hasChanged() {
		return
	}

	for _, e := range events {
		if matches(n, e) {
			trigger(e, n.toChange())
			return
		}
	}
}

// MatchingChanges returns changes that match (trigger) given events.
func MatchingChanges[E Event](n *DiffNode, events []E) Changes {
	m, _ := categorizeChanges(n, events, false)
	return m
}

// ConflictingChanges returns changes that are either excluded from
// all given events or have a conflicting actions (but their paths match).
func ConflictingChanges[E Event](n *DiffNode, events []E) Changes {
	_, c := categorizeChanges(n, events, false)
	return c
}

// categorizeChanges categorizes changes into matching and conflicting
// (non-matching) events.
func categorizeChanges[E Event](n *DiffNode, events []E, mismatch bool) (Changes, Changes) {
	mat := make(Changes, 0)
	con := make(Changes, 0)

	if !n.isRoot() {
		if !n.hasChanged() {
			return mat, con
		}

		if n.isLeaf() {
			if mismatch || excludes(n, events) {
				con = append(con, n.toChange())
			} else {
				mat = append(mat, n.toChange())
			}

			return mat, con
		}

		if conflicts(n, events) {
			mismatch = true
		}
	}

	for _, c := range n.children {
		cm, cc := categorizeChanges(c, events, mismatch)
		mat = append(mat, cm...)
		con = append(con, cc...)
	}

	return mat, con
}

// matches returns true if the path and action of the node match the
// path and action of the event.
func matches[E Event](n *DiffNode, e E) bool {
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
func conflicts[E Event](n *DiffNode, events []E) bool {
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

// excludes returns true if there is no event that both matches the action
// of the change and has a path that is a prefix of the path of the change.
func excludes[E Event](n *DiffNode, events []E) bool {
	for _, e := range events {
		a := e.Action()

		if a != UNKNOWN && a != n.action {
			continue
		}

		for _, p := range e.Paths() {
			if strings.HasPrefix(n.genericPath(), p) {
				return false
			}
		}
	}

	return true
}
