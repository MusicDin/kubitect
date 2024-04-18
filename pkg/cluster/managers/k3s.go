package managers

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/MusicDin/kubitect/pkg/cluster/event"
	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/MusicDin/kubitect/pkg/models/config"
	"github.com/MusicDin/kubitect/pkg/models/infra"
	"github.com/MusicDin/kubitect/pkg/tools/ansible"
	"github.com/MusicDin/kubitect/pkg/tools/git"
	"github.com/MusicDin/kubitect/pkg/tools/virtualenv"
	"github.com/MusicDin/kubitect/pkg/ui"
	"github.com/MusicDin/kubitect/pkg/utils/exec"
)

type k3s struct {
	common

	ProjectDir string
}

func (e *k3s) K8sVersion() string {
	return string(e.Config.Kubernetes.Version)
}

func (e *k3s) SshUser() string {
	return string(e.Config.Cluster.NodeTemplate.User)
}

func (e *k3s) SshPKey() string {
	return e.SshPrivateKeyPath
}

func NewK3sManager(
	clusterName string,
	clusterPath string,
	sshPrivateKeyPath string,
	configDir string,
	cacheDir string,
	sharedDir string,
	cfg *config.Config,
	infraCfg *infra.Config,
) *k3s {
	return &k3s{
		common: common{
			ClusterName:       clusterName,
			ClusterPath:       clusterPath,
			SshPrivateKeyPath: sshPrivateKeyPath,
			ConfigDir:         configDir,
			CacheDir:          cacheDir,
			SharedDir:         sharedDir,
			Config:            cfg,
			InfraConfig:       infraCfg,
		},
		ProjectDir: filepath.Join(clusterPath, "ansible", "k3s"),
	}
}

// Init clones k3s project, initializes virtual environment
// and generates Ansible hosts inventory.
func (e *k3s) Init() error {
	err := os.RemoveAll(e.ProjectDir)
	if err != nil {
		return err
	}

	// Clone repository with k3s playbooks.
	url := env.ConstK3sURL
	commitHash := env.ConstK3sVersion
	err = git.NewGitRepo(url).WithCommitHash(commitHash).Clone(e.ProjectDir)
	if err != nil {
		return err
	}

	if e.Ansible == nil {
		// Virtual environment.
		reqPath := filepath.Join(e.ClusterPath, "ansible/kubitect/requirements.txt")
		venvPath := path.Join(e.SharedDir, "venv", "k3s", env.ConstK3sVersion)
		err = virtualenv.NewVirtualEnv(venvPath, reqPath).Init()
		if err != nil {
			return fmt.Errorf("k3s: initialize virtual environment: %v", err)
		}

		ansibleBinDir := path.Join(venvPath, "bin")
		e.Ansible = ansible.NewAnsible(ansibleBinDir, e.CacheDir)
	}

	return nil
}

// Sync regenerates Ansible inventory.
func (e *k3s) Sync() error {
	nodes := struct {
		ConfigNodes config.Nodes
		InfraNodes  config.Nodes
	}{
		ConfigNodes: e.Config.Cluster.Nodes,
		InfraNodes:  e.InfraConfig.Nodes,
	}

	return NewTemplate("k3s/inventory.yaml", nodes).Write(filepath.Join(e.ConfigDir, "nodes.yaml"))
}

// Create creates a Kubernetes cluster by calling appropriate k3s
// playbooks.
func (e *k3s) Create() error {
	if err := e.HAProxy(); err != nil {
		return err
	}

	inventory := filepath.Join(e.ConfigDir, "nodes.yaml")
	err := e.K3sCreate(inventory)
	if err != nil {
		return err
	}

	err = e.Finalize()
	if err != nil {
		return err
	}

	if e.Config.Kubernetes.Other.MergeKubeconfig {
		err := e.mergeKubeconfig()
		if err != nil {
			// Just warn about failure, since deployment has succeeded.
			ui.Print(ui.WARN, "Failed to merge kubeconfig:", err)
		}
	}

	return nil
}

// Upgrades upgrades a Kubernetes cluster by calling appropriate k3s
// playbooks.
func (e *k3s) Upgrade() error {
	err := e.K3sUpgrade()
	if err != nil {
		return err
	}

	return e.Finalize()
}

// ScaleUp adds new nodes to the cluster.
func (e *k3s) ScaleUp(events event.Events) error {
	newNodes, err := extractNewNodes(events)
	if err != nil {
		return err
	}

	if len(newNodes) == 0 {
		// No removed nodes.
		return nil
	}

	nodes := make(map[string]config.Instance, len(newNodes))
	for _, n := range newNodes {
		name := fmt.Sprintf("%s-%s-%s", e.ClusterName, n.GetTypeName(), n.GetID())
		nodes[name] = n
	}

	inventory := filepath.Join(e.ConfigDir, "nodes_tmp.yaml")
	err = NewTemplate("k3s/inventory_partial.yaml", nodes).Write(inventory)
	if err != nil {
		return err
	}

	defer os.Remove(inventory)
	return e.K3sCreate(inventory)
}

// ScaleDown gracefully removes nodes from the cluster.
func (e *k3s) ScaleDown(events event.Events) error {
	rmNodes, err := extractRemovedNodes(events)
	if err != nil {
		return err
	}

	if len(rmNodes) == 0 {
		// No removed nodes.
		return nil
	}

	// Establish connection with one of the master nodes.
	leader := e.Config.Cluster.Nodes.Master.Instances[0]
	ssh := exec.NewSSHClient(e.SshUser(), string(leader.IP)).
		WithPrivateKeyFile(e.SshPKey()).
		WithSuperUser(true)

	ssh.SetCombinedStdout(os.Stdout)

	defer ssh.Close()

	for _, n := range rmNodes {
		name := fmt.Sprintf("%s-%s-%s", e.ClusterName, n.GetTypeName(), n.GetID())

		err = ssh.Run("kubectl", "cordon", name)
		if err != nil {
			return fmt.Errorf("cordon node %q: %v", name, err)
		}

		err = ssh.Run("kubectl", "drain", name, "--ignore-daemonsets", "--force")
		if err != nil {
			return fmt.Errorf("drain node %q: %v", name, err)
		}

		err = ssh.Run("kubectl", "delete", "node", name)
		if err != nil {
			return fmt.Errorf("delete node %q: %v", name, err)
		}
	}

	// No need for further cleanup. This instance will be removed.
	return nil
}
