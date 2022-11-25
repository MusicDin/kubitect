package git

import (
	"cli/env"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClone_EmptyVersion(t *testing.T) {
	err := Clone(t.TempDir(), env.ConstProjectUrl, "")
	assert.EqualError(t, err, "git clone: version not provided")
}

func TestClone_EmptyURL(t *testing.T) {
	err := Clone(t.TempDir(), "", "master")
	assert.EqualError(t, err, "git clone: URL not provided")
}

func TestClone_InvalidURL(t *testing.T) {
	err := Clone(t.TempDir(), env.ConstProjectUrl+"wrong", "master")
	assert.ErrorContains(t, err, ": authentication required")
}

func TestClone_InvalidDestination(t *testing.T) {
	err := Clone("", env.ConstProjectUrl, "master")
	assert.ErrorContains(t, err, ": no such file or directory")
}

func TestClone_Branch(t *testing.T) {
	err := Clone(t.TempDir(), env.ConstProjectUrl, "master")
	assert.NoError(t, err)
}

func TestClone_Version(t *testing.T) {
	err := Clone(t.TempDir(), env.ConstProjectUrl, "v2.0.0")
	assert.NoError(t, err)
}
