package cluster

import (
	"cli/app"
	"cli/cluster/executors"
	"cli/cluster/provisioner"
	"cli/config/modelconfig"
	"cli/ui"
	"cli/utils/template"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	clusterMock struct {
		Cluster
		ui ui.UiMock
	}
)

// mockConfigFile writes a sample test configuration file and
// returns its path.
func mockConfigFile(t *testing.T) string {
	cfgPath := path.Join(t.TempDir(), "config.yaml")
	cfg := template.TrimTemplate(`
		hosts:
			- name: localhost
				connection:
					type: local

		cluster:
			name: cluster-mock
			network:
				cidr: 192.168.113.0/24
			nodes:
				master:
					instances:
						- id: 1

		kubernetes:
			version: v1.24.7
			kubespray:
				version: v2.21.0
	`)

	err := ioutil.WriteFile(cfgPath, []byte(cfg), 0777)
	assert.NoError(t, err, "Failed to write test configuration file!")

	return cfgPath
}

func MockCluster(t *testing.T) *clusterMock {
	t.Helper()

	ctxOptions := app.AppContextOptions{
		Local:       false,
		AutoApprove: true,
	}
	ctx := app.MockAppContext(t, ctxOptions)

	c, err := NewCluster(ctx, mockConfigFile(t))
	assert.NoError(t, err)

	// Create empty SSH keys
	keyDir := t.TempDir()
	keyPath := path.Join(keyDir, "key")
	os.Create(keyPath)
	os.Create(keyPath + ".pub")
	c.NewConfig.Cluster.NodeTemplate.SSH.PrivateKeyPath = modelconfig.File(keyPath)

	c.exec = executors.MockExecutor(t)
	c.prov = provisioner.MockProvisioner(t)

	return &clusterMock{*c, ctx.Ui()}
}

func MockLocalCluster(t *testing.T) *clusterMock {
	t.Helper()

	c := MockCluster(t)
	c.Local = true

	return c
}

func TestNewCluster(t *testing.T) {
	_, err := NewCluster(app.MockAppContext(t), mockConfigFile(t))
	assert.NoError(t, err)
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

	// Make invalid infrastructure config
	err := os.MkdirAll(path.Dir(c.InfrastructureConfigPath()), 0777)
	assert.NoError(t, err)
	err = ioutil.WriteFile(c.InfrastructureConfigPath(), []byte("cluster:"), 0777)
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
