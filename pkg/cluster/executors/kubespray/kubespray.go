package kubespray

import (
	"fmt"
	"os"
	"path"

	"github.com/MusicDin/kubitect/pkg/cluster/event"
	"github.com/MusicDin/kubitect/pkg/cluster/executors"
	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/MusicDin/kubitect/pkg/models/config"
	"github.com/MusicDin/kubitect/pkg/models/infra"
	"github.com/MusicDin/kubitect/pkg/tools/ansible"
	"github.com/MusicDin/kubitect/pkg/tools/git"
	"github.com/MusicDin/kubitect/pkg/tools/virtualenv"
	"github.com/MusicDin/kubitect/pkg/ui"
	"github.com/MusicDin/kubitect/pkg/utils/file"
	"gopkg.in/yaml.v3"
)

type kubespray struct {
	ClusterName       string
	ClusterPath       string
	SshPrivateKeyPath string
	ConfigDir         string
	CacheDir          string
	SharedDir         string
	Config            *config.Config
	InfraConfig       *infra.Config
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
	configDir string,
	cacheDir string,
	sharedDir string,
	cfg *config.Config,
	infraCfg *infra.Config,
	virtualEnv virtualenv.VirtualEnv,
) executors.Executor {
	return &kubespray{
		ClusterName:       clusterName,
		ClusterPath:       clusterPath,
		SshPrivateKeyPath: sshPrivateKeyPath,
		ConfigDir:         configDir,
		CacheDir:          cacheDir,
		SharedDir:         sharedDir,
		Config:            cfg,
		InfraConfig:       infraCfg,
		VirtualEnv:        virtualEnv,
	}
}

// Init clones Kubespray project, initializes virtual environment
// and generates Ansible hosts inventory.
func (e *kubespray) Init() error {
	url := env.ConstKubesprayUrl
	ver := env.ConstKubesprayVersion
	dst := path.Join(e.ClusterPath, "ansible", "kubespray")

	if err := os.RemoveAll(dst); err != nil {
		return err
	}

	ui.Printf(ui.INFO, "Cloning Kubespray (%s)...\n", ver)
	if err := git.NewGitProject(url, ver).Clone(dst); err != nil {
		return err
	}

	// Patch: There is an issue with unsafe conditionals in
	// Kubespray with ansible-core version > 2.14.11.
	reqPath := path.Join(dst, "requirements.txt")
	reqPatch := []byte("ansible-core==2.14.11")
	if err := file.Append(reqPath, reqPatch); err != nil {
		return err
	}

	if err := e.VirtualEnv.Init(); err != nil {
		return fmt.Errorf("kubespray exec: initialize virtual environment: %v", err)
	}

	if e.Ansible == nil {
		ansibleBinDir := path.Join(e.VirtualEnv.Path(), "bin")
		e.Ansible = ansible.NewAnsible(ansibleBinDir, e.CacheDir)
	}

	return e.KubitectHostsSetup()
}

// Sync regenerates required Ansible inventories and Kubespray group
// variables.
func (e *kubespray) Sync() error {
	if err := e.generateHostsInventory(); err != nil {
		return err
	}

	if err := e.generateNodesInventory(); err != nil {
		return err
	}

	return e.generateGroupVars()
}

// Create creates a Kubernetes cluster by calling appropriate Kubespray
// playbooks.
func (e *kubespray) Create() error {
	if err := e.HAProxy(); err != nil {
		return err
	}

	if err := e.KubesprayCreate(); err != nil {
		return err
	}

	return e.KubitectFinalize()
}

// Upgrades upgrades a Kubernetes cluster by calling appropriate Kubespray
// playbooks.
func (e *kubespray) Upgrade() error {
	if err := e.KubesprayUpgrade(); err != nil {
		return err
	}

	return e.KubitectFinalize()
}

// ScaleUp adds new nodes to the cluster.
func (e *kubespray) ScaleUp(events event.Events) error {
	events = events.FilterByAction(event.Action_ScaleUp)

	if len(events) == 0 {
		return nil
	}

	if err := e.HAProxy(); err != nil {
		return err
	}

	return e.KubesprayScale()
}

// ScaleDown gracefully removes nodes from the cluster.
func (e *kubespray) ScaleDown(events event.Events) error {
	events = events.FilterByAction(event.Action_ScaleDown)

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

	if err := e.generateGroupVars(); err != nil {
		return err
	}

	if err := e.KubesprayRemoveNodes(names); err != nil {
		return err
	}

	return e.generateNodesInventory()
}

// extractRemovedNodes returns node instances from the event changes.
func extractRemovedNodes(events []event.Event) ([]config.Instance, error) {
	var nodes []config.Instance

	for _, e := range events {
		if i, ok := e.Change.ValueBefore.(config.Instance); ok {
			nodes = append(nodes, i)
			continue
		}

		return nil, fmt.Errorf("%v cannot be scaled", e.Change.ValueType.Name())
	}

	return nodes, nil
}

// generateNodesInventory creates an Ansible inventory of target nodes.
func (e *kubespray) generateNodesInventory() error {
	return NewNodesTemplate(e.ConfigDir, e.Config.Cluster.Nodes, e.InfraConfig.Nodes).Write()
}

// generateHostsInventory creates an Ansible inventory of target hosts.
func (e *kubespray) generateHostsInventory() error {
	return NewHostsTemplate(e.ConfigDir, e.Config.Hosts).Write()
}

// generateGroupVars creates a directory of Kubespray group variables.
func (e *kubespray) generateGroupVars() error {
	err := NewKubesprayAllTemplate(e.ConfigDir, e.InfraConfig.Nodes).Write()
	if err != nil {
		return err
	}

	err = NewKubesprayK8sClusterTemplate(e.ConfigDir, *e.Config).Write()
	if err != nil {
		return err
	}

	addons, err := yaml.Marshal(e.Config.Addons.Kubespray)
	if err != nil {
		return err
	}

	err = NewKubesprayAddonsTemplate(e.ConfigDir, string(addons)).Write()
	if err != nil {
		return err
	}

	return NewKubesprayEtcdTemplate(e.ConfigDir).Write()
}
