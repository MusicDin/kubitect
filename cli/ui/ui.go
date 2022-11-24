package ui

import (
	"cli/env"
	"fmt"
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

func (t Level) Color() Color {
	if env.NoColor {
		return Colors.NONE
	}

	switch t {
	case WARN:
		return Colors.YELLOW
	case ERROR:
		return Colors.RED
	default:
		return Colors.NONE
	}
}

type GlobalUi struct {
	streams *Streams
}

func (u *GlobalUi) OutStream(level Level) *OutputStream {
	switch level {
	case ERROR, WARN:
		return u.streams.Err
	default:
		return u.streams.Out
	}
}

// Ui singleton
var globalUi *GlobalUi

func Ui() *GlobalUi {
	if globalUi == nil {
		globalUi = &GlobalUi{
			streams: StandardStreams(),
		}
	}

	return globalUi
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

	Printf(INFO, "\n%s (yes/no) ", question)

	if _, err := fmt.Fscan(Ui().streams.In.File, &response); err != nil {
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

	w := Ui().OutStream(level).File

	fmt.Fprint(w, msg...)
}

func Println(level Level, msg ...any) {
	Print(level, msg...)
	Print(level, "\n")
}

func Printf(level Level, format string, args ...interface{}) {
	Print(level, fmt.Sprintf(format, args...))
}

func PrintBlockE(err ...error) {
	var eb ErrorBlock

	s := Ui().OutStream(eb.Severity)

	for _, e := range err {
		switch e.(type) {
		case ErrorBlock:
			eb = e.(ErrorBlock)
		default:
			eb = NewErrorBlock(ERROR,
				[]Content{
					NewErrorLine("Error:", fmt.Sprint(e)),
				},
			)
		}

		fmt.Fprintln(s.File, eb.Format(s, eb.Severity))
	}
}
