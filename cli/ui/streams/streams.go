package streams

import (
	"os"

	"golang.org/x/term"
)

const (
	defaultColumns    = 50
	defaultIsTerminal = false
)

type (
	Streams interface {
		Out() OutputStream
		Err() OutputStream
		In() InputStream
	}

	streams struct {
		out *outputStream
		err *outputStream
		in  *inputStream
	}
)

func (s *streams) Out() OutputStream {
	return s.out
}

func (s *streams) Err() OutputStream {
	return s.err
}

func (s *streams) In() InputStream {
	return s.in
}

type (
	OutputStream interface {
		File() *os.File
		IsTerminal() bool
		Columns() int
	}

	outputStream struct {
		file       *os.File
		isTerminal func(*os.File) bool
		columns    func(*os.File) int
	}
)

func (s *outputStream) File() *os.File {
	return s.file
}

func (s *outputStream) IsTerminal() bool {
	if s.isTerminal == nil {
		return defaultIsTerminal
	}

	return s.isTerminal(s.file)
}

func (s *outputStream) Columns() int {
	if s.columns == nil {
		return defaultColumns
	}

	return s.columns(s.file)
}

type (
	InputStream interface {
		File() *os.File
		IsTerminal() bool
	}

	inputStream struct {
		file       *os.File
		isTerminal func(*os.File) bool
	}
)

func (s *inputStream) File() *os.File {
	return s.file
}

func (s *inputStream) IsTerminal() bool {
	if s.isTerminal == nil {
		return defaultIsTerminal
	}

	return s.isTerminal(s.file)
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

func StandardStreams() Streams {
	return &streams{
		out: &outputStream{
			file:       os.Stdout,
			isTerminal: isTerminal,
			columns:    columns,
		},
		err: &outputStream{
			file:       os.Stderr,
			isTerminal: isTerminal,
			columns:    columns,
		},
		in: &inputStream{
			file:       os.Stdin,
			isTerminal: isTerminal,
		},
	}
}
