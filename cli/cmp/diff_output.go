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
)

// ToYaml returns the result of comparison in JSON format.
func (n *DiffNode) ToJson() string {
	return n.toJson(0)
}

// ToYaml returns the result of comparison in YAML format.
func (n *DiffNode) ToYaml() string {
	return n.toYaml(-1, false)
}

// toJson recursively creates a string of differences in YAML format.
func (n *DiffNode) toYaml(depth int, tagList bool) string {
	var output string

	key := n.stringKey() + ": "
	indent := fmt.Sprintf("%*s", depth*2, "")
	isListIndex := isListIndex(n.key)

	// Skip "[index]" nodes
	if isListIndex {
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

	if !isListIndex && depth >= 0 {
		output += fmt.Sprintf("%s%s\n", indent, key)
	}

	for i, k := range n.sortChildrenKeys() {
		v := n.getChild(k)
		tagList = (isListIndex && i == 0)
		output += v.toYaml(depth+1, tagList)
	}

	return output
}

// toJson recursively creates a string of differences in JSON format.
func (n *DiffNode) toJson(depth int) string {
	var output string

	key := n.key + ": "
	indent := fmt.Sprintf("%*s", depth*2, "")

	// Skip "[index]" nodes
	if isListIndex(n.key) || depth == 0 {
		key = ""
	}

	if n.isLeaf() {
		value := n.stringValue()
		return fmt.Sprintf("%s%s%s,\n", indent, key, value)
	}

	keys := n.sortChildrenKeys()
	isList := isListIndex(keys[0])

	annoA := "{"
	annoB := "}"

	if isList {
		annoA = "["
		annoB = "]"
	}

	output += fmt.Sprintf("%s%s%s\n", indent, key, annoA)

	for _, k := range keys {
		v := n.getChild(k)
		output += v.toJson(depth + 1)
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

// formatValue formats output value based on its type.
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
