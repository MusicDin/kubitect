package cmp

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
)

var (
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	green  = color.New(color.FgHiGreen).SprintFunc()
)

// ToYaml returns the result of comparison in YAML format.
func (n *DiffNode) ToYaml(diffOnly bool) string {
	return n.toYaml(0, false, diffOnly)
}

// ToJson returns the result of comparison in JSON format.
func (n *DiffNode) ToJson(diffOnly bool) string {
	return n.toJson(0, diffOnly)
}

// toYaml recursively creates a string of differences in YAML format.
func (n *DiffNode) toYaml(depth int, tagList, diffOnly bool) string {
	var output string

	if diffOnly && !n.hasChanged() {
		return ""
	}

	key := n.stringKey() + ": "
	indent := fmt.Sprintf("%*s", depth*2, "")
	isListIndex := isSliceKey(n.key)

	if isListIndex || n.isRoot() {
		key = ""
	}

	if tagList || (isListIndex && n.isLeaf()) {
		if len(indent) > 1 {
			indent = indent[:len(indent)-2]
		}
		indent += "- "
	}

	if n.isLeaf() {
		val := n.stringValue()
		return fmt.Sprintf("%s%s%s\n", indent, key, val)
	}

	if len(key) > 0 {
		output += fmt.Sprintf("%s%s\n", indent, key)
	}

	if !n.isRoot() {
		depth++
	}

	for i, k := range n.sortChildrenKeys() {
		v := n.getChild(k)
		tagList = (isListIndex && i == 0)
		output += v.toYaml(depth, tagList, diffOnly)
	}

	return output
}

// toJson recursively creates a string of differences in JSON format.
func (n *DiffNode) toJson(depth int, diffOnly bool) string {
	var output string

	if diffOnly && !n.hasChanged() {
		return ""
	}

	key := n.key + ": "
	indent := fmt.Sprintf("%*s", depth*2, "")

	if isSliceKey(n.key) || n.isRoot() {
		key = ""
	}

	if n.isLeaf() {
		value := n.stringValue()
		return fmt.Sprintf("%s%s%s,\n", indent, key, value)
	}

	keys := n.sortChildrenKeys()
	isList := isSliceKey(keys[0])

	annoA := "{"
	annoB := "}"

	if isList {
		annoA = "["
		annoB = "]"
	}

	output += fmt.Sprintf("%s%s%s\n", indent, key, annoA)

	for _, k := range keys {
		v := n.getChild(k)
		output += v.toJson(depth+1, diffOnly)
	}

	output += fmt.Sprintf("%s%s,\n", indent, annoB)

	return output
}

// sortChildrenKeys returns children keys sorted alphabetically.
func (n *DiffNode) sortChildrenKeys() []string {
	keys := make([]string, len(n.children))

	for i, v := range n.children {
		keys[i] += v.key
	}

	sort.Strings(keys)
	return keys
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
func (n *DiffNode) stringValue() string {
	bv := formatValue(n.before)
	av := formatValue(n.after)

	switch n.action {
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

// formatValue formats output value based on its type.
func formatValue(v interface{}) string {
	switch v.(type) {
	case string:
		return fmt.Sprintf("\"%v\"", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
