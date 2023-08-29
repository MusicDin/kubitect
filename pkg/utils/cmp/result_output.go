package cmp

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	green  = color.New(color.FgHiGreen).SprintFunc()
)

// FormatOptions contains configurations that control the result output.
type FormatOptions struct {
	// ShowDiffOnly omits unchanged fields.
	ShowDiffOnly bool

	// ShowColor emphasizes changes with colors.
	ShowColor bool

	// ShowChangeTypePrefix adds change type prefix to each shown line.
	ShowChangeTypePrefix bool
}

// toYaml recursively formats a DiffNode into a string of differences in
// YAML-like format.
func (opts *FormatOptions) toYaml(n *DiffNode, depth int, tagList bool) string {
	// Stop descending if node is nil or if we are looking only for
	// differences and there is no change in this particular node.
	// If node represents a slice id, show it anyway (for readability).
	if n == nil || (opts.ShowDiffOnly && !n.HasChanged() && !n.isSliceId) {
		return ""
	}

	lines := []string{}
	line := opts.formatValue(n, depth, tagList)
	if line != "" {
		lines = append(lines, line)
	}

	// Increase depth only if node is not root to prevent unnecessary
	// indentation for base type comparisons.
	if !n.IsRoot() {
		depth++
	}

	// tagList indicates that list prefix has to be appended to the next
	// line.
	tagList = n.isSliceIndex()

	// Descend into node's children.
	for _, k := range n.ChildrenKeysSorted() {
		c := n.Child(k)

		line := opts.toYaml(c, depth, tagList)
		if line != "" {
			tagList = false
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}

func (opts *FormatOptions) formatValue(n *DiffNode, depth int, tagList bool) string {
	// Slice indexes are not shown. If root node is not a leaf, there is
	// also nothing to show.
	if !n.IsLeaf() && (n.isSliceIndex() || n.IsRoot()) {
		return ""
	}

	// Evaluate indentation. If tagList is true, replace last 2
	// spaces with a list prefix.
	indent := strings.Repeat(" ", depth*2)
	if tagList {
		indent = strings.TrimSuffix(indent, "  ") + "- "
	}

	key := n.pathKey
	value := ""
	prefix := ""

	bv := formatString(n.valueBefore)
	av := formatString(n.valueAfter)

	// Format value.
	switch n.changeType {
	case Create:
		value = av
	case Modify:
		value = fmt.Sprintf("%v -> %v", bv, av)
	default:
		value = bv
	}

	showPrefix := false

	// Apply change type prefix.
	if opts.ShowChangeTypePrefix {
		showPrefix = true
		switch n.changeType {
		case Create:
			prefix = "+"
		case Delete:
			prefix = "-"
		case Modify:
			prefix = "~"
		default:
			prefix = " "
		}
	}

	// Apply colors to key and value.
	if opts.ShowColor {
		switch n.changeType {
		case Create:
			key = green(key)
			value = green(value)
			prefix = green(prefix)
		case Delete:
			key = red(key)
			value = red(value)
			prefix = red(prefix)
		case Modify:
			value = yellow(value)
			prefix = yellow(prefix)
		}
	}

	if showPrefix {
		verticalLine := "\u2502"
		prefix = fmt.Sprintf("%s %s ", prefix, verticalLine)
	}

	// Edge case: if we are comparing raw lists, we have to show the value
	// of a slide index node.
	if n.IsLeaf() && n.isSliceIndex() {
		return fmt.Sprintf("%s%s- %s", prefix, indent, value)
	}

	// If root node is also a leaf, return only its value.
	if n.IsLeaf() && n.IsRoot() {
		return fmt.Sprintf("%s%s%s", prefix, indent, value)
	}

	// If node is neither a leaf nor root, return only its key.
	if !n.IsLeaf() && !n.IsRoot() {
		return fmt.Sprintf("%s%s%s:", prefix, indent, key)
	}

	return fmt.Sprintf("%s%s%s: %s", prefix, indent, key, value)
}

// formatString formats a value based on its type and returns it as a string.
func formatString(v any) string {
	switch v.(type) {
	case string:
		return fmt.Sprintf("\"%v\"", v)
	default:
		return fmt.Sprint(v)
	}
}
