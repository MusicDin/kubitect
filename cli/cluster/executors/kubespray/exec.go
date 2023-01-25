package kubespray

import (
	"cli/cluster/event"
	"cli/cluster/executors"
	"cli/config/modelconfig"
	"cli/tools/virtualenv"
	"cli/ui"
	"fmt"
)

type VirtualEnvironments struct {
	MAIN      *virtualenv.VirtualEnv
	KUBESPRAY *virtualenv.VirtualEnv
}

type kubespray struct {
	ClusterName string
	ClusterPath string
	K8sVersion  string
	SshUser     string
	SshPKey     string

	Venvs VirtualEnvironments

	Ui *ui.Ui
}

func NewKubespray(
	clusterName string,
	clusterPath string,
	k8sVersion string,
	sshUser string,
	sshPKey string,

	venvs VirtualEnvironments,

	Ui *ui.Ui,
) executors.Executor {
	return &kubespray{
		ClusterName: clusterName,
		ClusterPath: clusterPath,
		K8sVersion:  k8sVersion,
		SshUser:     sshUser,
		SshPKey:     sshPKey,
		Venvs:       venvs,
		Ui:          Ui,
	}
}

func (e *kubespray) Init() error {
	if err := e.KubitectInit(TAG_INIT); err != nil {
		return err
	}

	return e.KubitectHostsSetup()
}

func (e *kubespray) Create() error {
	if err := e.KubitectInit(TAG_INIT, TAG_KUBESPRAY, TAG_GEN_NODES); err != nil {
		return err
	}

	if err := e.KubitectHostsSetup(); err != nil {
		return err
	}

	if err := e.HAProxy(); err != nil {
		return err
	}

	if err := e.KubesprayCreate(); err != nil {
		return err
	}

	return e.KubitectFinalize()
}

func (e *kubespray) Upgrade() error {
	if err := e.KubitectInit(TAG_INIT, TAG_KUBESPRAY, TAG_GEN_NODES); err != nil {
		return err
	}

	if err := e.KubitectHostsSetup(); err != nil {
		return err
	}

	if err := e.KubesprayUpgrade(); err != nil {
		return err
	}

	return e.KubitectFinalize()
}

// ScaleUp adds new nodes to the cluster.
func (e *kubespray) ScaleUp(events event.Events) error {
	events = events.OfType(event.SCALE_UP)

	if len(events) == 0 {
		return nil
	}

	if err := e.KubitectInit(TAG_KUBESPRAY, TAG_GEN_NODES); err != nil {
		return err
	}

	if err := e.HAProxy(); err != nil {
		return err
	}

	return e.KubesprayScale()
}

// scaleDown gracefully removes nodes from the cluster.
func (e *kubespray) ScaleDown(events event.Events) error {
	events = events.OfType(event.SCALE_DOWN)

	if len(events) == 0 {
		return nil
	}

	rmNodes, err := extractRemovedNodes(events)

	if err != nil {
		return err
	}

	if len(rmNodes) == 0 {
		return nil
	}

	var names []string

	for _, n := range rmNodes {
		name := fmt.Sprintf("%s-%s-%s", e.ClusterName, n.GetTypeName(), *n.GetID())
		names = append(names, name)
	}

	if err := e.KubitectInit(TAG_KUBESPRAY); err != nil {
		return err
	}

	return e.KubesprayRemoveNodes(names)
}

// extractRemovedNodes returns node instances from the event changes.
func extractRemovedNodes(events event.Events) ([]modelconfig.Instance, error) {
	var nodes []modelconfig.Instance

	for _, e := range events {
		for _, ch := range e.Changes() {
			if i, ok := ch.Before.(modelconfig.Instance); ok {
				nodes = append(nodes, i)
				continue
			}

			return nil, fmt.Errorf("%v cannot be scaled", ch.Type.Name())
		}
	}

	return nodes, nil
}
