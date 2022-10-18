package cmp

type ActionType string

const (
	UNKNOWN ActionType = ""
	NONE    ActionType = "none"
	CREATE  ActionType = "create"
	DELETE  ActionType = "delete"
	MODIFY  ActionType = "modify"
)

type DiffNode struct {
	key      string
	path     []string
	parent   *DiffNode
	children []*DiffNode
	action   ActionType
	before   interface{}
	after    interface{}
	isId     bool
}

// NewNode returns new empty node.
func NewNode() *DiffNode {
	node := &DiffNode{
		children: make([]*DiffNode, 0),
		action:   UNKNOWN,
	}
	return node
}

// addNode returns a new node that is linked to the current node.
func (n *DiffNode) addNode(key interface{}) *DiffNode {
	for _, c := range n.children {
		if c.key == key {
			return c
		}
	}

	var node *DiffNode

	node = NewNode()
	node.key = toString(key)
	node.parent = n

	if !node.isRoot() {
		path := make([]string, len(n.path))
		copy(path, n.path)
		node.path = append(path, node.key)
	}

	n.children = append(n.children, node)

	return node
}

// addLeaf returns a new leaf that is linked to the current node.
func (n *DiffNode) addLeaf(a ActionType, key, before, after interface{}) {
	node := n.addNode(key)
	node.action = a
	node.before = before
	node.after = after
	node.isId = (fromSliceKey(n.key) == before)

	n.setActionToRoot(a)
}

// setActionToRoot recursively propagates action across parent
// nodes (to the root node).
func (n *DiffNode) setActionToRoot(a ActionType) {
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

	if n.parent != nil {
		n.parent.setActionToRoot(a)
	}
}

// setActionToLeafs recursively propagates action across all
// children nodes.
func (n *DiffNode) setActionToLeafs(a ActionType) {
	n.action = a

	for _, v := range n.children {
		v.setActionToLeafs(a)
	}
}

// getChild returns a child node with a matching key and nil
// otherwise.
func (n *DiffNode) getChild(key interface{}) *DiffNode {
	for _, v := range n.children {
		if v.key == key {
			return v
		}
	}
	return nil
}

// isRoot returns true if node's key matches the root key.
func (n *DiffNode) isRoot() bool {
	return n.key == ROOT_KEY
}

// isLeaf returns true if node has no children.
func (n *DiffNode) isLeaf() bool {
	return len(n.children) == 0
}

// hasChanged returns true if node's action is not NIL or NONE.
func (n *DiffNode) hasChanged() bool {
	return !(n.action == NONE || n.action == UNKNOWN)
}
