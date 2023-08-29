package cmp

import (
	"fmt"
	"reflect"
	"sort"
)

// DiffNode forms the comparison tree, capturing the differences between
// corresponding values in two data structures. A single DiffNode can act as
// the root of the tree or as any subsequent node.
type DiffNode struct {
	changeType ChangeType
	parent     *DiffNode
	children   []*DiffNode

	// Fields related to the data being compared by this node.
	valueType   reflect.Type
	valueKind   reflect.Kind
	valueBefore any
	valueAfter  any

	// pathKey and pathStructKey are used to form a change path.
	pathKey       string
	pathStructKey string

	// isSliceId indicates if the node represents an element in a slice
	// with an 'id' tag option.
	isSliceId bool
}

func NewEmptyNode(t reflect.Type, k reflect.Kind) *DiffNode {
	return &DiffNode{
		children:   []*DiffNode{},
		changeType: Any,
		valueType:  t,
		valueKind:  k,
	}
}

func (c *Comparator) newNode(before any, after any) *DiffNode {
	var t reflect.Type
	var k reflect.Kind

	if before == nil && after == nil {
		t = reflect.TypeOf(nil)
		k = reflect.Invalid
	} else if before == nil {
		t = reflect.TypeOf(after)
		k = t.Kind()
	} else {
		t = reflect.TypeOf(before)
		k = t.Kind()
	}

	node := NewEmptyNode(t, k)
	node.valueBefore = before
	node.valueAfter = after
	return node
}

func (c *Comparator) newLeaf(a ChangeType, before any, after any) *DiffNode {
	if c.options.IgnoreEmptyChanges && before == nil && after == nil {
		return nil
	}

	node := c.newNode(before, after)
	node.changeType = a
	return node
}

func (n *DiffNode) addChild(c *DiffNode, key string, structKey string) {
	if c == nil {
		return
	}

	c.parent = n
	c.pathKey = key
	c.pathStructKey = structKey

	n.children = append(n.children, c)
	n.setChangeType(c.changeType)
}

// setChangeType sets node's change type based on the provided and current
// change type. If the current change type is Any, the node adopts the provided
// value. If the current and provided types differ, the change type is set to
// Modify.
func (n *DiffNode) setChangeType(t ChangeType) {
	if t == Any {
		return
	}

	if n.changeType == Any {
		n.changeType = t
	}

	// If the node's current changeType isn't 'Any', it's been set before.
	// Therefore, a differing new changeType implies a modification.
	if n.changeType != t {
		n.changeType = Modify
	}
}

// setChangeTypeOfChildren propagates change type across all children nodes.
func (n *DiffNode) setChangeTypeOfChildren(t ChangeType) {
	n.changeType = t

	for _, c := range n.children {
		c.setChangeTypeOfChildren(t)
	}
}

// ChangeType returns type of the node's change.
func (n *DiffNode) ChangeType() ChangeType {
	return n.changeType
}

// Parent returns node's parent or nil if node is root.
func (n *DiffNode) Parent() *DiffNode {
	return n.parent
}

// ParentByPath searches upwards from the current node to find a node matching
// the given path. Returns the matching node or nil if none is found.
func (n *DiffNode) ParentByPath(path string) *DiffNode {
	if n.Path() == path {
		return n
	}

	if n.parent != nil {
		return n.parent.ParentByPath(path)
	}

	return nil
}

// Child returns a Child node with a matching key and nil otherwise.
func (n *DiffNode) Child(key any) *DiffNode {
	for _, c := range n.children {
		if c.pathKey == key {
			return c
		}
	}

	return nil
}

// ChildrenKeys returns keys of node's children.
func (n *DiffNode) Children() []*DiffNode {
	return n.children
}

// ChildrenKeys returns keys of node's children.
func (n *DiffNode) ChildrenKeys() []string {
	keys := make([]string, len(n.children))
	for i, v := range n.children {
		keys[i] += v.pathKey
	}

	return keys
}

// ChildrenKeysSorted returns alphabetically sorted keys of node's children.
func (n *DiffNode) ChildrenKeysSorted() []string {
	keys := n.ChildrenKeys()
	sort.Strings(keys)
	return keys
}

// Path returns node's path as a string with each section being
// separated with a dot.
func (n *DiffNode) Path() string {
	if n.parent == nil || n.parent.IsRoot() {
		return n.pathKey
	}

	return fmt.Sprintf("%s.%s", n.parent.Path(), n.pathKey)
}

// structPath returns node's path as a string with each section being
// separated with a dot. Path is constructed from structKeys.
func (n *DiffNode) structPath() string {
	if n.parent == nil || n.parent.IsRoot() {
		return n.pathStructKey
	}

	return fmt.Sprintf("%s.%s", n.parent.structPath(), n.pathStructKey)
}

// IsRoot returns true if node has no parent.
func (n *DiffNode) IsRoot() bool {
	return n.parent == nil
}

// IsLeaf returns true if node has no children.
func (n *DiffNode) IsLeaf() bool {
	return len(n.children) == 0
}

// HasChanged returns true if node's change type indicates a change within the
// node itself or any of its children.
func (n *DiffNode) HasChanged() bool {
	return !(n.changeType == None || n.changeType == Any)
}

// isSliceIndex returns true if node represents a slice index. Node is
// considered a slice index if its parent is either a slice or an array.
func (n *DiffNode) isSliceIndex() bool {
	return n.parent != nil && (n.parent.valueKind == reflect.Slice || n.parent.valueKind == reflect.Array)
}
