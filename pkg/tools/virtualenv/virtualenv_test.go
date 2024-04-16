package virtualenv

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/MusicDin/kubitect/pkg/ui"

	"github.com/stretchr/testify/require"
)

func MockReqFile(t *testing.T) string {
	reqFile := "netaddr==0.7.19"
	reqPath := path.Join(t.TempDir(), "requirements.txt")

	err := os.WriteFile(reqPath, []byte(reqFile), os.ModePerm)
	require.NoError(t, err)

	return reqPath
}

func MockVirtualEnv(t *testing.T) *VirtualEnv {
	return &VirtualEnv{
		path:             path.Join(t.TempDir(), "env"),
		requirementsPath: MockReqFile(t),
	}
}

func TestCreate(t *testing.T) {
	ui.MockGlobalUi(t, ui.UiOptions{Debug: true})

	env := MockVirtualEnv(t)

	require.NoError(t, env.create())
	require.Equal(t, "env", path.Base(env.path))
}

func TestInstallPipReq(t *testing.T) {
	ui.MockGlobalUi(t, ui.UiOptions{Debug: true})

	env := MockVirtualEnv(t)

	require.NoError(t, env.create())
	require.NoError(t, env.installPipReq())

	require.DirExists(t, filepath.Join(env.path, "bin"))
	require.DirExists(t, filepath.Join(env.path, "lib"))
}

func TestInit(t *testing.T) {
	env := NewVirtualEnv(t.TempDir(), MockReqFile(t))
	require.NoError(t, env.Init())
	require.NoError(t, env.Init()) // Instant, since environment already exists
}

func TestInit_InvalidReqPath(t *testing.T) {
	env := NewVirtualEnv(t.TempDir(), "")
	require.ErrorContains(t, env.Init(), "failed to install pip3 requirements:")
}
