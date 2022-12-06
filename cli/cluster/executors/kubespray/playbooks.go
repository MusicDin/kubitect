package kubespray

import (
	"cli/tools/ansible"
	"cli/tools/virtualenv"
	"fmt"
	"path/filepath"
	"strings"
)

// play executes a playbook using ansible from a given virtual environment.
func (e *KubesprayExecutor) play(ve *virtualenv.VirtualEnv, pb ansible.Playbook) error {
	if err := ve.Init(); err != nil {
		return fmt.Errorf("kubespray exec: failed initializing virtual environment: %v", err)
	}

	ansible := ansible.Ansible{
		BinPath:    filepath.Join(ve.Path, "bin", "ansible-playbook"),
		WorkingDir: e.ClusterPath,
		Ui:         e.Ui,
	}

	return ansible.Exec(pb)
}

type PlaybookTag string

const (
	TAG_INIT      PlaybookTag = "init"
	TAG_KUBESPRAY PlaybookTag = "kubespray"
	TAG_GEN_NODES PlaybookTag = "gen_nodes"
)

// KubitectInit function calls Ansible playbook that is responsible for initializing
// the cluster. Different workflow is executed based on the provided tags.
func (e *KubesprayExecutor) KubitectInit(tags ...PlaybookTag) error {
	var sTags []string

	for _, s := range tags {
		sTags = append(sTags, string(s))
	}

	pb := ansible.Playbook{
		PlaybookFile: filepath.Join(e.ClusterPath, "ansible/kubitect/init.yaml"),
		Tags:         sTags,
		Local:        true,
	}

	return e.play(e.Venvs.MAIN, pb)
}

// KubitectHostsSetup function calls an Ansible playbook that ensures Kubitect target
// hosts meet all the requirements before cluster is created.
func (e *KubesprayExecutor) KubitectHostsSetup() error {
	pb := ansible.Playbook{
		PlaybookFile: filepath.Join(e.ClusterPath, "ansible/kubitect/hosts-setup.yaml"),
		Inventory:    filepath.Join(e.ClusterPath, "config/hosts.yaml"),
		Local:        true,
	}

	return e.play(e.Venvs.MAIN, pb)
}

// KubitectFinalize function calls an Ansible playbook that finalizes Kubernetes
// cluster installation.
func (e *KubesprayExecutor) KubitectFinalize() error {
	pb := ansible.Playbook{
		PlaybookFile: filepath.Join(e.ClusterPath, "ansible/kubitect/finalize.yaml"),
		Inventory:    filepath.Join(e.ClusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         e.SshUser,
		PrivateKey:   e.SshPKey,
		Timeout:      3000,
	}

	return e.play(e.Venvs.MAIN, pb)
}

// HAProxy function calls an Ansible playbook that configures HAProxy
// load balancers.
func (e *KubesprayExecutor) HAProxy() error {
	pb := ansible.Playbook{
		PlaybookFile: filepath.Join(e.ClusterPath, "ansible/kubitect/haproxy.yaml"),
		Inventory:    filepath.Join(e.ClusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         e.SshUser,
		PrivateKey:   e.SshPKey,
		Timeout:      3000,
	}

	return e.play(e.Venvs.MAIN, pb)
}

// KubesprayCreate function calls an Ansible playbook that configures Kubernetes
// cluster.
func (e *KubesprayExecutor) KubesprayCreate() error {
	vars := []string{
		"kube_version=" + e.K8sVersion,
	}

	pb := ansible.Playbook{
		PlaybookFile: filepath.Join(e.ClusterPath, "ansible/kubespray/cluster.yml"),
		Inventory:    filepath.Join(e.ClusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         e.SshUser,
		PrivateKey:   e.SshPKey,
		Timeout:      3000,
		ExtraVars:    vars,
	}

	return e.play(e.Venvs.KUBESPRAY, pb)
}

// KubesprayUpgrade function calls an Ansible playbook that upgrades Kubernetes
// nodes to a newer version.
func (e *KubesprayExecutor) KubesprayUpgrade() error {
	vars := []string{
		"kube_version=" + e.K8sVersion,
	}

	pb := ansible.Playbook{
		PlaybookFile: filepath.Join(e.ClusterPath, "ansible/kubespray/upgrade-cluster.yml"),
		Inventory:    filepath.Join(e.ClusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         e.SshUser,
		PrivateKey:   e.SshPKey,
		Timeout:      3000,
		ExtraVars:    vars,
	}

	return e.play(e.Venvs.KUBESPRAY, pb)
}

// KubesprayScale function calls an Ansible playbook that configures virtual machines
// that are freshly added to the cluster.
func (e *KubesprayExecutor) KubesprayScale() error {
	vars := []string{
		"kube_version=" + e.K8sVersion,
	}

	pb := ansible.Playbook{
		PlaybookFile: filepath.Join(e.ClusterPath, "ansible/kubespray/scale.yml"),
		Inventory:    filepath.Join(e.ClusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         e.SshUser,
		PrivateKey:   e.SshPKey,
		Timeout:      3000,
		ExtraVars:    vars,
	}

	return e.play(e.Venvs.KUBESPRAY, pb)
}

// KubesprayRemoveNodes function calls an Ansible playbook that removes the nodes with
// the provided names.
func (e *KubesprayExecutor) KubesprayRemoveNodes(removedNodeNames []string) error {
	vars := []string{
		"skip_confirmation=yes",
		"delete_nodes_confirmation=yes",
		"node=" + strings.Join(removedNodeNames, ","),
	}

	pb := ansible.Playbook{
		PlaybookFile: filepath.Join(e.ClusterPath, "ansible/kubespray/remove-node.yml"),
		Inventory:    filepath.Join(e.ClusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         e.SshUser,
		PrivateKey:   e.SshPKey,
		Timeout:      3000,
		ExtraVars:    vars,
	}

	return e.play(e.Venvs.KUBESPRAY, pb)
}
