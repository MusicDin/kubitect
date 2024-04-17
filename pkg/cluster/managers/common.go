package managers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MusicDin/kubitect/pkg/cluster/event"
	"github.com/MusicDin/kubitect/pkg/models/config"
	"github.com/MusicDin/kubitect/pkg/models/infra"
	"github.com/MusicDin/kubitect/pkg/tools/ansible"
	"github.com/MusicDin/kubitect/pkg/utils/exec"
)

type common struct {
	ClusterName       string
	ClusterPath       string
	SshPrivateKeyPath string
	ConfigDir         string
	CacheDir          string
	SharedDir         string
	Config            *config.Config
	InfraConfig       *infra.Config

	Ansible ansible.Ansible
}

func (e common) K8sVersion() string {
	return string(e.Config.Kubernetes.Version)
}

func (e common) SshUser() string {
	return string(e.Config.Cluster.NodeTemplate.User)
}

func (e common) SshPKey() string {
	return e.SshPrivateKeyPath
}

// mergeKubeconfig merges cluster kubeconfig with config default
// config in user directory (~/.kube/config). Note that if kubectl
// is not present locally, the command will fail.
func (e common) mergeKubeconfig() error {
	// Get directory of the current user.
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	defConfigPath := filepath.Join(home, ".kube", "config")
	clsConfigPath := filepath.Join(e.ConfigDir, "admin.conf")

	exec := exec.NewLocalClient()
	exec.SetEnv("KUBECONFIG", fmt.Sprintf("%s:%s", clsConfigPath, defConfigPath))

	config, err := exec.Output("kubectl", "config", "view", "--raw")
	if err != nil {
		return err
	}

	return os.WriteFile(defConfigPath, config, 0600)
}

// extractRemovedNodes returns removed node instances extracted from the event changes.
func extractRemovedNodes(events []event.Event) ([]config.Instance, error) {
	var nodes []config.Instance
	for _, e := range events {
		if e.Rule.ActionType != event.Action_ScaleDown {
			continue
		}

		node, ok := e.Change.ValueBefore.(config.Instance)
		if ok {
			nodes = append(nodes, node)
			continue
		}

		return nil, fmt.Errorf("%v cannot be scaled", e.Change.ValueType.Name())
	}

	return nodes, nil
}

// extractNewNodes returns new node instances extracted from the event changes.
func extractNewNodes(events []event.Event) ([]config.Instance, error) {
	var nodes []config.Instance
	for _, e := range events {
		if e.Rule.ActionType != event.Action_ScaleUp {
			continue
		}

		node, ok := e.Change.ValueAfter.(config.Instance)
		if ok {
			nodes = append(nodes, node)
			continue
		}

		return nil, fmt.Errorf("%v cannot be scaled", e.Change.ValueType.Name())
	}

	return nodes, nil
}
