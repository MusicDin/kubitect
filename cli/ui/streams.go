package ui

import (
	"os"

	"golang.org/x/term"
)

const (
	defaultColumns    = 50
	defaultIsTerminal = false
)

type Streams struct {
	Out *OutputStream
	Err *OutputStream
	In  *InputStream
}

func StandardStreams() *Streams {
	return &Streams{
		Out: &OutputStream{
			File:       os.Stdout,
			isTerminal: isTerminal,
			columns:    columns,
		},
		Err: &OutputStream{
			File:       os.Stderr,
			isTerminal: isTerminal,
			columns:    columns,
		},
		In: &InputStream{
			File:       os.Stdin,
			isTerminal: isTerminal,
		},
	}
}

type OutputStream struct {
	File *os.File

	isTerminal func(*os.File) bool
	columns    func(*os.File) int
}

func (s *OutputStream) IsTerminal() bool {
	if s.isTerminal == nil {
		return defaultIsTerminal
	}

	return s.isTerminal(s.File)
}

func (s *OutputStream) Columns() int {
	if s.columns == nil {
		return defaultColumns
	}

	return s.columns(s.File)
}

type InputStream struct {
	File *os.File

	isTerminal func(*os.File) bool
}

func (s *InputStream) IsTerminal() bool {
	if s.isTerminal == nil {
		return defaultIsTerminal
	}

	return s.isTerminal(s.File)
}

func isTerminal(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

func columns(f *os.File) int {
	width, _, err := term.GetSize(int(f.Fd()))

	if err != nil {
		return defaultColumns
	}

	return width
}
