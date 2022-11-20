package ui

import (
	"cli/env"
	"fmt"
	"os"
	"strings"
)

type Level uint8

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return ""
	}
}

type Ui struct {
}

// Ask asks user for confirmation. If user confirms with either "y" or "yes"
// nil is returned. Otherwise, if user enters "n" or "no" an error is returned.
func Ask(msg ...string) error {
	var question string
	var response string

	// Automatically approve if '--auto-approve' flag is used
	if env.AutoApprove {
		return nil
	}

	if len(msg) == 0 {
		question = "Would you like to continue?"
	} else {
		question = strings.Join(msg, " ")
	}

	Printf(DEBUG, "\n%s (yes/no)", question)

	if _, err := fmt.Scan(&response); err != nil {
		return fmt.Errorf("ask: %v", err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return nil
	case "n", "no":
		return fmt.Errorf("User aborted...")
	default:
		return Ask(msg...)
	}
}

func Print(level Level, msg ...any) {
	if level == DEBUG && !env.Debug {
		return
	}

	w := os.Stdout

	if level == ERROR || level == WARN {
		w = os.Stderr
	}

	fmt.Fprint(w, msg...)
}

func Println(level Level, msg ...any) {
	Print(level, msg...)
	Print(level, "\n")
}

func Printf(level Level, format string, args ...interface{}) {
	Print(level, fmt.Sprintf(format, args...))
}

func PrintBlock(err ...error) {
	var es []ErrorBlock

	for _, e := range err {
		if b, ok := e.(ErrorBlock); ok {
			es = append(es, b)
			continue
		}

		es = append(es, ErrorBlock{
			Level: ERROR,
			Content: []ErrorContent{
				NewErrorLine("Error:", fmt.Sprint(e)),
			},
		})
	}

	for _, e := range es {
		if env.NoColor {
			e.Level = INFO
		}

		fmt.Fprintln(os.Stderr, e)
	}
}
