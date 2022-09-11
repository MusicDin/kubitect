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
	none   = color.New(color.Reset).SprintFunc()
)

type DiffType int

const (
	NIL  DiffType = iota // Unknown
	NONE                 // No change
	CREATE
	DELETE
	MODIFY
)

type DiffNode struct {
	key    string
	parent *DiffNode
	nodes  []*DiffNode
	leafs  []*DiffLeaf
	action DiffType
}

type DiffLeaf struct {
	key    string
	parent *DiffNode
	before interface{}
	after  interface{}
	action DiffType
}

// NewNode returns new node.
func NewNode() *DiffNode {
	node := &DiffNode{
		nodes: make([]*DiffNode, 0),
		leafs: make([]*DiffLeaf, 0),
	}
	return node
}

// AddNode returns a new node that is linked to the current node.
func (n *DiffNode) AddNode(key interface{}) *DiffNode {
	node := NewNode()
	node.key = toString(key)
	node.parent = n
	node.action = NIL

	n.nodes = append(n.nodes, node)

	return node
}

// AddLeaf returns a new leaf that is linked to the current node.
func (n *DiffNode) AddLeaf(a DiffType, key, before, after interface{}) {

	leaf := &DiffLeaf{
		key:    toString(key),
		parent: n,
		action: a,
		after:  after,
		before: before,
	}

	n.leafs = append(n.leafs, leaf)

	n.setActionToRoot(a)
}

// setActionToRoot recursively propagates action across parent
// nodes (to root node).
func (n *DiffNode) setActionToRoot(a DiffType) {

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
// nodes and leafs.
func (n *DiffNode) setActionToLeafs(a DiffType) {

	n.action = a

	for _, l := range n.leafs {
		l.action = a
	}

	for _, v := range n.nodes {
		v.setActionToLeafs(a)
	}
}

func (n *DiffNode) getNode(key string) *DiffNode {
	for _, v := range n.nodes {
		if v.key == key {
			return v
		}
	}
	return nil
}

func (n *DiffNode) getNodeKeys() []string {
	keys := make([]string, len(n.nodes))

	for i, v := range n.nodes {
		keys[i] += v.key
	}

	sort.Strings(keys)
	return keys
}

func (n *DiffNode) getLeaf(key string) *DiffLeaf {
	for _, v := range n.leafs {
		if v.key == key {
			return v
		}
	}
	return nil
}

func (n *DiffNode) getLeafKeys() []string {
	keys := make([]string, len(n.leafs))

	for i, v := range n.leafs {
		keys[i] += v.key
	}

	sort.Strings(keys)
	return keys
}

func (n *DiffNode) getPath() string {
	if n.parent == nil {
		return n.key
	}
	return n.parent.getPath() + "." + n.key
}

func (l *DiffLeaf) getPath() string {
	if l.parent == nil {
		return l.key
	}
	return l.parent.getPath() + "." + l.key
}

// Print outputs differences in YAML format.
func (n *DiffNode) Print() {
	fmt.Print(n.string(-1, false))
}

// string recursively creates an output of differences in YAML format.
func (n *DiffNode) string(depth int, isListElem bool) string {
	var output string

	if len(n.leafs) > 0 {
		keys := n.getLeafKeys()
		isList := isListIndex(keys[0])

		for _, k := range keys {
			v := n.getLeaf(k)

			indent := evalIndent(depth, isListElem)
			isListElem = false

			s := v.toString(isList)
			output += fmt.Sprintf("%s%s\n", indent, s)
		}
	}

	if len(n.nodes) > 0 {
		keys := n.getNodeKeys()
		isList := isListIndex(keys[0])

		for _, k := range keys {
			v := n.getNode(k)

			indent := evalIndent(depth, isListElem)
			isListElem = false

			if !isList && depth >= 0 {
				s := v.toString(isList)
				output += fmt.Sprintf("%s%s\n", indent, s)
			}
			output += v.string(depth+1, isList)
		}
	}

	return output
}

// toString returns node's key as a string.
func (n *DiffNode) toString(isList bool) string {
	if isList {
		return ""
	}

	switch n.action {
	case CREATE:
		return green(fmt.Sprintf("%s:", n.key))
	case DELETE:
		return red(fmt.Sprintf("%s:", n.key))
	default:
		return fmt.Sprintf("%s:", n.key)
	}
}

// toString returns leaf's key and value as a string.
func (l *DiffLeaf) toString(isList bool) string {
	var key string

	if !isList {
		key = l.key + ": "
	}

	switch l.action {
	case CREATE:
		return green(fmt.Sprintf("%s%v", key, l.after))
	case DELETE:
		return red(fmt.Sprintf("%s%v", key, l.before))
	case MODIFY:
		return fmt.Sprintf("%s%v", key, yellow(fmt.Sprintf("%v -> %v", l.before, l.after)))
	default:
		return fmt.Sprintf("%s%v", key, l.before)
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
func evalIndent(d int, tagList bool) string {
	var indent string

	for i := 0; i < d*2; i++ {
		indent += " "
	}

	if tagList {
		indent = indent[:len(indent)-2]
		indent += "- "
	}

	return indent
}
