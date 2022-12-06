package ui

import (
	"fmt"
	"strings"
)

// Global Ui singleton
var globalUi *Ui

func GlobalUi() *Ui {
	if globalUi == nil {
		globalUi = &Ui{
			Streams: StandardStreams(),
		}
	}

	return globalUi
}

type Level uint8

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

type UiOptions struct {
	NoColor     bool
	Debug       bool
	AutoApprove bool
}

type Ui struct {
	Streams *Streams

	NoColor     bool
	Debug       bool
	autoApprove bool
}

func NewUi(o UiOptions) *Ui {
	return &Ui{
		Streams:     StandardStreams(),
		NoColor:     o.NoColor,
		Debug:       o.Debug,
		autoApprove: o.AutoApprove,
	}
}

func (u *Ui) Stream(level Level) *OutputStream {
	switch level {
	case ERROR, WARN:
		return u.Streams.Err
	default:
		return u.Streams.Out
	}
}

func (ui *Ui) Color(l Level) Color {
	if ui.NoColor {
		return Colors.NONE
	}

	switch l {
	case WARN:
		return Colors.YELLOW
	case ERROR:
		return Colors.RED
	default:
		return Colors.NONE
	}
}

// Ask asks user for confirmation. If user confirms with either "y" or "yes"
// nil is returned. Otherwise, if user enters "n" or "no" an error is returned.
func (u *Ui) Ask(msg ...string) error {
	var question string
	var response string

	// Automatically approve if '--auto-approve' flag is used
	if u.autoApprove {
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
	if level == DEBUG && !u.Debug {
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

		fmt.Fprintln(s.File, eb.Format(s, u.Color(eb.Severity)))
	}
}
