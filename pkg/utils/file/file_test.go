package file

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tmpFile(t *testing.T, name string, content ...string) string {
	t.Helper()

	if name == "" {
		name = "temp.file"
	}

	fPath := path.Join(t.TempDir(), name)
	fData := strings.Join(content, " ")

	err := os.WriteFile(fPath, []byte(fData), os.ModePerm)
	require.NoErrorf(t, err, "failed creating tmp file (%s): %v", name, err)

	return fPath
}

func tmpPath(t *testing.T) string {
	return path.Join(t.TempDir(), "tmp")
}

func TestExists(t *testing.T) {
	assert.False(t, Exists(tmpPath(t)))
	assert.True(t, Exists(tmpFile(t, "test")))
}

func TestRead(t *testing.T) {
	out, err := Read(tmpFile(t, "test.file", "test"))
	require.NoError(t, err)
	assert.Equal(t, "test", out)
}

func TestRead_Fail(t *testing.T) {
	_, err := Read(tmpPath(t))
	assert.ErrorContains(t, err, ": no such file or directory")
}

func TestCopy(t *testing.T) {
	src := tmpFile(t, "test.file", "test")
	dst := tmpPath(t)

	err := Copy(src, dst, os.ModePerm)
	require.NoError(t, err)
	assert.FileExists(t, src)
	assert.FileExists(t, dst)

	out, err := Read(dst)
	require.NoError(t, err)
	assert.Equal(t, "test", out)
}

func TestCopy_SourceNotFound(t *testing.T) {
	err := Copy(tmpPath(t), tmpPath(t), os.ModePerm)
	assert.ErrorContains(t, err, "no such file or directory")
}

func TestCopy_FileExists(t *testing.T) {
	src := tmpFile(t, "src.file")
	dst := tmpFile(t, "dst.file")
	err := Copy(src, dst, os.ModePerm)
	assert.ErrorContains(t, err, "destination file already exists")
}

func TestForceCopy(t *testing.T) {
	src := tmpFile(t, "src.file", "source")
	dst := tmpFile(t, "dst.file", "destination")

	err := ForceCopy(src, dst, os.ModePerm)
	require.NoError(t, err)

	out, err := Read(dst)
	require.NoError(t, err)
	assert.Equal(t, "source", out)
}

func TestAppend(t *testing.T) {
	src := tmpFile(t, "src.file", "source\n")

	require.NoError(t, Append(src, []byte("test")))
	require.NoError(t, Append(src, []byte("x=1.2.3")))
	require.NoError(t, Append(src, []byte("$ \n \\ ")))

	out, err := Read(src)
	require.NoError(t, err)
	assert.Equal(t, "source\ntest\nx=1.2.3\n$ \n \\ \n", out)
}

func TestWriteYaml(t *testing.T) {
	type T struct {
		Value int `yaml:"value"`
	}

	fPath := path.Join(t.TempDir(), "file.yaml")

	err := WriteYaml(T{7}, fPath, os.ModePerm)
	require.NoError(t, err)

	f, err := os.ReadFile(fPath)
	assert.Equal(t, "value: 7\n", string(f))
}

func TestReadYaml(t *testing.T) {
	type T struct {
		Value int `yaml:"value"`
	}

	type S struct {
		Test T `yaml:"test"`
	}

	fPath := path.Join(t.TempDir(), "file.yaml")

	err := WriteYaml(S{T{7}}, fPath, os.ModePerm)
	require.NoError(t, err)

	s, err := ReadYaml(fPath, S{})
	require.NoError(t, err, "failed reading YAML file")
	assert.Equal(t, 7, s.Test.Value)
}

func TestReadYaml_InvalidContent(t *testing.T) {
	type S struct {
		Value int
	}

	_, err := ReadYaml(tmpFile(t, "tmp.file", "\t"), S{})
	assert.ErrorContains(t, err, "yaml: found character that cannot start any token")
}

func TestReadYaml_FileNotExist(t *testing.T) {
	type S struct {
		Value int
	}

	_, err := ReadYaml(path.Join(t.TempDir(), "invalid"), S{})
	assert.ErrorContains(t, err, "no such file or directory")
}

func TestReadYamlStrict(t *testing.T) {
	type A struct {
		ValueA int `yaml:"a"`
	}

	type B struct {
		A      `yaml:",inline"`
		ValueB int `yaml:"b"`
	}

	tests := []struct {
		Content string
		Expect  B
		Error   string
	}{
		{
			Content: "a: 1",
			Expect:  B{A: A{ValueA: 1}},
		},
		{
			Content: "b: 1",
			Expect:  B{ValueB: 1},
		},
		{
			Content: "a: 1\nb: 1",
			Expect:  B{A: A{ValueA: 1}, ValueB: 1},
		},
		{
			Content: "b: 1\nc: 1",
			Error:   "field c not found in type",
		},
	}

	for _, test := range tests {
		path := filepath.Join(t.TempDir(), "tmp")
		_ = os.WriteFile(path, []byte(test.Content), 0777)

		out, err := ReadYamlStrict(path, B{})
		if test.Error == "" {
			require.NoError(t, err, "Expected no error!")
			assert.Equal(t, test.Expect, *out)
		} else {
			assert.ErrorContains(t, err, test.Error)
		}
	}
}
