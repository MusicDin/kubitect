package ui

import (
	"cli/ui/streams"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testFormat ignores pivot for cleaner tests
func testFormat(o streams.OutputStream, s string, i int) string {
	lines, _ := Format(o, s, i, 0)
	return strings.Join(lines, "\n")
}

// testFmtLine ignores pivot for cleaner tests
func testFmtLine(s string, w int) string {
	lines, _ := fmtLine(s, w, 0)
	return strings.Join(lines, "\n")
}

func TestFmtLine_Empty(t *testing.T) {
	assert.Equal(t, "", testFmtLine("", -1))
	assert.Equal(t, "", testFmtLine("", 0))
	assert.Equal(t, "", testFmtLine("", 1))
}

func TestFmtLine_Basic(t *testing.T) {
	assert.Equal(t, "123", testFmtLine("123", 3))
	assert.Equal(t, "1\n2\n3", testFmtLine("123", 1))
	assert.Equal(t, "\n1\n2\n3", testFmtLine("123", 0))
	assert.Equal(t, "abcd\n1234", testFmtLine("abcd1234", 4))
	assert.Equal(t, "thi\ns \nis \na t\nest", testFmtLine("this is a test", 3))
	assert.Equal(t, "this\nis a\ntest", testFmtLine("this is a test", 4))
	assert.Equal(t, "this \nis a \ntest", testFmtLine("this is a test", 5))
	assert.Equal(t, "this is a test", testFmtLine("this is a test", 100))
}

func TestFmtLine_pivot(t *testing.T) {
	lines, pivot := fmtLine("", 10, 0)
	assert.Equal(t, "", strings.Join(lines, "\n"))
	assert.Equal(t, 0, pivot)

	lines, pivot = fmtLine("test", 10, 10)
	assert.Equal(t, "test", strings.Join(lines, "\n"))
	assert.Equal(t, 4, pivot)

	lines, pivot = fmtLine("test", 10, 1)
	assert.Equal(t, " test", strings.Join(lines, "\n"))
	assert.Equal(t, 6, pivot)

	lines, pivot = fmtLine("test", 10, 9)
	assert.Equal(t, " \ntest", strings.Join(lines, "\n"))
	assert.Equal(t, 4, pivot)

	lines, pivot = fmtLine("\ntest test test", 10, 5)
	assert.Equal(t, " \n\ntest test\ntest", strings.Join(lines, "\n"))
	assert.Equal(t, 4, pivot)

	lines, pivot = fmtLine("\ntest test test", 10, 4)
	assert.Equal(t, " \ntest\ntest test", strings.Join(lines, "\n"))
	assert.Equal(t, 9, pivot)
}

func TestFmtLine_ContinuosFormat(t *testing.T) {
	lines, pivot := fmtLine("this is a test", 10, 0)
	assert.Equal(t, "this is a \ntest", strings.Join(lines, "\n"))
	assert.Equal(t, 4, pivot)

	lines, pivot = fmtLine("this is a test too", 10, pivot)
	assert.Equal(t, " this \nis a test \ntoo", strings.Join(lines, "\n"))
	assert.Equal(t, 3, pivot)
}

func TestFormat_Empty(t *testing.T) {
	assert.Equal(t, "", testFormat(nil, "", 0))
	assert.Equal(t, "", testFormat(nil, "", 1))
	assert.Equal(t, "", testFormat(nil, "", -1))
}

func TestFormat_TerminalStream(t *testing.T) {
	s := streams.MockTerminalStreams(t).Out()

	assert.Equal(t, "test", testFormat(s, "test", 0))
	assert.Equal(t, "te\nst", testFormat(s, "test", s.Columns()-2))
	assert.Equal(t, "t\ne\ns\nt", testFormat(s, "test", s.Columns()-1))
	assert.Equal(t, "\nt\ne\ns\nt", testFormat(s, "test", s.Columns())) // Maybe if (width <= indentation) => defaultCols
	assert.Equal(t, "te\nst", testFormat(s, "te\nst", 0))
	assert.Equal(t, "te\n\nst", testFormat(s, "te\n\nst", 0))
	assert.Equal(t, "te\nst\n", testFormat(s, "test\n", s.Columns()-2))
}

func TestFormat_NonTerminalStream(t *testing.T) {
	s := streams.MockStreams(t).Out()

	assert.Equal(t, "test", testFormat(s, "test", 0))
	assert.Equal(t, "test", testFormat(s, "test", s.Columns()-2))
	assert.Equal(t, "test", testFormat(s, "test", s.Columns()-1))
	assert.Equal(t, "test", testFormat(s, "test", s.Columns()))
	assert.Equal(t, "\ntest", testFormat(s, "\ntest", 0))
	assert.Equal(t, "te\nst", testFormat(s, "te\nst", 0))
	assert.Equal(t, "test\n", testFormat(s, "test\n", 0))
	assert.Equal(t, "te\nst", testFormat(s, "te\nst", s.Columns()))

	lines, pivot := Format(s, "test", 0, 42)
	assert.Equal(t, " test", strings.Join(lines, "\n"))
	assert.Equal(t, 1, pivot) // When stream is nil or non-terminal pivot is always 1
}
