package git

import (
	"github.com/MusicDin/kubitect/cli/pkg/env"
	ui2 "github.com/MusicDin/kubitect/cli/pkg/ui"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGit(t *testing.T) {
	p := NewGitProject(env.ConstProjectUrl, env.ConstProjectVersion)
	assert.Equal(t, env.ConstProjectUrl, p.Url())
	assert.Equal(t, env.ConstProjectVersion, p.Version())
}

func TestClone_Branch(t *testing.T) {
	u := ui2.MockGlobalUi(t)
	p := NewGitProject(env.ConstProjectUrl, "master")

	assert.NoError(t, p.Clone(t.TempDir()))
	assert.Equal(t, "", u.ReadStdout(t))
}

func TestClone_Version(t *testing.T) {
	o := ui2.UiOptions{Debug: true}
	u := ui2.MockGlobalUi(t, o)
	p := NewGitProject(env.ConstProjectUrl, "v2.0.0")

	assert.NoError(t, p.Clone(t.TempDir()))
	assert.Contains(t, u.ReadStdout(t), "Compressing objects")
}

func TestClone_EmptyVersion(t *testing.T) {
	p := NewGitProject(env.ConstProjectUrl, "")
	assert.EqualError(t, p.Clone(t.TempDir()), "git clone: project version not set")
}

func TestClone_EmptyURL(t *testing.T) {
	p := NewGitProject("", "master")
	assert.EqualError(t, p.Clone(t.TempDir()), "git clone: project URL not set")
}

func TestClone_InvalidURL(t *testing.T) {
	p := NewGitProject(env.ConstProjectUrl+"invalid", "master")
	assert.ErrorContains(t, p.Clone(t.TempDir()), ": authentication required")
}

func TestClone_InvalidDestination(t *testing.T) {
	p := NewGitProject(env.ConstProjectUrl, "master")
	assert.ErrorContains(t, p.Clone(""), ": no such file or directory")
}
