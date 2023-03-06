package app

import (
	"github.com/MusicDin/kubitect/pkg/env"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {
	ctx := MockAppContext(t, AppContextOptions{})
	assert.NotNil(t, ctx)
	assert.NotNil(t, ctx.Ui())
}

func TestNewAppContext(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)

	home, err := os.UserHomeDir()
	assert.NoError(t, err)

	ctx := AppContextOptions{}.AppContext()
	assert.Equal(t, wd, ctx.WorkingDir())
	assert.Equal(t, path.Join(home, defaultHomeDir), ctx.HomeDir())
	assert.Equal(t, path.Join(home, defaultHomeDir, defaultShareDir), ctx.ShareDir())
	assert.Equal(t, path.Join(home, defaultHomeDir, defaultClustersDir), ctx.ClustersDir())
	assert.Equal(t, path.Join(wd, defaultHomeDir, defaultClustersDir), ctx.LocalClustersDir())
	assert.False(t, ctx.Local())
	assert.False(t, ctx.ShowTerraformPlan())
}

func TestVerifyRequirements(t *testing.T) {
	assert.NoError(t, appContext{}.VerifyRequirements())
}

func TestVerifyRequirements_Missing(t *testing.T) {
	tmp := env.ProjectRequiredApps
	env.ProjectRequiredApps = append(env.ProjectRequiredApps, "invalid-app")

	err := appContext{}.VerifyRequirements()
	assert.EqualError(t, err, "Some requirements are not met: [invalid-app]")

	env.ProjectRequiredApps = tmp
}
