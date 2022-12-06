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

func TestFile_Exists(t *testing.T) {
	assert.False(t, Exists(tmpPath(t)))
	assert.True(t, Exists(tmpFile(t, "test").Name()))
}

func TestFile_MakeDir(t *testing.T) {
	dir := t.TempDir() + "a/b/c/d/e/f"
	err := MakeDir(dir)

	assert.NoError(t, err)
	assert.True(t, Exists(dir))
}

func TestFile_MakeDirFail(t *testing.T) {
	err := MakeDir("")
	assert.ErrorContains(t, err, ": no such file or directory")
}

func TestFile_Read(t *testing.T) {
	f := tmpFile(t, "test")

	out, err := Read(f.Name())

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "test", out)
}

func TestFile_ReadFail(t *testing.T) {
	tmp := tmpPath(t)

	_, err := Read(tmp)
	assert.ErrorContains(t, err, ": no such file or directory")
}

func TestFile_Copy(t *testing.T) {
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

func TestFile_CopyFileExists(t *testing.T) {
	tmp := tmpFile(t, "wrong")
	f := tmpFile(t, "test")

	assert.ErrorContains(t, Copy(f.Name(), tmp.Name()), ": destination already exists")
}

func TestFile_CopyFail(t *testing.T) {
	assert.ErrorContains(t, Copy("", ""), ": no such file or directory")
}

func TestFile_ForceCopy(t *testing.T) {
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

func TestFile_ForceCopyFail(t *testing.T) {
	assert.ErrorContains(t, ForceCopy("", ""), ": no such file or directory")
}

func TestFile_Move(t *testing.T) {
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

func TestFile_MoveFail(t *testing.T) {
	assert.ErrorContains(t, Move("", ""), ": no such file or directory")

	tmp := t.TempDir() + string(os.PathSeparator) + "."
	assert.ErrorContains(t, Move("", tmp), ": invalid argument")
}

func TestFile_Remove(t *testing.T) {
	f := tmpFile(t, "test")

	assert.True(t, Exists(f.Name()))

	if err := Remove(f.Name()); err != nil {
		t.Error(err)
	}

	assert.False(t, Exists(f.Name()))
}

func TestFile_RemoveFail(t *testing.T) {
	tmp := t.TempDir() + string(os.PathSeparator) + "."
	assert.ErrorContains(t, Remove(tmp), ": invalid argument")
}

func TestFile_Yaml(t *testing.T) {
	type T struct {
		Value int `yaml:"value"`
	}

	type S struct {
		Test T `yaml:"test"`
	}

	yaml, err := MarshalYaml(S{T{7}})

	assert.NoError(t, err)
	assert.Equal(t, "test:\n    value: 7\n", yaml)

	f := tmpFile(t, yaml)

	s, err := ReadYaml(f.Name(), S{})

	if err != nil {
		t.Errorf("failed reading YAML file: %v", err)
	}

	assert.Equal(t, 7, s.Test.Value)
}

func TestFile_YamlFail(t *testing.T) {
	type S struct {
		Value int
	}

	_, err := ReadYaml(tmpFile(t, "\t").Name(), S{})

	assert.ErrorContains(t, err, "yaml: found character that cannot start any token")
}
