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

// Global Ui singleton
var globalUi *Ui

type Ui struct {
	Streams *Streams
}

func (u *Ui) Stream(level Level) *OutputStream {
	switch level {
	case ERROR, WARN:
		return u.Streams.Err
	default:
		return u.Streams.Out
	}
}

func GlobalUi() *Ui {
	if globalUi == nil {
		globalUi = &Ui{
			Streams: StandardStreams(),
		}
	}

	return globalUi
}

// Ask asks user for confirmation. If user confirms with either "y" or "yes"
// nil is returned. Otherwise, if user enters "n" or "no" an error is returned.
func (u *Ui) Ask(msg ...string) error {
	var question string
	var response string

	// Automatically approve if '--auto-approve' flag is used
	if env.AutoApprove {
		return nil
	}

	si := u.Streams.In

	// Auto approve if stdin is not a terminal
	if si == nil || !si.IsTerminal() {
		return nil
	}

	if len(msg) == 0 {
		question = "Would you like to continue?"
	} else {
		question = strings.Join(msg, " ")
	}

	u.Printf(INFO, "\n%s (yes/no) ", question)

	if _, err := fmt.Fscan(si.File, &response); err != nil {
		return fmt.Errorf("ask: %v", err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return nil
	case "n", "no":
		return fmt.Errorf("User aborted...")
	default:
		return u.Ask(msg...)
	}
}

func (u *Ui) Print(level Level, msg ...any) {
	if level == DEBUG && !env.Debug {
		return
	}

	w := u.Stream(level).File

	fmt.Fprint(w, msg...)
}

func (u *Ui) Println(level Level, msg ...any) {
	u.Print(level, msg...)
	u.Print(level, "\n")
}

func (u *Ui) Printf(level Level, format string, args ...interface{}) {
	u.Print(level, fmt.Sprintf(format, args...))
}

func (u *Ui) PrintBlockE(err ...error) {
	var eb ErrorBlock

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

		s := u.Stream(eb.Severity)

		fmt.Fprintln(s.File, eb.Format(s, eb.Severity))
	}
}
