package git

import (
	"cli/env"
	"cli/ui"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClone_Branch(t *testing.T) {
	p := &GitProject{
		Url:     env.ConstProjectUrl,
		Version: "master",
		Path:    t.TempDir(),
		Ui:      ui.MockUi(t),
	}

	assert.NoError(t, p.Clone())
}

func TestClone_Version(t *testing.T) {
	p := &GitProject{
		Url:     env.ConstProjectUrl,
		Version: "v2.0.0",
		Path:    t.TempDir(),
		Ui:      ui.MockUi(t),
	}

	p.Ui.Debug = true

	assert.NoError(t, p.Clone())
}

func TestClone_EmptyVersion(t *testing.T) {
	p := &GitProject{
		Url:     env.ConstProjectUrl,
		Version: "",
		Path:    t.TempDir(),
		Ui:      ui.MockUi(t),
	}

	assert.EqualError(t, p.Clone(), "git clone: project version not set")
}

func TestClone_EmptyURL(t *testing.T) {
	p := &GitProject{
		Url:     "",
		Version: "master",
		Path:    t.TempDir(),
		Ui:      ui.MockUi(t),
	}

	assert.EqualError(t, p.Clone(), "git clone: project URL not set")
}

func TestClone_InvalidURL(t *testing.T) {
	p := &GitProject{
		Url:     env.ConstProjectUrl + "wrong",
		Version: "master",
		Path:    t.TempDir(),
		Ui:      ui.MockUi(t),
	}

	assert.ErrorContains(t, p.Clone(), ": authentication required")
}

func TestClone_InvalidDestination(t *testing.T) {
	p := &GitProject{
		Url:     env.ConstProjectUrl,
		Version: "master",
		Path:    "",
		Ui:      ui.MockUi(t),
	}

	assert.ErrorContains(t, p.Clone(), ": no such file or directory")
}
