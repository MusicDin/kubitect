package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"

	"gopkg.in/yaml.v3"
)

// Exists returns true if file metadata is obtained without errors.
// Otherwise, false is returned.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Read reads a file on the given path and returns it as a string.
func Read(path string) (string, error) {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		return "", fmt.Errorf("read file '%s': %v", path, err)
	}

	return string(file), nil
}

// MakeDir creates a directory named path, along with any necessary parents.
func MakeDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		return fmt.Errorf("make directory '%s': %v", path, err)
	}

	return nil
}

// Copy gracefully copies a file from source to destination path. Files
// that already exist on the destination path will be left untouched.
func Copy(srcPath, dstPath string) error {
	opt := copy.Options{}

	opt.OnDirExists = func(src, dest string) copy.DirExistsAction {
		return copy.Untouchable
	}

	err := copy.Copy(srcPath, dstPath, opt)

	if err != nil {
		return fmt.Errorf("copy from '%s' to '%s': %v", srcPath, dstPath, err)
	}

	return nil
}

// ForceCopy copies a file from source to destination path. Files that
// already exist on the destination path will be replaced (overwritten).
func ForceCopy(srcPath, dstPath string) error {
	opt := copy.Options{}

	opt.OnDirExists = func(src, dest string) copy.DirExistsAction {
		return copy.Replace
	}

	err := copy.Copy(srcPath, dstPath, opt)

	if err != nil {
		return fmt.Errorf("force copy from '%s' to '%s': %v", srcPath, dstPath, err)
	}

	return nil
}

// Move moves a file or directory to a specified location. First the
// destination file or directory is removed, and then the source
// content is moved.
func Move(srcPath string, dstPath string) error {
	if err := Remove(dstPath); err != nil {
		return fmt.Errorf("force move: %v", err)
	}

	if err := MakeDir(filepath.Dir(dstPath)); err != nil {
		return fmt.Errorf("force move: %v", err)
	}

	if err := os.Rename(srcPath, dstPath); err != nil {
		return fmt.Errorf("force move from '%s' to '%s': %v", srcPath, dstPath, err)
	}

	return nil
}

// Remove removes directory and any children it contains.
func Remove(path string) error {
	err := os.RemoveAll(path)

	if err != nil {
		return fmt.Errorf("remove '%s': %v", path, err)
	}

	return nil
}

// ReadYaml reads yaml file on the given path and unmarshals it into the given type.
func ReadYaml[T interface{}](path string, typ T) (*T, error) {
	yml, err := Read(path)

	if err != nil {
		return nil, err
	}

	return UnmarshalYaml(string(yml), typ)
}

// UnmarshalYaml unmarshals yaml string to a given type.
func UnmarshalYaml[T interface{}](yml string, typ T) (*T, error) {
	err := yaml.Unmarshal([]byte(yml), &typ)

	if err != nil {
		return nil, err
	}

	return &typ, nil
}

// MarshalYaml marshals given object into a string.
func MarshalYaml(value interface{}) (string, error) {
	arr, err := yaml.Marshal(value)
	return string(arr), err
}
