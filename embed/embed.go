package embed

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
)

//go:embed templates presets ansible terraform
var efs embed.FS

type Resource struct {
	Name    string
	Path    string
	Content []byte
}

// GetResource returns an embedded resource on a given path.
// If the resource is not found, an error is thrown.
func GetResource(filePath string) (*Resource, error) {
	content, err := efs.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return &Resource{
		Name:    path.Base(filePath),
		Path:    filePath,
		Content: content,
	}, nil
}

// getAllResources recursively searches for all embedded files
// on the given path.
func getAllResources(rootDir string) ([]Resource, error) {
	var resources []Resource

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

		res := Resource{
			Name:    f.Name(),
			Path:    fPath,
			Content: content,
		}
		resources = append(resources, res)

		return nil
	}

	err := fs.WalkDir(efs, rootDir, filter)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

// MirrorResource copies an embedded resource to the given destination path,
// retaining (mirroring) its full path.
func MirrorResource(resPath, dstPath string) error {
	resPath = path.Clean(resPath)

	f, err := efs.Open(resPath)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		content, err := efs.ReadFile(resPath)
		if err != nil {
			return err
		}

		resDstPath := path.Join(dstPath, resPath)
		os.MkdirAll(path.Dir(resDstPath), os.ModePerm)
		return ioutil.WriteFile(resDstPath, content, os.ModePerm)
	}

	mirrorDir := func(fPath string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		resDstPath := path.Join(dstPath, fPath)

		if f.IsDir() {
			return os.MkdirAll(resDstPath, os.ModePerm)
		}

		content, err := efs.ReadFile(fPath)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(resDstPath, content, os.ModePerm)
	}

	return fs.WalkDir(efs, resPath, mirrorDir)
}

// GetResource returns a resource on a given path. If resource is not
// found, an error is thrown
func GetTemplate(tplPath string) (*Resource, error) {
	return GetResource(path.Join("templates", tplPath))
}

// Presets returns a list of all available presets.
func Presets() ([]Resource, error) {
	return getAllResources("presets")
}

// GetPreset returns a preset with a given name. If preset is not found,
// an error is returned.
func GetPreset(name string) (*Resource, error) {
	return GetResource(path.Join("presets", name))
}
