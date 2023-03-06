package template

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockTemplateFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "test.tpl")
	err := os.WriteFile(path, []byte(content), os.ModePerm)
	assert.NoError(t, err)

	return path
}

type TemplateMock struct{ Value string }
type InvalidTemplateMock struct{ TemplateMock }
type InvalidFieldTemplateMock struct{ TemplateMock }
type CustomDelimsTemplateMock struct{ TemplateMock }
type CustomFuncsTemplateMock struct{ TemplateMock }

func (t TemplateMock) Name() string                             { return "test.tpl" }
func (t TemplateMock) Template() string                         { return "Test {{ .Value }}" }
func (t InvalidTemplateMock) Template() string                  { return "{{ \\ }}" }
func (t InvalidFieldTemplateMock) Template() string             { return "Test {{ .Invalid }}" }
func (t CustomDelimsTemplateMock) Template() string             { return "<< if true >>success<< end >>" }
func (t CustomDelimsTemplateMock) Delimiters() (string, string) { return "<<", ">>" }
func (t CustomFuncsTemplateMock) Template() string              { return "{{ alwaysTrue }}" }
func (t CustomFuncsTemplateMock) Functions() map[string]any {
	return map[string]any{
		"alwaysTrue": func() bool { return true },
	}
}

func TestPopulate(t *testing.T) {
	tpl := TemplateMock{Value: "test"}
	str, err := Populate(tpl)
	assert.NoError(t, err)
	assert.Equal(t, "Test test", str)
}

func TestPopulate_Invalid(t *testing.T) {
	_, err := Populate(InvalidTemplateMock{})
	assert.ErrorContains(t, err, "unexpected")
}

func TestPopulate_InvalidField(t *testing.T) {
	_, err := Populate(InvalidFieldTemplateMock{})
	assert.ErrorContains(t, err, "can't evaluate field Invalid")
}

func TestPopulate_CustomDelims(t *testing.T) {
	tpl := CustomDelimsTemplateMock{}
	str, err := Populate(tpl)
	assert.NoError(t, err)
	assert.Equal(t, "success", str)
}

func TestPopulate_CustomFunctions(t *testing.T) {
	tpl := CustomFuncsTemplateMock{}
	str, err := Populate(tpl)
	assert.NoError(t, err)
	assert.Equal(t, "true", str)
}

func TestPopulateFrom(t *testing.T) {
	path := MockTemplateFile(t, "{{ .Value }}")

	tpl := TemplateMock{Value: "test"}
	str, err := PopulateFrom(tpl, path)
	assert.NoError(t, err)
	assert.Equal(t, "test", str)
}

func TestPopulateFrom_Invalid(t *testing.T) {
	path := MockTemplateFile(t, "{{ \\ }}")

	_, err := PopulateFrom(TemplateMock{}, path)
	assert.ErrorContains(t, err, "unexpected")
}

func TestPopulateFrom_InvalidPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "invalid.tpl")

	_, err := PopulateFrom(TemplateMock{}, path)
	assert.ErrorContains(t, err, "no such file or directory")
}

func TestWrite_InvalidDir(t *testing.T) {
	dir := path.Join(t.TempDir(), "dir")

	// Create file on dir location
	_, err := os.Create(dir)
	assert.NoError(t, err)

	assert.ErrorContains(t, write(path.Join(dir, "tpl"), nil), "not a directory")
}

func TestWrite(t *testing.T) {
	tpl := TemplateMock{Value: "test"}
	tplPath := path.Join(t.TempDir(), "test")

	err := Write(tpl, tplPath)
	assert.NoError(t, err)

	str, err := os.ReadFile(tplPath)
	assert.NoError(t, err)
	assert.Equal(t, "Test test", string(str))
}

func TestWrite_Invalid(t *testing.T) {
	tpl := InvalidTemplateMock{}
	tplPath := path.Join(t.TempDir(), "test")

	err := Write(tpl, tplPath)
	assert.ErrorContains(t, err, "unexpected")
}

func TestWriteFrom(t *testing.T) {
	tpl := TemplateMock{Value: "test"}
	srcPath := MockTemplateFile(t, "{{ .Value }}")
	dstPath := path.Join(t.TempDir(), "test")

	err := WriteFrom(tpl, srcPath, dstPath)
	assert.NoError(t, err)

	str, err := os.ReadFile(dstPath)
	assert.NoError(t, err)
	assert.Equal(t, "test", string(str))
}

