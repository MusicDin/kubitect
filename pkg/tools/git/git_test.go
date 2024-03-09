package git

import (
	"testing"

	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/MusicDin/kubitect/pkg/ui"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGit(t *testing.T) {
	p := NewGitProject(env.ConstProjectUrl, env.ConstProjectVersion)
	assert.Equal(t, env.ConstProjectUrl, p.Url())
	assert.Equal(t, env.ConstProjectVersion, p.Version())
}

func TestClone_Branch(t *testing.T) {
	u := ui.MockGlobalUi(t)
	p := NewGitProject(env.ConstProjectUrl, "main")

	require.NoError(t, p.Clone(t.TempDir()))
	assert.Equal(t, "", u.ReadStdout(t))
}

func TestClone_Version(t *testing.T) {
	o := ui.UiOptions{Debug: true}
	u := ui.MockGlobalUi(t, o)
	p := NewGitProject(env.ConstProjectUrl, "v2.0.0")

	require.NoError(t, p.Clone(t.TempDir()))
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
