package cmp

import (
	"reflect"
	"strings"
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
	parent    *DiffNode
	children  []*DiffNode
	action    ActionType
	before    interface{}
	after     interface{}
	typ       reflect.Type
	kind      reflect.Kind
}

func NewEmptyNode() *DiffNode {
	node := &DiffNode{
		children: make([]*DiffNode, 0),
		action:   UNKNOWN,
	}
	return node
}

func NewNode(before, after interface{}) *DiffNode {
	node := NewEmptyNode()
	node.before = before
	node.after = after

	if before == nil && after == nil {
		node.typ = reflect.TypeOf(nil)
		node.kind = reflect.Invalid
	} else if before == nil {
		node.typ = reflect.TypeOf(after)
		node.kind = node.typ.Kind()
	} else {
		node.typ = reflect.TypeOf(before)
		node.kind = node.typ.Kind()
	}

	return node
}

func NewLeaf(a ActionType, before, after interface{}) *DiffNode {
	node := NewNode(before, after)
	node.action = a

	return node
}

func (n *DiffNode) addChild(c *DiffNode, key, structKey interface{}) {
	c.parent = n
	c.key = toString(key)
	c.structKey = toString(structKey)

	n.children = append(n.children, c)

	n.setAction(c.action)
}

// setAction sets node action with consideration to the current action.
func (n *DiffNode) setAction(a ActionType) {
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

// getChild returns a child node with a matching key and nil otherwise.
func (n *DiffNode) getChild(key interface{}) *DiffNode {
	for _, v := range n.children {
		if v.key == key {
			return v
		}
	}
	return nil
}

// path returns node's path as a slice of strings.
func (n *DiffNode) path() []string {
	if n.parent == nil || n.parent.isRoot() {
		return []string{n.key}
	}

	return append(n.parent.path(), n.key)
}

// exactPath returns node's path as a string with each section being
// separated with a dot.
func (n *DiffNode) exactPath() string {
	return strings.Join(n.path(), ".")
}

// genericPath returns the path as a string with all slice keys replaced
// by an asterisk (*).
func (n *DiffNode) genericPath() string {
	path := make([]string, 0)

	for _, s := range n.path() {
		if isSliceKey(s) {
			path = append(path, "[*]")
		} else {
			path = append(path, s)
		}
	}

	return strings.Join(path, ".")
}

// isRoot returns true if node's key is empty.
func (n *DiffNode) isRoot() bool {
	return n.key == ""
}

// isLeaf returns true if node has no children.
func (n *DiffNode) isLeaf() bool {
	return len(n.children) == 0
}

// hasChanged returns true if node's action indicates a change within the
// node or any of its children.
func (n *DiffNode) hasChanged() bool {
	return !(n.action == NONE || n.action == UNKNOWN)
}
