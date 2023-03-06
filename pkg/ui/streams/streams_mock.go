package streams

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockStreams(t *testing.T) *streams {
	t.Helper()

	dir := t.TempDir()

	fIn := tmpFile(t, dir, "in")
	fOut := tmpFile(t, dir, "out")
	fErr := tmpFile(t, dir, "err")

	return &streams{
		in: &inputStream{
			file: fIn,
		},
		out: &outputStream{
			file: fOut,
		},
		err: &outputStream{
			file: fErr,
		},
	}
}

func MockTerminalStreams(t *testing.T) *streams {
	t.Helper()

	isTerminal := func(f *os.File) bool { return true }
	columns := func(f *os.File) int { return 42 }

	streams := MockStreams(t)
	streams.in.isTerminal = isTerminal
	streams.out.isTerminal = isTerminal
	streams.out.columns = columns
	streams.err.isTerminal = isTerminal
	streams.err.columns = columns

	return streams
}

func MockEmptyStreams(t *testing.T) *streams {
	t.Helper()

	return &streams{
		in: &inputStream{
			file: nil,
		},
		out: &outputStream{
			file: nil,
		},
		err: &outputStream{
			file: nil,
		},
	}
}

func tmpFile(t *testing.T, dir, name string) *os.File {
	t.Helper()

	f, err := os.CreateTemp(dir, name)
	assert.NoErrorf(t, err, "failed creating tmp file (%s)", name)
	return f
}
