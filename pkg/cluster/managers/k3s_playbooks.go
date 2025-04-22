package managers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/MusicDin/kubitect/pkg/tools/ansible"
)

// K3sCreate function calls an Ansible playbook that configures Kubernetes
// cluster.
func (e *k3s) K3sCreate(inventory string) error {
	// Use hashed cluster name as token. This is not perfect, as it makes
	// the token predictable, but removes the necessity to set it manually.
	hash := sha256.Sum256([]byte(e.ClusterName))
	token := hex.EncodeToString(hash[:])

	vars := map[string]string{
		"k3s_version":       fmt.Sprintf("%s+k3s1", e.K8sVersion()),
		"token":             string(token),
		"api_endpoint":      string(e.InfraConfig.Nodes.LoadBalancer.VIP),
		"api_port":          "6443",
		"user_kubectl":      "true", // Set to false to kubectl via root user.
		"cluster_context":   "default",
		"kubeconfig":        filepath.Join(e.ConfigDir, "admin.conf"),
		"extra_server_args": "",
		"extra_agent_args":  "",
	}

	pb := ansible.Playbook{
		WorkingDir: e.ProjectDir,
		Path:       filepath.Join(e.ProjectDir, "playbook/site.yml"),
		Inventory:  inventory,
		Become:     true,
		User:       e.SshUser(),
		PrivateKey: e.SshPKey(),
		Timeout:    600,
		ExtraVars:  vars,
	}

	return e.Ansible.Exec(pb)
}

// K3sUpdate function calls an Ansible playbook that configures Kubernetes
// cluster.
func (e *k3s) K3sUpgrade() error {
	vars := map[string]string{
		"k3s_version":       fmt.Sprintf("%s+k3s1", e.K8sVersion()),
		"api_endpoint":      string(e.InfraConfig.Nodes.LoadBalancer.VIP),
		"api_port":          "6443",
		"user_kubectl":      "true", // Set to false to kubectl via root user.
		"extra_server_args": "",
		"extra_agent_args":  "",
	}

	pb := ansible.Playbook{
		WorkingDir: e.ProjectDir,
		Path:       filepath.Join(e.ProjectDir, "playbook/upgrade.yml"),
		Inventory:  filepath.Join(e.ConfigDir, "nodes.yaml"),
		Become:     true,
		User:       e.SshUser(),
		PrivateKey: e.SshPKey(),
		Timeout:    600,
		ExtraVars:  vars,
	}

	return e.Ansible.Exec(pb)
}
