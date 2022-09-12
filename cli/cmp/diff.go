package cmp

type ActionType int

const (
	NIL  ActionType = iota // Unknown
	NONE                   // No change
	CREATE
	DELETE
	MODIFY
)

type DiffNode struct {
	key      string
	parent   *DiffNode
	children []*DiffNode
	action   ActionType
	before   interface{}
	after    interface{}
}

// NewNode returns new node.
func NewNode() *DiffNode {
	node := &DiffNode{
		children: make([]*DiffNode, 0),
	}
	return node
}

// addNode returns a new node that is linked to the current node.
func (n *DiffNode) addNode(key interface{}) *DiffNode {
	node := NewNode()
	node.key = toString(key)
	node.parent = n
	node.action = NIL

	n.children = append(n.children, node)

	return node
}

// addLeaf returns a new leaf that is linked to the current node.
func (n *DiffNode) addLeaf(a ActionType, key, before, after interface{}) {
	node := NewNode()
	node.key = toString(key)
	node.parent = n
	node.action = a
	node.before = before
	node.after = after

	n.children = append(n.children, node)

	n.setActionToRoot(a)
}

// setActionToRoot recursively propagates action across parent
// nodes (to the root node).
func (n *DiffNode) setActionToRoot(a ActionType) {
	switch n.action {
	case CREATE:
		if a == DELETE {
			n.action = MODIFY
		} else {
			n.action = a
		}
	case DELETE:
		if a == CREATE {
			n.action = MODIFY
		} else {
			n.action = a
		}
	case NONE:
		if a == MODIFY {
			n.action = a
		}
	case NIL:
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

// isLeaf returns true if node has no children.
func (n *DiffNode) isLeaf() bool {
	return len(n.children) == 0
}

// getPath returns all keys to the root node, separated by a dot.
func (n *DiffNode) getPath() string {
	if n.parent == nil {
		return n.key
	}
	return n.parent.getPath() + "." + n.key
}
