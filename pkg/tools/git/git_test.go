package git

import (
	"testing"

	"github.com/MusicDin/kubitect/pkg/env"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGit(t *testing.T) {
	repo := NewGitRepo(env.ConstProjectUrl).WithRef(env.ConstProjectVersion)
	assert.Equal(t, env.ConstProjectUrl, repo.url)
	assert.Equal(t, env.ConstProjectVersion, repo.version)
}

func TestClone_Branch(t *testing.T) {
	repo := NewGitRepo(env.ConstProjectUrl).WithRef("main")
	require.NoError(t, repo.Clone(t.TempDir()))
}

func TestClone_Tag(t *testing.T) {
	repo := NewGitRepo(env.ConstProjectUrl).WithRef("v2.0.0")
	require.NoError(t, repo.Clone(t.TempDir()))
}

func TestClone_CommitHash(t *testing.T) {
	repo := NewGitRepo(env.ConstProjectUrl).WithCommitHash("c45d60ebc11e6925be8aebfaef1f6b025772c509")
	assert.NoError(t, repo.Clone(t.TempDir()))
}

func TestClone_HEAD(t *testing.T) {
	// HEAD = empty ref.
	repo := NewGitRepo(env.ConstProjectUrl)
	assert.NoError(t, repo.Clone(t.TempDir()))
}

func TestClone_EmptyURL(t *testing.T) {
	repo := NewGitRepo("").WithRef("master")
	assert.ErrorIs(t, repo.Clone(t.TempDir()), ErrInvalidRepositoryURL)
}

func TestClone_InvalidURL(t *testing.T) {
	repo := NewGitRepo(env.ConstProjectUrl + "invalid").WithRef("master")
	assert.ErrorContains(t, repo.Clone(t.TempDir()), "authentication required")
}

func TestClone_InvalidDestination(t *testing.T) {
	repo := NewGitRepo(env.ConstProjectUrl).WithRef("master")
	assert.ErrorContains(t, repo.Clone(""), ": no such file or directory")
}
