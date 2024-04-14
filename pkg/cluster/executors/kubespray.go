package executors

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/MusicDin/kubitect/pkg/cluster/event"
	"github.com/MusicDin/kubitect/pkg/cluster/interfaces"
	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/MusicDin/kubitect/pkg/models/config"
	"github.com/MusicDin/kubitect/pkg/models/infra"
	"github.com/MusicDin/kubitect/pkg/tools/ansible"
	"github.com/MusicDin/kubitect/pkg/tools/git"
	"github.com/MusicDin/kubitect/pkg/tools/virtualenv"
	"github.com/MusicDin/kubitect/pkg/ui"
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
) interfaces.Executor {
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
	err := os.RemoveAll(dst)
	if err != nil {
		return err
	}

	ui.Printf(ui.INFO, "Cloning Kubespray (%s)...\n", ver)

	err = git.NewGitRepo(url).WithRef(ver).Clone(dst)
	if err != nil {
		return err
	}

	err = e.VirtualEnv.Init()
	if err != nil {
		return fmt.Errorf("kubespray exec: initialize virtual environment: %v", err)
	}

	if e.Ansible == nil {
		ansibleBinDir := path.Join(e.VirtualEnv.Path(), "bin")
		e.Ansible = ansible.NewAnsible(ansibleBinDir, e.CacheDir)
	}

	return nil
}

// Sync regenerates required Ansible inventories and Kubespray group
// variables.
func (e *kubespray) Sync() error {
	err := e.generateInventory()
	if err != nil {
		return err
	}

	return e.generateGroupVars()
}

// Create creates a Kubernetes cluster by calling appropriate Kubespray
// playbooks.
func (e *kubespray) Create() error {
	err := e.HAProxy()
	if err != nil {
		return err
	}

	err = e.KubesprayCreate()
	if err != nil {
		return err
	}

	return e.KubitectFinalize()
}

// Upgrades upgrades a Kubernetes cluster by calling appropriate Kubespray
// playbooks.
func (e *kubespray) Upgrade() error {
	err := e.KubesprayUpgrade()
	if err != nil {
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

	err := e.HAProxy()
	if err != nil {
		return err
	}

	return e.KubesprayScale()
}

// ScaleDown gracefully removes nodes from the cluster.
func (e *kubespray) ScaleDown(events event.Events) error {
	rmNodes, err := extractRemovedNodes(events)
	if err != nil {
		return err
	}

	if len(rmNodes) == 0 {
		// No removed nodes.
		return nil
	}

	var names []string
	for _, n := range rmNodes {
		name := fmt.Sprintf("%s-%s-%s", e.ClusterName, n.GetTypeName(), n.GetID())
		names = append(names, name)
	}

	err = e.generateGroupVars()
	if err != nil {
		return err
	}

	err = e.KubesprayRemoveNodes(names)
	if err != nil {
		return err
	}

	return e.generateInventory()
}

// generateInventory creates an Ansible inventory containing cluster nodes.
func (e *kubespray) generateInventory() error {
	nodes := struct {
		ConfigNodes config.Nodes
		InfraNodes  config.Nodes
	}{
		ConfigNodes: e.Config.Cluster.Nodes,
		InfraNodes:  e.InfraConfig.Nodes,
	}

	return NewTemplate("kubespray/inventory.yaml", nodes).Write(filepath.Join(e.ConfigDir, "nodes.yaml"))
}

// generateGroupVars creates a directory of Kubespray group variables.
func (e *kubespray) generateGroupVars() error {
	groupVarsDir := filepath.Join(e.ConfigDir, "group_vars")

	err := NewTemplate("kubespray/all.yaml", e.InfraConfig.Nodes).Write(filepath.Join(groupVarsDir, "all", "all.yml"))
	if err != nil {
		return err
	}

	err = NewTemplate("kubespray/k8s-cluster.yaml", *e.Config).Write(filepath.Join(groupVarsDir, "k8s_cluster", "k8s-cluster.yaml"))
	if err != nil {
		return err
	}

	addons, err := yaml.Marshal(e.Config.Addons.Kubespray)
	if err != nil {
		return err
	}

	addonsPath := filepath.Join(groupVarsDir, "k8s_cluster", "addons.yaml")
	err = os.WriteFile(addonsPath, addons, 0644)
	if err != nil {
		return err
	}

	err = NewTemplate("kubespray/etcd.yaml", "").Write(filepath.Join(groupVarsDir, "etcd.yaml"))
	if err != nil {
		return err
	}

	return nil
}
