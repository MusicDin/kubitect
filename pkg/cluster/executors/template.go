package executors

import (
	"fmt"
	"path/filepath"

	"github.com/MusicDin/kubitect/embed"
	"github.com/MusicDin/kubitect/pkg/utils/template"
)

type Template[T any] struct {
	path   string
	values T
}

// NewTemplate initializes new embedded template on the given path
// with the provided values.
func NewTemplate[T any](templatePath string, values T) Template[T] {
	return Template[T]{
		path:   filepath.Clean(templatePath),
		values: values,
	}
}

func (t Template[T]) Name() string {
	return filepath.Base(t.path)
}

func (t Template[T]) Path() string {
	return t.path
}

func (t Template[T]) Values() T {
	return t.values
}

func (t Template[T]) Template() (string, error) {
	tpl, err := embed.GetTemplate(t.path + ".tpl")
	if err != nil {
		return "", err
	}

	return template.TrimTemplate(string(tpl.Content)), nil
}

// Write writes the populated template to the given path.
func (t Template[T]) Write(dstPath string) error {
	err := template.Write(t, filepath.Join(dstPath, t.Name()))
	if err != nil {
		return fmt.Errorf("Failed writing template %q: %v", t.Name(), err)
	}

	return nil
}
