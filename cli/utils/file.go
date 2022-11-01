package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"

	"gopkg.in/yaml.v3"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ReadFile reads a file from provided source path and returns it as a string.
func ReadFile(path string) (string, error) {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		return "", fmt.Errorf("Failed to read a file on path '%s': %v", path, err)
	}

	return string(file), nil
}

// CopyFile copies a file from source to destination path. If file already exists
// on the destination path, it will be overwritten.
func Copy(srcPath, dstPath string) error {
	opt := copy.Options{}

	opt.OnDirExists = func(src, dest string) copy.DirExistsAction {
		return copy.Untouchable
	}

	return copy.Copy(srcPath, dstPath, opt)
}

func ForceCopy(srcPath, dstPath string) error {
	opt := copy.Options{}

	opt.OnDirExists = func(src, dest string) copy.DirExistsAction {
		return copy.Replace
	}

	return copy.Copy(srcPath, dstPath, opt)
}

func CopyFile(srcPath, dstPath string) error {
	return copy.Copy(srcPath, dstPath)
}

// CopyFile copies a file from source to destination path. If file already exists
// on the destination path, it will be overwritten.
func CopyFile1(srcPath, dstPath string) error {
	sFile, err := os.Open(srcPath)

	if err != nil {
		return err
	}

	defer sFile.Close()

	dstDir := filepath.Dir(dstPath)

	err = os.MkdirAll(dstDir, os.ModePerm)

	if err != nil {
		return fmt.Errorf("Failed to create destination directory (%s) while copying a file: %v", dstDir, err)
	}

	dFile, err := os.Create(dstPath)

	if err != nil {
		return err
	}

	defer dFile.Close()

	_, err = io.Copy(dFile, sFile)
	return err
}

// ForceMove forcibly moves a file or directory to a specified location.
// First the destination file or directory is removed, and then the contents
// are moved there.
func ForceMove(srcPath string, dstPath string) error {
	err := os.RemoveAll(dstPath)

	if err != nil {
		return fmt.Errorf("Failed to force remove destination file: %w", err)
	}

	dstDir := filepath.Dir(dstPath)

	err = os.MkdirAll(dstDir, os.ModePerm)

	if err != nil {
		return fmt.Errorf("Failed to create destination directory (%s) while moving files: %v", dstDir, err)
	}

	err = os.Rename(srcPath, dstPath)

	if err != nil {
		return fmt.Errorf("Failed to move file from src (%s) to dst (%s) path: %w", srcPath, dstPath, err)
	}

	return nil
}

// ReadYaml reads yaml file on the given path and unmarshals it into the given type.
func ReadYaml[T interface{}](path string, typ T) (*T, error) {
	yml, err := ReadFile(path)

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

	if err != nil {
		return "", err
	}

	return string(arr), nil
}
