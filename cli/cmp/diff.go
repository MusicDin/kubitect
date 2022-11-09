package cmp

import (
	"fmt"
	"reflect"
)

type ActionType string

const (
	UNKNOWN ActionType = ""
	NONE    ActionType = "none"
	CREATE  ActionType = "create"
	DELETE  ActionType = "delete"
	MODIFY  ActionType = "modify"
)

type DiffNode struct {
	key       string
	structKey string
	typ       reflect.Type
	kind      reflect.Kind
	parent    *DiffNode
	children  []*DiffNode
	action    ActionType
	before    interface{}
	after     interface{}

	isSliceId bool
}

func NewEmptyNode(t reflect.Type, k reflect.Kind) *DiffNode {
	node := &DiffNode{
		children: make([]*DiffNode, 0),
		action:   UNKNOWN,
		typ:      t,
		kind:     k,
	}
	return node
}

func NewNode(before, after interface{}) *DiffNode {
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
	node.before = before
	node.after = after

	return node
}

func NewLeaf(a ActionType, before, after interface{}) *DiffNode {
	node := NewNode(before, after)
	node.action = a

	return node
}

func (n *DiffNode) addChild(c *DiffNode, key, structKey interface{}) {
	if c == nil {
		return
	}

	c.parent = n
	c.key = toString(key)
	c.structKey = toString(structKey)

	n.children = append(n.children, c)

	n.setAction(c.action)
}

// setAction sets node action with consideration to the current action.
func (n *DiffNode) setAction(a ActionType) {
	if a == UNKNOWN {
		return
	}

	switch n.action {
	case CREATE:
		if a != CREATE {
			n.action = MODIFY
		}
	case DELETE:
		if a != DELETE {
			n.action = MODIFY
		}
	case NONE:
		if a != NONE {
			n.action = MODIFY
		}
	case UNKNOWN:
		n.action = a
	}
}

// setActionToLeafs recursively propagates action across all children nodes.
func (n *DiffNode) setActionToLeafs(a ActionType) {
	n.action = a

	for _, v := range n.children {
		v.setActionToLeafs(a)
	}
}

// child returns a child node with a matching key and nil otherwise.
func (n *DiffNode) child(key interface{}) *DiffNode {
	for _, v := range n.children {
		if v.key == key {
			return v
		}
	}
	return nil
}

// exactPath returns node's path as a string with each section being
// separated with a dot.
func (n *DiffNode) exactPath() string {
	if n.parent == nil || n.parent.isRoot() {
		return n.key
	}

	return fmt.Sprintf("%s.%s", n.parent.exactPath(), n.key)
}

// structPath returns node's path as a string with each section being
// separated with a dot. Path is constructed from structKeys.
func (n *DiffNode) structPath() string {
	if n.parent == nil || n.parent.isRoot() {
		return n.structKey
	}

	return fmt.Sprintf("%s.%s", n.parent.exactPath(), n.structKey)
}

// genericPath returns the path as a string with all slice keys replaced
// by an asterisk (*).
func (n *DiffNode) genericPath() string {
	key := n.structKey

	if n.isSliceElem() {
		key = "*"
	}

	if n.parent == nil || n.parent.isRoot() {
		return key
	}

	return fmt.Sprintf("%s.%s", n.parent.genericPath(), key)
}

// isRoot returns true if node's key is empty.
func (n *DiffNode) isRoot() bool {
	return n.key == ""
}

// isLeaf returns true if node has no children.
func (n *DiffNode) isLeaf() bool {
	return len(n.children) == 0
}

// isSlice returns true if node's kind is either slice or array.
func (n *DiffNode) isSlice() bool {
	return (n.kind == reflect.Slice || n.kind == reflect.Array)
}

// isSliceElem returns true if node's parent is either a slice or an array.
func (n *DiffNode) isSliceElem() bool {
	return (n.parent != nil && n.parent.isSlice())
}

// isSliceElem returns true if node's parent is either a slice or an array.
func (n *DiffNode) isEmpty() bool {
	return (n.isLeaf() && n.before == nil && n.after == nil)
}

// hasChanged returns true if node's action indicates a change within the
// node or any of its children.
func (n *DiffNode) hasChanged() bool {
	return !(n.action == NONE || n.action == UNKNOWN)
}
