package cmp

import (
	"fmt"
	"strings"
)

type Change struct {
	Path   []string
	Before interface{}
	After  interface{}
	Action ActionType
}

type Changes []Change

// Changes returns list of changes extracted from the comparison node.
func (n *DiffNode) Changes() Changes {
	changes := make([]Change, 0)

	if n.isLeaf() && n.hasChanged() {
		changes = append(changes, Change{
			Path:   n.path,
			Before: n.before,
			After:  n.after,
			Action: n.action,
		})
		return changes
	}

	for _, c := range n.children {
		changes = append(changes, c.Changes()...)
	}

	if len(changes) == 0 {
		return nil
	}

	return changes
}

// String returns change as a string.
func (c Change) String() string {
	return fmt.Sprintf("(%s) %s: %v -> %v",
		c.Action,
		strings.Join(c.Path, "."),
		c.Before,
		c.After,
	)
}

// String returns all changes as a string.
func (cs Changes) String() string {
	var out []string

	for _, c := range cs {
		out = append(out, c.String())
	}

	return strings.Join(out, "\n")
}
