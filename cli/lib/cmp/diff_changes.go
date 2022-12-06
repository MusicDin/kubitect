package cmp

import (
	"fmt"
	"reflect"
	"strings"
)

type Change struct {
	Path        string
	StructPath  string
	GenericPath string
	Before      interface{}
	After       interface{}
	Type        reflect.Type
	Kind        reflect.Kind
	Action      ActionType
}

type Changes []Change

// Changes returns list of changes extracted from the comparison node.
func (n *DiffNode) Changes() Changes {
	changes := make([]Change, 0)

	if n.isLeaf() && n.hasChanged() {
		changes = append(changes, n.toChange())
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

// toChange converts node into a change (by stripping
// away parent and children references)
func (n DiffNode) toChange() Change {
	return Change{
		Path:        n.path(),
		StructPath:  n.structPath(),
		GenericPath: n.genericPath(),
		Type:        n.typ,
		Kind:        n.kind,
		Before:      n.before,
		After:       n.after,
		Action:      n.action,
	}
}

// String returns change as a string.
func (c Change) String() string {
	if c.Kind == reflect.Struct {
		return fmt.Sprintf("(%s) %s", c.Action, c.Path)
	}
	return fmt.Sprintf("(%s) %s: %v -> %v", c.Action, c.Path, c.Before, c.After)
}

// String returns all changes as a string.
func (cs Changes) String() string {
	var out []string

	for _, c := range cs {
		out = append(out, c.String())
	}

	return strings.Join(out, "\n")
}
