package cmp

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
)

var (
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	green  = color.New(color.FgHiGreen).SprintFunc()
	blue   = color.New(color.FgHiBlue).SprintFunc()
	none   = color.New(color.Reset).SprintFunc()
)

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

// AddNode returns a new node that is linked to the current node.
func (n *DiffNode) AddNode(key interface{}) *DiffNode {
	node := NewNode()
	node.key = toString(key)
	node.parent = n
	node.action = NIL

	n.children = append(n.children, node)

	return node
}

// AddLeaf returns a new leaf that is linked to the current node.
func (n *DiffNode) AddLeaf(a ActionType, key, before, after interface{}) {
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

// setActionToRoot recursively propagates action across all children
// nodes.
func (n *DiffNode) setActionToLeafs(a ActionType) {
	n.action = a

	for _, v := range n.children {
		v.setActionToLeafs(a)
	}
}

// ToYaml returns the result of comparison in JSON format.
func (n *DiffNode) ToJson() string {
	return n.toJson(-2, false)
}

// ToYaml returns the result of comparison in YAML format.
func (n *DiffNode) ToYaml() string {
	return n.toYaml(-2, false)
}

// toJson recursively creates a string of differences in YAML format.
func (n *DiffNode) toYaml(depth int, tagList bool) string {
	var output string

	key := n.stringKey() + ": "
	indent := evalIndent(depth)
	isList := isListIndex(n.key)

	if tagList {
		indent = indent[:len(indent)-2]
		indent += "- "
	}

	if n.isLeaf() {
		val := n.stringValue()
		return fmt.Sprintf("%s%s%s\n", indent, key, val)
	}

	if !isList && depth >= 0 {
		output += fmt.Sprintf("%s%s\n", indent, key)
	}

	for i, k := range n.sortChildrenKeys() {
		v := n.getChild(k)
		tagList = isList && i == 0
		output += v.toYaml(depth+1, tagList)
	}

	return output
}

// toJson recursively creates a string of differences in JSON format.
func (n *DiffNode) toJson(depth int, isListElem bool) string {
	var output string

	indent := evalIndent(depth)
	key := n.stringKey() + ": "
	value := n.stringValue()

	// Leaf
	if len(n.children) == 0 {
		return fmt.Sprintf("%s%s%s,\n", indent, key, value)
	}

	keys := n.sortChildrenKeys()
	isList := isListIndex(keys[0])

	if isListElem {
		output += fmt.Sprintf("%s{\n", indent)
	} else if isList {
		output += fmt.Sprintf("%s%s[\n", indent, key)
	} else {
		output += fmt.Sprintf("%s%s{\n", indent, key)
	}

	for _, k := range keys {
		v := n.getChild(k)

		output += v.toJson(depth+1, isList)
	}

	if isList {
		output += fmt.Sprintf("%s],\n", indent)
	} else {
		output += fmt.Sprintf("%s},\n", indent)
	}

	return output
}

// stringKey returns node's key as a string.
func (n *DiffNode) stringKey() string {
	switch n.action {
	case CREATE:
		return green(n.key)
	case DELETE:
		return red(n.key)
	default:
		return n.key
	}
}

// stringValue returns node's value change as a string.
func (l *DiffNode) stringValue() string {
	bv := formatValue(l.before)
	av := formatValue(l.after)

	switch l.action {
	case CREATE:
		return green(av)
	case DELETE:
		return red(bv)
	case MODIFY:
		return yellow(fmt.Sprintf("%v -> %v", bv, av))
	default:
		return bv
	}
}

// getChild returns a child node that matches a key and nil otherwise.
func (n *DiffNode) getChild(key string) *DiffNode {
	for _, v := range n.children {
		if v.key == key {
			return v
		}
	}
	return nil
}

func (n *DiffNode) sortChildrenKeys() []string {
	keys := make([]string, len(n.children))

	for i, v := range n.children {
		keys[i] += v.key
	}

	sort.Strings(keys)
	return keys
}

func (n *DiffNode) isLeaf() bool {
	return len(n.children) == 0
}

func (n *DiffNode) getPath() string {
	if n.parent == nil {
		return n.key
	}
	return n.parent.getPath() + "." + n.key
}

func formatValue(v interface{}) string {
	switch v.(type) {
	case string:
		return fmt.Sprintf("\"%v\"", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// isListIndex checks whether given key represents a list index,
// which means that it starts with "[" and ends with "]".
func isListIndex(k interface{}) bool {
	s := toString(k)
	return strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]")
}

// evalIndent returns an indentation (for yaml) based on the given depth.
// If isList is set to true, it will include yaml list identifier (-).
func evalIndent(depth int) string {
	var indent string

	for i := 0; i < depth*2; i++ {
		indent += " "
	}

	// if tagList {
	// 	indent = indent[:len(indent)-2]
	// 	indent += "- "
	// }

	return indent
}
