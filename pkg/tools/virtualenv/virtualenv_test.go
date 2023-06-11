package virtualenv

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/MusicDin/kubitect/pkg/ui"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func MockReqFile(t *testing.T) string {
	reqFile := "netaddr==0.7.19"
	reqPath := path.Join(t.TempDir(), "requirements.txt")

	err := ioutil.WriteFile(reqPath, []byte(reqFile), os.ModePerm)
	require.NoError(t, err)

	return reqPath
}

func MockVirtualEnv(t *testing.T) *virtualEnv {
	tmpDir := t.TempDir()

	return &virtualEnv{
		path:             path.Join(tmpDir, "env"),
		workingDir:       tmpDir,
		requirementsPath: MockReqFile(t),
	}
}

func TestCreate(t *testing.T) {
	ui.MockGlobalUi(t, ui.UiOptions{Debug: true})

	env := MockVirtualEnv(t)

	assert.NoError(t, env.create())
	assert.Equal(t, "env", path.Base(env.Path()))
}

func TestInstallPipReq(t *testing.T) {
	ui.MockGlobalUi(t, ui.UiOptions{Debug: true})

	env := MockVirtualEnv(t)

	assert.NoError(t, env.create())
	assert.NoError(t, env.installPipReq())
}

func TestInit(t *testing.T) {
	tmpDir := t.TempDir()
	env := NewVirtualEnv(tmpDir, tmpDir, MockReqFile(t))

	assert.NoError(t, env.Init())
	assert.NoError(t, env.Init()) // Instant, since environment already exists
}

func TestInit_InvalidReqPath(t *testing.T) {
	tmpDir := t.TempDir()
	env := NewVirtualEnv(tmpDir, tmpDir, "")

	assert.ErrorContains(t, env.Init(), "failed to install pip3 requirements:")
}

func TestInit_InvalidWorkingDir(t *testing.T) {
	tmpDir := t.TempDir()
	env := NewVirtualEnv(tmpDir, tmpDir+"invalid", "")

	assert.ErrorContains(t, env.Init(), "failed to create virtual environment:")
}
