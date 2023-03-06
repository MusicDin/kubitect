package ui

import (
	"fmt"
	"github.com/MusicDin/kubitect/pkg/ui/streams"
	"strings"
	"sync"
)

// Global Ui singleton
var (
	instance *ui
	once     sync.Once
)

func GlobalUi(opts ...UiOptions) Ui {
	if instance == nil {
		once.Do(func() {
			instance = &ui{
				streams: streams.StandardStreams(),
			}

			if len(opts) > 0 {
				o := opts[0]
				instance.autoApprove = o.AutoApprove
				instance.debug = o.Debug
				instance.noColor = o.NoColor
			}
		})
	}

	return instance
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

type (
	Ui interface {
		Ask(msg ...string) error
		Print(level Level, msg ...any)
		Printf(level Level, format string, args ...any)
		Println(level Level, msg ...any)
		PrintBlockE(err ...error)

		Streams() streams.Streams

		HasColor() bool
		Debug() bool
		AutoApprove() bool
	}

	ui struct {
		streams streams.Streams

		noColor     bool
		debug       bool
		autoApprove bool
	}
)

func HasColor() bool {
	return GlobalUi().HasColor()
}

func (u *ui) HasColor() bool {
	return !u.noColor
}

func Debug() bool {
	return GlobalUi().Debug()
}

func (u *ui) Debug() bool {
	return u.debug
}

func AutoApprove() bool {
	return GlobalUi().AutoApprove()
}

func (u *ui) AutoApprove() bool {
	return u.debug
}

func Streams() streams.Streams {
	return GlobalUi().Streams()
}

func (u *ui) Streams() streams.Streams {
	return u.streams
}

func (u *ui) outputStream(level Level) streams.OutputStream {
	switch level {
	case ERROR, WARN:
		return u.streams.Err()
	default:
		return u.streams.Out()
	}
}

func (u *ui) outputColor(l Level) Color {
	if u.noColor {
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

func Ask(msg ...string) error {
	return GlobalUi().Ask(msg...)
}

// Ask asks user for confirmation. If user confirms with either "y" or "yes"
// nil is returned. Otherwise, if user enters "n" or "no" an error is returned.
func (u *ui) Ask(msg ...string) error {
	var question string
	var response string

	// Automatically approve if '--auto-approve' flag is used
	if u.autoApprove {
		return nil
	}

	si := u.streams.In()

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

	if _, err := fmt.Fscan(si.File(), &response); err != nil {
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

func Print(level Level, msg ...any) {
	GlobalUi().Print(level, msg...)
}

func (u *ui) Print(level Level, msg ...any) {
	if level == DEBUG && !u.debug {
		return
	}

	w := u.outputStream(level).File()

	fmt.Fprint(w, msg...)
}

func Println(level Level, msg ...any) {
	GlobalUi().Println(level, msg...)
}

func (u *ui) Println(level Level, msg ...any) {
	u.Print(level, msg...)
	u.Print(level, "\n")
}

func Printf(level Level, format string, args ...any) {
	GlobalUi().Printf(level, format, args...)
}

func (u *ui) Printf(level Level, format string, args ...any) {
	u.Print(level, fmt.Sprintf(format, args...))
}

func PrintBlockE(errs ...error) {
	GlobalUi().PrintBlockE(errs...)
}

func (u *ui) PrintBlockE(errs ...error) {
	var eb ErrorBlock

	for _, e := range errs {
		switch e.(type) {
		case ErrorBlock:
			eb = e.(ErrorBlock)
		default:
			content := []Content{
				NewErrorLine("Error:", fmt.Sprint(e)),
			}

			eb = NewErrorBlock(ERROR, content)
		}

		s := u.outputStream(eb.Severity())
		c := u.outputColor(eb.Severity())

		fmt.Fprintln(s.File(), eb.Format(s, c))
	}
}
