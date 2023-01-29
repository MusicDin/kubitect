package kubespray

import (
	"cli/cluster/event"
	"cli/cluster/executors"
	"cli/config/modelconfig"
	"cli/tools/ansible"
	"cli/tools/virtualenv"
	"cli/utils/cmp"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type virtualEnvMock struct{}

func (e *virtualEnvMock) Init() error {
	return nil
}

func (e *virtualEnvMock) Path() string {
	return ""
}

func MockVenv(t *testing.T) virtualenv.VirtualEnv {
	return &virtualEnvMock{}
}

type invalidVirtualEnvMock struct {
	virtualEnvMock
}

func (e *invalidVirtualEnvMock) Init() error {
	return fmt.Errorf("error")
}

func MockInvalidVenv(t *testing.T) virtualenv.VirtualEnv {
	return &invalidVirtualEnvMock{}
}

type ansibleMock struct{}

func (a *ansibleMock) Exec(ansible.Playbook) error {
	return nil
}

type invalidAnsibleMock struct{}

func (a *invalidAnsibleMock) Exec(ansible.Playbook) error {
	return fmt.Errorf("error")
}

func MockExecutor(t *testing.T) executors.Executor {
	tmpDir := t.TempDir()

	return &kubespray{
		ClusterName: "mock",
		ClusterPath: tmpDir,
		K8sVersion:  "v1.23.0",
		SshUser:     "test",
		SshPKey:     tmpDir,
		Ansible:     &ansibleMock{},
	}
}

func MockInvalidExecutor(t *testing.T) executors.Executor {
	tmpDir := t.TempDir()

	return &kubespray{
		ClusterName: "mock",
		ClusterPath: tmpDir,
		K8sVersion:  "v1.23.0",
		SshUser:     "test",
		SshPKey:     tmpDir,
		Ansible:     &invalidAnsibleMock{},
	}
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
		event.MockEvent(t, eType, changes),
	}
}

func TestNewExecutor(t *testing.T) {
	_, err := NewKubesprayExecutor("clsName", "clsPath", "vX.Y.Z", "user", "pk", MockVenv(t))
	assert.NoError(t, err)
}

func TestNewExecutor_InvalidVenv(t *testing.T) {
	_, err := NewKubesprayExecutor("clsName", "clsPath", "vX.Y.Z", "user", "pk", MockInvalidVenv(t))
	assert.EqualError(t, err, "kubespray exec: initialize virtual environment: error")
}

func TestActions(t *testing.T) {
	e := MockExecutor(t)
	assert.NoError(t, e.Init())
	assert.NoError(t, e.Create())
	assert.NoError(t, e.Upgrade())
}

func TestActions_Invalid(t *testing.T) {
	e := MockInvalidExecutor(t)
	assert.EqualError(t, e.Init(), "error")
	assert.EqualError(t, e.Create(), "error")
	assert.EqualError(t, e.Upgrade(), "error")
}

func TestExtractRemovedNodes(t *testing.T) {
	wName := "worker"
	w := modelconfig.WorkerInstance{
		Id: &wName,
	}

	events := MockEvents(t, w, event.SCALE_DOWN)

	rmNodes, err := extractRemovedNodes(events)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []modelconfig.Instance{w}, rmNodes)
}

func TestScaleDown(t *testing.T) {
	wName := "worker"
	w := modelconfig.WorkerInstance{
		Id: &wName,
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
	events := MockEvents(t, modelconfig.Host{}, event.SCALE_DOWN)

	err := MockExecutor(t).ScaleDown(events)
	assert.EqualError(t, err, "Host cannot be scaled")
}

func TestScaleUp(t *testing.T) {
	wName := "worker"
	w := modelconfig.WorkerInstance{
		Id: &wName,
	}

	events := MockEvents(t, w, event.SCALE_UP)

	err := MockExecutor(t).ScaleUp(events)
	assert.NoError(t, err)
}

func TestScaleUp_NoEvents(t *testing.T) {
	err := MockExecutor(t).ScaleUp(nil)
	assert.NoError(t, err)
}
