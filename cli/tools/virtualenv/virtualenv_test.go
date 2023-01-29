package virtualenv

import (
	"cli/ui"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockReqFile(t *testing.T) string {
	reqFile := "netaddr==0.7.19"
	reqPath := path.Join(t.TempDir(), "requirements.txt")

	err := ioutil.WriteFile(reqPath, []byte(reqFile), os.ModePerm)
	assert.NoError(t, err)

	return reqPath
}

func MockVirtualEnv(t *testing.T) *virtualEnv {
	tmpDir := t.TempDir()

	return &virtualEnv{
		name:             "mock",
		path:             tmpDir + "env",
		workingDir:       tmpDir,
		requirementsPath: MockReqFile(t),
	}
}

func TestCreate(t *testing.T) {
	ui.MockGlobalUi(t, ui.UiOptions{Debug: true})

	env := MockVirtualEnv(t)

	assert.NoError(t, env.create())
}

func TestInstallPipReq(t *testing.T) {
	ui.MockGlobalUi(t, ui.UiOptions{Debug: true})

	env := MockVirtualEnv(t)

	assert.NoError(t, env.create())
	assert.NoError(t, env.installPipReq())
}

func TestInit(t *testing.T) {
	tmpDir := t.TempDir()
	env := NewVirtualEnv("test", tmpDir, tmpDir, MockReqFile(t))

	assert.NoError(t, env.Init())
	assert.NoError(t, env.Init()) // Instant, since environment already exists
}

func TestInit_InvalidReqPath(t *testing.T) {
	tmpDir := t.TempDir()
	env := NewVirtualEnv("test", tmpDir, tmpDir, "")

	assert.ErrorContains(t, env.Init(), "failed to install pip3 requirements:")
}

func TestInit_InvalidWorkingDir(t *testing.T) {
	tmpDir := t.TempDir()
	env := NewVirtualEnv("test", tmpDir, tmpDir+"invalid", "")

	assert.ErrorContains(t, env.Init(), "failed to create virtual environment:")
}
