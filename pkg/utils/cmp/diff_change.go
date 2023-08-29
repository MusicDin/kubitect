package cmp

import (
	"fmt"
	"reflect"
)

// ChangeType represents the type of change detected between two values.
type ChangeType string

const (
	// Any indicates that the change type is unknown or unspecified.
	// In context of event rules, this change type is used to match any
	// other change type.
	Any ChangeType = ""

	// None indicates that no change has occurred between the two values.
	None ChangeType = "none"

	// Create indicates that a new field or value has been added to
	// the structure.
	Create ChangeType = "create"

	// Delete indicates that a field or value has been removed from
	// the structure.
	Delete ChangeType = "delete"

	// Modify indicates that an existing field or value the has been
	// altered.
	Modify ChangeType = "modify"
)

type Change struct {
	Type ChangeType

	// Change paths.
	Path       string
	StructPath string

	// Metadata of compared values.
	ValueType   reflect.Type `cmp:"-"`
	ValueKind   reflect.Kind `cmp:"-"`
	ValueBefore any
	ValueAfter  any
}

// String returns change as a string.
func (c Change) String() string {
	if c.ValueKind == reflect.Struct {
		return fmt.Sprintf("(%s) %s", c.Type, c.Path)
	}

	return fmt.Sprintf("(%s) %s: %v -> %v", c.Type, c.Path, c.ValueBefore, c.ValueAfter)
}

// ToChange returns the diff node as Change.
func (n *DiffNode) ToChange() Change {
	return Change{
		Type:        n.changeType,
		Path:        n.Path(),
		StructPath:  n.structPath(),
		ValueType:   n.valueType,
		ValueKind:   n.valueKind,
		ValueBefore: n.valueBefore,
		ValueAfter:  n.valueAfter,
	}
}

// leafChanges extracts changes from the leaf nodes of the tree, ignoring
// changes in intermediary nodes.
func (n *DiffNode) leafChanges(changes []Change) []Change {
	if !n.HasChanged() {
		return changes
	}

	if n.IsLeaf() {
		changes = append(changes, n.ToChange())
		return changes
	}

	for _, c := range n.children {
		changes = c.leafChanges(changes)
	}

	return changes
}

// distinctChanges extracts changes from the result tree, filtering out
// propagated changes. The extraction behavior is determined by the change
// type:
//   - Create/Delete: The change is returned without further traversal into
//     child nodes, as all descendants have the same change type.
//   - Modify: If the node is a leaf, the change is returned; otherwise, the
//     function descends further into child nodes.
func (n *DiffNode) distinctChanges(changes []Change) []Change {
	if !n.HasChanged() {
		return changes
	}

	if n.IsLeaf() || (n.changeType == Create || n.changeType == Delete) {
		changes = append(changes, n.ToChange())
		return changes
	}

	for _, c := range n.children {
		changes = c.distinctChanges(changes)
	}

	return changes
}
