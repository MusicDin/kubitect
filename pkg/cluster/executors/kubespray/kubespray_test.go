package kubespray

import (
	"fmt"
	"path"
	"reflect"
	"testing"

	"github.com/MusicDin/kubitect/pkg/cluster/event"
	"github.com/MusicDin/kubitect/pkg/cluster/executors"
	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/MusicDin/kubitect/pkg/models/config"
	"github.com/MusicDin/kubitect/pkg/models/infra"
	"github.com/MusicDin/kubitect/pkg/tools/ansible"
	"github.com/MusicDin/kubitect/pkg/utils/cmp"

	"github.com/stretchr/testify/assert"
)

type virtualEnvMock struct{}
type invalidVirtualEnvMock struct{ virtualEnvMock }

func (e *virtualEnvMock) Init() error        { return nil }
func (e *virtualEnvMock) Path() string       { return "" }
func (e *invalidVirtualEnvMock) Init() error { return fmt.Errorf("error") }

type ansibleMock struct{}
type invalidAnsibleMock struct{}

func (a *ansibleMock) Exec(ansible.Playbook) error        { return nil }
func (a *invalidAnsibleMock) Exec(ansible.Playbook) error { return fmt.Errorf("error") }

func MockExecutor(t *testing.T) *kubespray {
	tmpDir := t.TempDir()

	cfg := &config.Config{}
	cfg.Kubernetes.Version = env.ConstKubernetesVersion
	cfg.Cluster.NodeTemplate.User = config.User("test")
	cfg.Cluster.NodeTemplate.SSH.PrivateKeyPath = config.File(path.Join(tmpDir, ".ssh", "id_rsa"))

	iCfg := &infra.Config{}

	return &kubespray{
		ClusterName: "mock",
		ClusterPath: tmpDir,
		Config:      cfg,
		ConfigDir:   path.Join(tmpDir, "config"),
		InfraConfig: iCfg,
		VirtualEnv:  &virtualEnvMock{},
		Ansible:     &ansibleMock{},
	}
}

func MockInvalidExecutor(t *testing.T) executors.Executor {
	ks := MockExecutor(t)
	ks.VirtualEnv = &invalidVirtualEnvMock{}
	ks.Ansible = &invalidAnsibleMock{}
	return ks
}

func MockEvents(t *testing.T, obj interface{}, eType event.EventType) []event.Event {
	changes := []cmp.Change{
		{
			Type:   reflect.TypeOf(obj),
			Before: obj,
			After:  obj,
		},
	}

	return []event.Event{
		event.MockEvent(t, eType, cmp.MODIFY, changes),
	}
}

func TestNewExecutor(t *testing.T) {
	tmpDir := t.TempDir()
	clsName := "clsName"
	e := NewKubesprayExecutor(
		clsName,
		path.Join(tmpDir, clsName),
		path.Join(tmpDir, "id_rsa"),
		path.Join(tmpDir, "config"),
		path.Join(tmpDir, "cache"),
		path.Join(tmpDir, "share"),
		&config.Config{},
		&infra.Config{},
		&virtualEnvMock{},
	)
	assert.NotNil(t, e)
}

func TestInit(t *testing.T) {
	e := MockExecutor(t)
	assert.NoError(t, e.Init())
}

func TestInit_InvalidVenv(t *testing.T) {
	e := MockInvalidExecutor(t)
	assert.EqualError(t, e.Init(), "kubespray exec: initialize virtual environment: error")
}

func TestCreateAndUpgrade(t *testing.T) {
	e := MockExecutor(t)
	assert.NoError(t, e.Create())
	assert.NoError(t, e.Upgrade())
}

func TestCreateAndUpgrade_Invalid(t *testing.T) {
	e := MockInvalidExecutor(t)
	assert.EqualError(t, e.Create(), "error")
	assert.EqualError(t, e.Upgrade(), "error")
}

func TestExtractRemovedNodes(t *testing.T) {
	w := config.WorkerInstance{
		Id: "worker",
	}

	events := MockEvents(t, w, event.SCALE_DOWN)

	rmNodes, err := extractRemovedNodes(events)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []config.Instance{w}, rmNodes)
}

func TestScaleDown(t *testing.T) {
	w := config.WorkerInstance{
		Id: "worker",
	}

	events := MockEvents(t, w, event.SCALE_DOWN)

	err := MockExecutor(t).ScaleDown(events)
	assert.NoError(t, err)
}

func TestScaleDown_NoEvents(t *testing.T) {
	err := MockExecutor(t).ScaleDown(nil)
	assert.NoError(t, err)
}

func TestScaleDown_InvalidEvent(t *testing.T) {
	events := MockEvents(t, config.Host{}, event.SCALE_DOWN)

	err := MockExecutor(t).ScaleDown(events)
	assert.EqualError(t, err, "Host cannot be scaled")
}

func TestScaleUp(t *testing.T) {
	w := config.WorkerInstance{
		Id: "worker",
	}

	events := MockEvents(t, w, event.SCALE_UP)

	err := MockExecutor(t).ScaleUp(events)
	assert.NoError(t, err)
}

func TestScaleUp_NoEvents(t *testing.T) {
	err := MockExecutor(t).ScaleUp(nil)
	assert.NoError(t, err)
}
