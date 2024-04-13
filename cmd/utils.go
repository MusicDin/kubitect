package main

import (
	"fmt"
	"path"
	"strings"
)

// LongDesc trims alls leading and trailing spaces from each line.
func LongDesc(s string) string {
	var trimmed []string

	for _, s := range strings.Split(s, "\n") {
		trimmed = append(trimmed, strings.TrimSpace(s))
	}

	return strings.TrimSpace(strings.Join(trimmed, "\n"))
}

// Example trims all leading and trailing spaces from each line
// and indents it with some spaces.
func Example(s string) string {
	var trimmed []string

	indent := 2

	for i, s := range strings.Split(s, "\n") {
		if i == 0 && s == "" {
			continue
		}

		t := fmt.Sprintf("%*s%s", indent, "", strings.TrimSpace(s))
		trimmed = append(trimmed, t)
	}

	out := strings.Join(trimmed, "\n")
	return strings.TrimRight(out, " ")
}

// presetName extracts file name without extension from the given path.
func presetName(fPath string) string {
	base := path.Base(fPath)
	ext := path.Ext(base)

	if base == ext {
		return base
	}

	return base[:len(base)-len(ext)]
}
