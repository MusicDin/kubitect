package playbook

import (
	"cli/env"
	"cli/helpers"
	"path/filepath"
	"strings"
)

// KubitectInit function calls an Ansible playbook that creates a cluster config
// snapshot and generates Terraform main script.
func KubitectInit() error {

	extravars := []string{
		"kubitect_home=" + env.ProjectHomePath,
		"kubitect_cluster_action=" + env.ClusterAction,
		"kubitect_cluster_name=" + env.ClusterName,
		"kubitect_cluster_path=" + env.ClusterPath,
	}

	// Set custom config path if provided.
	if env.IsCustomConfig {
		extravars = append(extravars, "config_path="+env.ConfigPath)
	}

	err := helpers.ExecAnsiblePlaybookLocal(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
		Venv:         helpers.Venvs.Main,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/kubitect/init.yaml"),
		Extravars:    extravars,
	})

	if err != nil {
		return err
	}

	return nil
}

// KubitectHostsSetup function calls an Ansible playbook that ensures Kubitect target
// hosts meet all the requirements before cluster is created.
func KubitectHostsSetup() error {

	extravars := []string{
		"kubitect_cluster_path=" + env.ClusterPath,
	}

	err := helpers.ExecAnsiblePlaybook(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
		Venv:            helpers.Venvs.Main,
		PlaybookFile:    filepath.Join(env.ClusterPath, "ansible/kubitect/hosts-setup.yaml"),
		Inventory:       filepath.Join(env.ClusterPath, "config/hosts.yaml"),
		ConnectionLocal: true,
		Extravars:       extravars,
	})

	if err != nil {
		return err
	}

	return nil
}

// KubitectKubespraySetup functions calls an Ansible playbook that prepares Kubespray
// configuration files (all.yaml, k8s_cluster.yaml, ...) and clones Kubespray
// git project.
func KubitectKubespraySetup() error {

	extravars := []string{
		"kubitect_cluster_path=" + env.ClusterPath,
	}

	err := helpers.ExecAnsiblePlaybookLocal(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
		Venv:         helpers.Venvs.Main,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/kubitect/kubespray-setup.yaml"),
		Extravars:    extravars,
	})

	if err != nil {
		return err
	}

	return nil
}

// KubitectFinalize function calls an Ansible playbook that finalizes Kubernetes
// cluster installation.
func KubitectFinalize(sshUser string, sshPKey string) error {

	extravars := []string{
		"kubitect_cluster_path=" + env.ClusterPath,
	}

	err := helpers.ExecAnsiblePlaybook(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
		Venv:         helpers.Venvs.Main,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/kubitect/finalize.yaml"),
		Inventory:    filepath.Join(env.ClusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
		Extravars:    extravars,
	})

	if err != nil {
		return err
	}

	return nil
}

// HAProxySetup function calls an Ansible playbook that configures HAProxy
// load balancers.
func HAProxySetup(sshUser string, sshPKey string) error {

	extravars := []string{
		"kubitect_cluster_path=" + env.ClusterPath,
	}

	err := helpers.ExecAnsiblePlaybook(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
		Venv:         helpers.Venvs.Main,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/kubitect/haproxy.yaml"),
		Inventory:    filepath.Join(env.ClusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
		Extravars:    extravars,
	})

	if err != nil {
		return err
	}

	return nil
}

// KubesprayCreate function calls an Ansible playbook that configures Kubernetes
// cluster.
func KubesprayCreate(sshUser string, sshPKey string, k8sVersion string) error {

	extravars := []string{
		"kube_version=" + k8sVersion,
	}

	err := helpers.ExecAnsiblePlaybook(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
		Venv:         helpers.Venvs.Kubespray,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/kubespray/cluster.yml"),
		Inventory:    filepath.Join(env.ClusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
		Extravars:    extravars,
	})

	if err != nil {
		return err
	}

	return nil
}

// KubesprayUpgrade function calls an Ansible playbook that upgrades Kubernetes
// nodes to a newer version.
func KubesprayUpgrade(sshUser string, sshPKey string, k8sVersion string) error {

	extravars := []string{
		"kube_version=" + k8sVersion,
	}

	err := helpers.ExecAnsiblePlaybook(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
		Venv:         helpers.Venvs.Kubespray,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/kubespray/upgrade-cluster.yml"),
		Inventory:    filepath.Join(env.ClusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
		Extravars:    extravars,
	})
	if err != nil {
		return err
	}

	return nil
}

// KubesprayScale function calls an Ansible playbook that configures virtual machines
// that are freshly added to the cluster.
func KubesprayScale(sshUser string, sshPKey string, k8sVersion string) error {

	extravars := []string{
		"kube_version=" + k8sVersion,
	}

	err := helpers.ExecAnsiblePlaybook(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
		Venv:         helpers.Venvs.Kubespray,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/kubespray/scale.yml"),
		Inventory:    filepath.Join(env.ClusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
		Extravars:    extravars,
	})
	if err != nil {
		return err
	}

	return nil
}

// KubesprayRemoveNodes function calls an Ansible playbook that removes the nodes with
// the provided names.
func KubesprayRemoveNodes(sshUser string, sshPKey string, removedNodeNames []string) error {

	extravars := []string{
		"skip_confirmation=yes",
		"delete_nodes_confirmation=yes",
		"node=" + strings.Join(removedNodeNames, ","),
	}

	err := helpers.ExecAnsiblePlaybook(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
		Venv:         helpers.Venvs.Kubespray,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/kubespray/remove-node.yml"),
		Inventory:    filepath.Join(env.ClusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
		Extravars:    extravars,
	})

	if err != nil {
		return err
	}

	return nil
}