func TestWriteFrom_Invalid(t *testing.T) {
	srcPath := MockTemplateFile(t, "{{ \\ }}")
	dstPath := path.Join(t.TempDir(), "test")

	err := WriteFrom(TemplateMock{}, srcPath, dstPath)
	assert.ErrorContains(t, err, "unexpected")
}

func TestEmpty(t *testing.T) {
	assert.True(t, empty(""))
	assert.True(t, empty("  "))
	assert.True(t, empty(" \n "))
	assert.True(t, empty(" \n "))
	assert.True(t, empty(" \t "))
	assert.False(t, empty("a"))
}

func TestCountLeadingSpaces(t *testing.T) {
	assert.Equal(t, leadingSpaces(""), 0)
	assert.Equal(t, leadingSpaces(" "), 1)
	assert.Equal(t, leadingSpaces("aa "), 0)
	assert.Equal(t, leadingSpaces(" aa "), 1)
}

func TestTrimLeadingEmptyLines(t *testing.T) {
	assert.Equal(t, []string{}, trimLeadingEmptyLines([]string{}))
	assert.Equal(t, []string{}, trimLeadingEmptyLines([]string{"", "", ""}))
	assert.Equal(t, []string{}, trimLeadingEmptyLines([]string{"", " ", "\n", "\t"}))
	assert.Equal(t, []string{"a"}, trimLeadingEmptyLines([]string{"", "a"}))
	assert.Equal(t, []string{"a", ""}, trimLeadingEmptyLines([]string{"a", ""}))
	assert.Equal(t, []string{"a", ""}, trimLeadingEmptyLines([]string{"", "a", ""}))
	assert.Equal(t, []string{"a", ""}, trimLeadingEmptyLines([]string{" ", "a", ""}))
}

func TestTrimTrailingEmptyLines(t *testing.T) {
	assert.Equal(t, []string{}, trimTrailingEmptyLines([]string{}))
	assert.Equal(t, []string{}, trimTrailingEmptyLines([]string{"", " ", "\n", "\t"}))
	assert.Equal(t, []string{"a"}, trimTrailingEmptyLines([]string{"a", ""}))
	assert.Equal(t, []string{"", "a"}, trimTrailingEmptyLines([]string{"", "a"}))
	assert.Equal(t, []string{"", "a"}, trimTrailingEmptyLines([]string{"", "a", ""}))
	assert.Equal(t, []string{"", "a"}, trimTrailingEmptyLines([]string{"", "a", " "}))
}

func TestTrimEmptyLines(t *testing.T) {
	assert.Equal(t, []string{}, trimEmptyLines([]string{}))
	assert.Equal(t, []string{}, trimEmptyLines([]string{"", " ", "\n", "\t"}))
	assert.Equal(t, []string{"a"}, trimEmptyLines([]string{"a", ""}))
	assert.Equal(t, []string{"a"}, trimEmptyLines([]string{"", "a"}))
	assert.Equal(t, []string{"a"}, trimEmptyLines([]string{"", "a", ""}))
	assert.Equal(t, []string{"a"}, trimEmptyLines([]string{" ", "a", " "}))
}

func TestTrimLines(t *testing.T) {
	assert.Equal(t, []string{}, trimLines([]string{}))
	assert.Equal(t, []string{}, trimLines([]string{"", " ", "\n", "\t"}))
	assert.Equal(t, []string{"a"}, trimLines([]string{" a ", ""}))
	assert.Equal(t, []string{"a"}, trimLines([]string{" ", " a ", ""}))
	assert.Equal(t, []string{"\na"}, trimLines([]string{"", " \na "}))
	assert.Equal(t, []string{"\ta"}, trimLines([]string{"", " \ta ", ""}))
	assert.Equal(t, []string{"a", "", "a"}, trimLines([]string{"", "a", "", "a", ""}))
}

func TestTrimTemplate(t *testing.T) {
	in := `
		test:
			value: 7
	`

	assert.Equal(t, "test:\n  value: 7", TrimTemplate(in))
}
