package file

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tmpFile(t *testing.T, content string) *os.File {
	t.Helper()

	fName := "file"

	f, err := os.CreateTemp(t.TempDir(), fName)
	if err != nil {
		t.Errorf("failed creating tmp file (%s): %v", fName, err)
	}

	f.Write([]byte(content))

	return f
}

func tmpPath(t *testing.T) string {
	return path.Join(t.TempDir(), "tmp")
}

func TestExists(t *testing.T) {
	assert.False(t, Exists(tmpPath(t)))
	assert.True(t, Exists(tmpFile(t, "test").Name()))
}

func TestMakeDir(t *testing.T) {
	dir := t.TempDir() + "a/b/c/d/e/f"
	err := MakeDir(dir)

	assert.NoError(t, err)
	assert.True(t, Exists(dir))
}

func TestMakeDir_Fail(t *testing.T) {
	err := MakeDir("")
	assert.ErrorContains(t, err, ": no such file or directory")
}

func TestRead(t *testing.T) {
	f := tmpFile(t, "test")

	out, err := Read(f.Name())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "test", out)
}

func TestRead_Fail(t *testing.T) {
	tmp := tmpPath(t)

	_, err := Read(tmp)
	assert.ErrorContains(t, err, ": no such file or directory")
}

func TestCopy(t *testing.T) {
	tmp := tmpPath(t)
	f := tmpFile(t, "test")

	if err := Copy(f.Name(), tmp); err != nil {
		t.Error(err)
	}

	out, err := Read(tmp)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "test", out)
	assert.True(t, Exists(f.Name()))
	assert.True(t, Exists(tmp))
}

func TestCopy_FileExists(t *testing.T) {
	tmp := tmpFile(t, "wrong")
	f := tmpFile(t, "test")

	assert.ErrorContains(t, Copy(f.Name(), tmp.Name()), ": destination already exists")
}

func TestCopy_Fail(t *testing.T) {
	assert.ErrorContains(t, Copy("", ""), ": no such file or directory")
}

func TestForceCopy(t *testing.T) {
	tmp := tmpFile(t, "wrong")
	f := tmpFile(t, "test")

	if err := ForceCopy(f.Name(), tmp.Name()); err != nil {
		t.Error(err)
	}

	out, err := Read(tmp.Name())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "test", out)
	assert.True(t, Exists(f.Name()))
	assert.True(t, Exists(tmp.Name()))
}

func TestForceCopy_Fail(t *testing.T) {
	assert.ErrorContains(t, ForceCopy("", ""), ": no such file or directory")
}

func TestMove(t *testing.T) {
	tmp := tmpFile(t, "wrong")
	f := tmpFile(t, "test")

	if err := Move(f.Name(), tmp.Name()); err != nil {
		t.Error(err)
	}

	out, err := Read(tmp.Name())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "test", out)
	assert.False(t, Exists(f.Name()))
	assert.True(t, Exists(tmp.Name()))
}

func TestMove_Fail(t *testing.T) {
	assert.ErrorContains(t, Move("", ""), ": no such file or directory")

	tmp := t.TempDir() + string(os.PathSeparator) + "."
	assert.ErrorContains(t, Move("", tmp), ": invalid argument")
}

func TestRemove(t *testing.T) {
	f := tmpFile(t, "test")

	assert.True(t, Exists(f.Name()))

	if err := Remove(f.Name()); err != nil {
		t.Error(err)
	}

	assert.False(t, Exists(f.Name()))
}

func TestRemove_Fail(t *testing.T) {
	tmp := t.TempDir() + string(os.PathSeparator) + "."
	assert.ErrorContains(t, Remove(tmp), ": invalid argument")
}

func TestWriteYaml(t *testing.T) {
	type T struct {
		Value int `yaml:"value"`
	}

	fPath := path.Join(t.TempDir(), "file.yaml")

	err := WriteYaml(T{7}, fPath, os.ModePerm)
	assert.NoError(t, err)

	f, err := Read(fPath)
	assert.Equal(t, "value: 7\n", f)
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
	assert.NoError(t, err)

	s, err := ReadYaml(fPath, S{})
	assert.NoError(t, err, "failed reading YAML file")
	assert.Equal(t, 7, s.Test.Value)
}

func TestReadYaml_InvalidContent(t *testing.T) {
	type S struct {
		Value int
	}

	_, err := ReadYaml(tmpFile(t, "\t").Name(), S{})
	assert.ErrorContains(t, err, "yaml: found character that cannot start any token")
}

func TestReadYaml_FileNotExist(t *testing.T) {
	type S struct {
		Value int
	}

	_, err := ReadYaml(path.Join(t.TempDir(), "invalid"), S{})
	assert.ErrorContains(t, err, "no such file or directory")
}
