package file

import (
	"fmt"
	"io/fs"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

// Exists returns true if file metadata is obtained without errors.
// Otherwise, false is returned.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Read reads file on a given path and returns its content as a string.
func Read(path string) (string, error) {
	file, err := os.ReadFile(path)

	if err != nil {
		return "", fmt.Errorf("read file '%s': %v", path, err)
	}

	return string(file), nil
}

// Copy reads file from the source path and writes it to the destination path
// with the specified permissions. All required subdirectories are also
// created. An error is thrown if destination file exists.
func Copy(srcPath, dstPath string, mode fs.FileMode) error {
	if Exists(dstPath) {
		return fmt.Errorf("copy file: destination file already exists %s", dstPath)
	}

	return ForceCopy(srcPath, dstPath, mode)
}

// Append appends the provided byte content with an additional new line to the file.
func Append(path string, content []byte) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(append(content, '\n'))
	return err
}

// ForceCopy reads the file located at the source path and writes it to the
// destination path with the specified permissions. All required subdirectories
// are also created.
func ForceCopy(srcPath, dstPath string, mode fs.FileMode) error {
	file, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("copy file: %v", err)
	}

	if err := os.MkdirAll(path.Dir(dstPath), os.ModePerm); err != nil {
		return err
	}

	err2 := os.WriteFile(dstPath, file, mode)
	if err2 != nil {
		return fmt.Errorf("copy file: %v", err2)
	}

	return nil
}

// ReadYaml reads yaml file on the given path and unmarshals it into the given
// type.
func ReadYaml[T any](path string, typ T) (*T, error) {
	yml, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yml, &typ)
	return &typ, err
}

// ReadYaml reads yaml file on the given path and unmarshals it into the given
// type. File is not allowed to contain any unused/extra keys.
func ReadYamlStrict[T any](path string, typ T) (*T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	err = decoder.Decode(&typ)
	if err != nil {
		return nil, err
	}

	return &typ, nil
}

// WriteYaml writes a given object as a yaml file to the given path.
func WriteYaml(obj any, path string, perm fs.FileMode) error {
	yml, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	return os.WriteFile(path, yml, perm)
}
