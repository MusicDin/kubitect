package embed

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"path"
)

//go:embed presets
var all embed.FS

type Preset struct {
	Name    string
	Path    string
	Content []byte
}

func Presets() ([]Preset, error) {
	var presets []Preset

	filter := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			content, err := all.ReadFile(path)
			if err != nil {
				return err
			}

			f := Preset{
				Name:    presetName(path),
				Path:    path,
				Content: content,
			}
			presets = append(presets, f)
		}

		return nil
	}

	err := fs.WalkDir(all, "presets", filter)
	if err != nil {
		return nil, err
	}

	return presets, nil
}

func GetPreset(name string) (*Preset, error) {
	presets, err := Presets()
	if err != nil {
		return nil, err
	}

	for _, p := range presets {
		if p.Name == name {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("get preset: preset %s not found", name)
}

// presetName converts file name without extension form the given path.
func presetName(fPath string) string {
	return path.Base(fPath[:len(fPath)-len(path.Ext(fPath))])
}
