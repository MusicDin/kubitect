package cluster

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/MusicDin/kubitect/pkg/app"
	"github.com/MusicDin/kubitect/pkg/utils/template"

	"github.com/stretchr/testify/assert"
)

func TestNewCluster(t *testing.T) {
	ctx := app.MockAppContext(t)

	c, err := NewCluster(ctx, ConfigMock{}.Write(t))
	assert.NoError(t, err)
	assert.Equal(t, "cluster-mock", c.Name)
}

func TestNewCluster_Local(t *testing.T) {
	ctx := app.MockAppContext(t, app.AppContextOptions{Local: true})

	c, err := NewCluster(ctx, ConfigMock{}.Write(t))
	assert.NoError(t, err)
	assert.Equal(t, "local-cluster-mock", c.Name)
}

func TestNewCluster_ClusterNameAlreadyPrefixed(t *testing.T) {
	ctx := app.MockAppContext(t, app.AppContextOptions{Local: true})
	cfg := ConfigMock{ClusterName: "local-cluster-mock"}
	cfgPath := cfg.Write(t)

	c, err := NewCluster(ctx, cfgPath)
	assert.NoError(t, err)
	assert.Equal(t, "local-cluster-mock", c.Name)
}

func TestNewCluster_InvalidClusterName(t *testing.T) {
	cfg := ConfigMock{ClusterName: "local-cluster-mock"}
	cfgPath := cfg.Write(t)

	_, err := NewCluster(app.MockAppContext(t), cfgPath)
	assert.ErrorContains(t, err, "cluster name contains the prefix 'local'")
}

func TestNewCluster_ConfigNotExists(t *testing.T) {
	_, err := NewCluster(app.MockAppContext(t), "config.yaml")
	assert.EqualError(t, err, "file 'config.yaml' does not exist")
}

func TestNewCluster_InvalidConfig(t *testing.T) {
	// Create empty configuration file
	cfgPath := path.Join(t.TempDir(), "config.yaml")
	_, err := os.Create(cfgPath)
	assert.NoError(t, err)

	_, err = NewCluster(app.MockAppContext(t), cfgPath)
	assert.ErrorContains(t, err, "invalid configuration file")
}

func TestSync_FailReadingAppliedConfig(t *testing.T) {
	c := MockCluster(t)

	// Make directory on path of applied config
	err := os.MkdirAll(c.AppliedConfigPath(), 0777)
	assert.NoError(t, err)

	assert.ErrorContains(t, c.Sync(), "failed to read previously applied configuration file")
}

func TestSync_FailReadingInfraConfig(t *testing.T) {
	c := MockCluster(t)

	// Make directory on path of applied config
	err := os.MkdirAll(c.InfrastructureConfigPath(), 0777)
	assert.NoError(t, err)

	assert.ErrorContains(t, c.Sync(), "failed to read infrastructure file")
}

func TestSync_InvalidInfraConfig(t *testing.T) {
	c := MockCluster(t)

	// Invalid infrastructure config
	cfg := template.TrimTemplate(`
		nodes:
			master:
				instances:
					- id: "1"
					  ip: "192.168.113.10"
					- id: "2"
					  ip: "192.168.113.10"
	`)

	err := os.MkdirAll(path.Dir(c.InfrastructureConfigPath()), 0777)
	assert.NoError(t, err)
	err = ioutil.WriteFile(c.InfrastructureConfigPath(), []byte(cfg), 0777)
	assert.NoError(t, err)

	assert.ErrorContains(t, c.Sync(), "infrastructure file (produced by Terraform) is invalid")
}

func TestApplyNewConfig(t *testing.T) {
	c := MockCluster(t)

	assert.NoFileExists(t, c.AppliedConfigPath())
	assert.NoError(t, c.ApplyNewConfig())
	assert.FileExists(t, c.AppliedConfigPath())
}

func TestStoreNewConfigFile(t *testing.T) {
	c := MockCluster(t)

	archiveFile := path.Join(c.Path, DefaultConfigDir, DefaultNewConfigFilename)
	assert.NoFileExists(t, archiveFile)
	assert.NoError(t, c.StoreNewConfig())
	assert.FileExists(t, archiveFile)
}
