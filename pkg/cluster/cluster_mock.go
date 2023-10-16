package cluster

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/MusicDin/kubitect/pkg/app"
	"github.com/MusicDin/kubitect/pkg/cluster/executors"
	"github.com/MusicDin/kubitect/pkg/cluster/provisioner"
	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/MusicDin/kubitect/pkg/models/config"
	"github.com/MusicDin/kubitect/pkg/ui"
	"github.com/MusicDin/kubitect/pkg/utils/defaults"
	"github.com/MusicDin/kubitect/pkg/utils/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	c, err := NewCluster(ctx, ConfigMock{}.Write(t))
	require.NoError(t, err)

	// Create an empty SSH key pair
	keyDir := t.TempDir()
	keyPath := path.Join(keyDir, "key")
	os.Create(keyPath)
	os.Create(keyPath + ".pub")
	c.NewConfig.Cluster.NodeTemplate.SSH.PrivateKeyPath = config.File(keyPath)

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

type ConfigMock struct {
	ClusterName string
}

func (c ConfigMock) Name() string {
	return c.ClusterName
}

func (c *ConfigMock) SetDefaults() {
	c.ClusterName = defaults.Default(c.ClusterName, "cluster-mock")
}

func (c ConfigMock) Template() string {
	return template.TrimTemplate(fmt.Sprintf(`
		hosts:
			- name: localhost
				connection:
					type: local

		cluster:
			name: {{ .ClusterName }}
			network:
				cidr: 192.168.113.0/24
			nodes:
				master:
					instances:
						- id: 1

		kubernetes:
			version: %s
	`, env.ConstKubernetesVersion))
}

func (c ConfigMock) Write(t *testing.T) string {
	cfgPath := path.Join(t.TempDir(), "config.yaml")
	c.SetDefaults()

	err := template.Write(c, cfgPath)
	assert.NoError(t, err)

	return cfgPath
}
