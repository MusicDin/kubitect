package embed

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

//go:embed resources templates presets
var efs embed.FS

type EmbeddedFile struct {
	Name    string
	Path    string
	Content []byte
}

// get returns an embedded file on the given path.
func get(rootDir, filePath string) (*EmbeddedFile, error) {
	content, err := efs.ReadFile(path.Join(rootDir, filePath))
	if err != nil {
		return nil, err
	}

	return &EmbeddedFile{
		Name:    path.Base(filePath),
		Path:    filePath,
		Content: content,
	}, nil
}

// getAll recursively searches for all embedded files on the given path.
func getAll(rootDir string) ([]EmbeddedFile, error) {
	var files []EmbeddedFile

	filter := func(fPath string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		content, err := efs.ReadFile(fPath)
		if err != nil {
			return err
		}

		ef := EmbeddedFile{
			Name:    f.Name(),
			Path:    fPath,
			Content: content,
		}
		files = append(files, ef)

		return nil
	}

	err := fs.WalkDir(efs, rootDir, filter)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// mirror copies the embedded file to the given destination path,
// retaining (mirroring) its full path.
func mirror(rootDir, resPath, dstPath string) error {
	resAbsPath := path.Join(rootDir, resPath)

	f, err := efs.Open(resAbsPath)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		content, err := efs.ReadFile(resAbsPath)
		if err != nil {
			return err
		}

		resDstPath := path.Join(dstPath, resPath)
		os.MkdirAll(path.Dir(resDstPath), os.ModePerm)
		return ioutil.WriteFile(resDstPath, content, os.ModePerm)
	}

	copyDir := func(path string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		resRelPath, _ := filepath.Rel(rootDir, path)
		resDstPath := filepath.Join(dstPath, resRelPath)

		if f.IsDir() {
			return os.MkdirAll(resDstPath, os.ModePerm)
		}

		content, err := efs.ReadFile(path)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(resDstPath, content, os.ModePerm)
	}

	return fs.WalkDir(efs, resAbsPath, copyDir)
}

// GetResource returns a resource on a given path. If resource is not
// found, an error is thrown
func GetTemplate(tplPath string) (*EmbeddedFile, error) {
	return get("templates", tplPath)
}

// GetResource returns a resource on a given path. If resource is not
// found, an error is thrown
func GetResource(resPath string) (*EmbeddedFile, error) {
	return get("resources", resPath)
}

// MirrorResource copies the resource to the given destination path,
// retaining (mirroring) its full path.
func MirrorResource(resPath string, dstPath string) error {
	return mirror("resources", resPath, dstPath)
}

// Presets returns a list of all available presets.
func Presets() ([]EmbeddedFile, error) {
	return getAll("presets")
}

// GetPreset returns a preset with a given name. If preset is not found,
// an error is returned.
func GetPreset(name string) (*EmbeddedFile, error) {
	return get("presets", name)
}
