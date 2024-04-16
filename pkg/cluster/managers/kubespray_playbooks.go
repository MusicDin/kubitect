package managers

import (
	"path/filepath"
	"strings"

	"github.com/MusicDin/kubitect/pkg/tools/ansible"
)

// KubesprayCreate function calls an Ansible playbook that configures Kubernetes
// cluster.
func (e *kubespray) KubesprayCreate() error {
	vars := map[string]string{
		"kube_version": e.K8sVersion(),
	}

	pb := ansible.Playbook{
		Path:       filepath.Join(e.ClusterPath, "ansible/kubespray/cluster.yml"),
		Inventory:  filepath.Join(e.ClusterPath, "config/nodes.yaml"),
		Become:     true,
		User:       e.SshUser(),
		PrivateKey: e.SshPKey(),
		Timeout:    3000,
		ExtraVars:  vars,
	}

	return e.Ansible.Exec(pb)
}

// KubesprayUpgrade function calls an Ansible playbook that upgrades Kubernetes
// nodes to a newer version.
func (e *kubespray) KubesprayUpgrade() error {
	vars := map[string]string{
		"kube_version": e.K8sVersion(),
	}

	pb := ansible.Playbook{
		Path:       filepath.Join(e.ClusterPath, "ansible/kubespray/upgrade-cluster.yml"),
		Inventory:  filepath.Join(e.ClusterPath, "config/nodes.yaml"),
		Become:     true,
		User:       e.SshUser(),
		PrivateKey: e.SshPKey(),
		Timeout:    3000,
		ExtraVars:  vars,
	}

	return e.Ansible.Exec(pb)
}

// KubesprayScale function calls an Ansible playbook that configures virtual machines
// that are freshly added to the cluster.
func (e *kubespray) KubesprayScale() error {
	vars := map[string]string{
		"kube_version": e.K8sVersion(),
	}

	pb := ansible.Playbook{
		Path:       filepath.Join(e.ClusterPath, "ansible/kubespray/scale.yml"),
		Inventory:  filepath.Join(e.ClusterPath, "config/nodes.yaml"),
		Become:     true,
		User:       e.SshUser(),
		PrivateKey: e.SshPKey(),
		Timeout:    3000,
		ExtraVars:  vars,
	}

	return e.Ansible.Exec(pb)
}

// KubesprayRemoveNodes function calls an Ansible playbook that removes the nodes with
// the provided names.
func (e *kubespray) KubesprayRemoveNodes(removedNodeNames []string) error {
	vars := map[string]string{
		"skip_confirmation":         "yes",
		"delete_nodes_confirmation": "yes",
		"node":                      strings.Join(removedNodeNames, ","),
	}

	pb := ansible.Playbook{
		Path:       filepath.Join(e.ClusterPath, "ansible/kubespray/remove-node.yml"),
		Inventory:  filepath.Join(e.ClusterPath, "config/nodes.yaml"),
		Become:     true,
		User:       e.SshUser(),
		PrivateKey: e.SshPKey(),
		Timeout:    3000,
		ExtraVars:  vars,
	}

	return e.Ansible.Exec(pb)
}
