package cmd

import (
	"fmt"
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

// Example trims alls leading and trailing spaces from each line
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

	return strings.Join(trimmed, "\n")
}
