package cluster

import (
	"github.com/MusicDin/kubitect/cli/pkg/app"
	"github.com/MusicDin/kubitect/cli/pkg/cluster/executors"
	"github.com/MusicDin/kubitect/cli/pkg/cluster/provisioner"
	"github.com/MusicDin/kubitect/cli/pkg/config/modelconfig"
	"github.com/MusicDin/kubitect/cli/pkg/ui"
	"github.com/MusicDin/kubitect/cli/pkg/utils/template"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	ClusterMock struct {
		*Cluster
		appContextMock app.AppContextMock
	}
)

func (m *ClusterMock) AppContext() app.AppContextMock {
	return m.appContextMock
}

func (m *ClusterMock) Ui() ui.UiMock {
	return m.appContextMock.Ui()
}

func MockCluster(t *testing.T) *ClusterMock {
	t.Helper()

	ctxOptions := app.AppContextOptions{
		Local:       false,
		AutoApprove: true,
	}
	ctx := app.MockAppContext(t, ctxOptions)

	c, err := NewCluster(ctx, mockConfigFile(t))
	assert.NoError(t, err)

	// Create an empty SSH key pair
	keyDir := t.TempDir()
	keyPath := path.Join(keyDir, "key")
	os.Create(keyPath)
	os.Create(keyPath + ".pub")
	c.NewConfig.Cluster.NodeTemplate.SSH.PrivateKeyPath = modelconfig.File(keyPath)

	c.exec = executors.MockExecutor(t)
	c.prov = provisioner.MockProvisioner(t)

	return &ClusterMock{c, ctx}
}

func MockLocalCluster(t *testing.T) *ClusterMock {
	t.Helper()

	c := MockCluster(t)
	c.Local = true

	return c
}

// mockConfigFile writes a sample test configuration file and
// returns its path.
func mockConfigFile(t *testing.T) string {
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
	return writeConfigFile(t, cfg)
}

func writeConfigFile(t *testing.T, cfg string) string {
	cfgPath := path.Join(t.TempDir(), "config.yaml")
	err := ioutil.WriteFile(cfgPath, []byte(cfg), 0777)
	assert.NoError(t, err, "Failed to write configuration file!")
	return cfgPath
}
