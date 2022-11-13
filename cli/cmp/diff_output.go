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

// ToYaml returns comparison result in YAML like format.
func (n *DiffNode) ToYaml() string {
	return n.toYaml(0, false, false)
}

// ToYamlDiff returns only differences of the comparison in YAML like format.
func (n *DiffNode) ToYamlDiff() string {
	return n.toYaml(0, false, true)
}

// toYaml recursively creates a string of differences in YAML format.
func (n *DiffNode) toYaml(depth int, tagList, diffOnly bool) string {
	if diffOnly && !n.isDifferent() {
		return ""
	}

	var output []string
	var key string

	isSliceElem := n.isSliceElem()
	indent := fmt.Sprintf("%*s", depth*2, "")

	if !isSliceElem && !n.isRoot() {
		key = n.stringKey() + ": "
	}

	if tagList || (isSliceElem && n.isLeaf()) {
		if len(indent) > 1 {
			indent = indent[:len(indent)-2]
		}
		indent += "- "
	}

	if n.isLeaf() {
		val := n.stringValue()
		return fmt.Sprintf("%s%s%s", indent, key, val)
	}

	if key != "" {
		output = append(output, fmt.Sprintf("%s%s", indent, key))
	}

	if !n.isRoot() {
		depth++
	}

	first := true

	for _, k := range n.sortChildrenKeys() {
		v := n.child(k)

		if diffOnly && !v.isDifferent() {
			continue
		}

		tagList = (isSliceElem && first)
		output = append(output, v.toYaml(depth, tagList, diffOnly))
		first = false
	}

	return strings.Join(output, "\n")
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

// isDifferent returns true if node is either a slice id or it has
// changed.
func (n *DiffNode) isDifferent() bool {
	return n.isSliceId || n.hasChanged()
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
