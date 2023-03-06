package kubespray

import (
	"path/filepath"
	"strings"

	"github.com/MusicDin/kubitect/pkg/tools/ansible"
)

type PlaybookTag string

const (
	TAG_INIT      PlaybookTag = "init"
	TAG_KUBESPRAY PlaybookTag = "kubespray"
	TAG_GEN_NODES PlaybookTag = "gen_nodes"
)

// KubitectHostsSetup function calls an Ansible playbook that ensures Kubitect target
// hosts meet all the requirements before cluster is created.
func (e *kubespray) KubitectHostsSetup() error {
	pb := ansible.Playbook{
		Path:      filepath.Join(e.ClusterPath, "ansible/kubitect/hosts-setup.yaml"),
		Inventory: filepath.Join(e.ClusterPath, "config/hosts.yaml"),
		Local:     true,
	}

	return e.Ansible.Exec(pb)
}

// KubitectFinalize function calls an Ansible playbook that finalizes Kubernetes
// cluster installation.
func (e *kubespray) KubitectFinalize() error {
	pb := ansible.Playbook{
		Path:       filepath.Join(e.ClusterPath, "ansible/kubitect/finalize.yaml"),
		Inventory:  filepath.Join(e.ClusterPath, "config/nodes.yaml"),
		Become:     true,
		User:       e.SshUser(),
		PrivateKey: e.SshPKey(),
		Timeout:    3000,
	}

	return e.Ansible.Exec(pb)
}

// HAProxy function calls an Ansible playbook that configures HAProxy
// load balancers.
func (e *kubespray) HAProxy() error {
	pb := ansible.Playbook{
		Path:       filepath.Join(e.ClusterPath, "ansible/kubitect/haproxy.yaml"),
		Inventory:  filepath.Join(e.ClusterPath, "config/nodes.yaml"),
		Become:     true,
		User:       e.SshUser(),
		PrivateKey: e.SshPKey(),
		Timeout:    3000,
	}

	return e.Ansible.Exec(pb)
}

// KubesprayCreate function calls an Ansible playbook that configures Kubernetes
// cluster.
func (e *kubespray) KubesprayCreate() error {
	vars := []string{
		"kube_version=" + e.K8sVersion(),
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
	vars := []string{
		"kube_version=" + e.K8sVersion(),
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
	vars := []string{
		"kube_version=" + e.K8sVersion(),
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
	vars := []string{
		"skip_confirmation=yes",
		"delete_nodes_confirmation=yes",
		"node=" + strings.Join(removedNodeNames, ","),
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
