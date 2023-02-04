package kubespray

import (
	"cli/cluster/event"
	"cli/cluster/executors"
	"cli/config/modelconfig"
	"cli/config/modelinfra"
	"cli/tools/ansible"
	"cli/tools/virtualenv"
	"fmt"
	"path"
)

type kubespray struct {
	ClusterName       string
	ClusterPath       string
	SshPrivateKeyPath string
	Config            *modelconfig.Config
	InfraConfig       *modelinfra.Config
	VirtualEnv        virtualenv.VirtualEnv
	Ansible           ansible.Ansible
}

func (e *kubespray) K8sVersion() string {
	return string(e.Config.Kubernetes.Version)
}

func (e *kubespray) SshUser() string {
	return string(e.Config.Cluster.NodeTemplate.User)
}

func (e *kubespray) SshPKey() string {
	return e.SshPrivateKeyPath
}

func NewKubesprayExecutor(
	clusterName string,
	clusterPath string,
	sshPrivateKeyPath string,
	cfg *modelconfig.Config,
	infraCfg *modelinfra.Config,
	virtualEnv virtualenv.VirtualEnv,
) executors.Executor {
	return &kubespray{
		ClusterName:       clusterName,
		ClusterPath:       clusterPath,
		SshPrivateKeyPath: sshPrivateKeyPath,
		Config:            cfg,
		InfraConfig:       infraCfg,
		VirtualEnv:        virtualEnv,
	}
}

func (e *kubespray) Init() error {
	err := e.VirtualEnv.Init()
	if err != nil {
		return fmt.Errorf("kubespray exec: initialize virtual environment: %v", err)
	}

	if e.Ansible == nil {
		ansibleBinDir := path.Join(e.VirtualEnv.Path(), "bin")
		e.Ansible = ansible.NewAnsible(ansibleBinDir)
	}

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
	if err != nil || len(rmNodes) == 0 {
		return err
	}

	var names []string

	for _, n := range rmNodes {
		name := fmt.Sprintf("%s-%s-%s", e.ClusterName, n.GetTypeName(), n.GetID())
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
