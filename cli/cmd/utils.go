package cmd

import (
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
