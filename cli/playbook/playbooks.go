package playbook

import (
	"cli/env"
	"cli/helpers"
	"path/filepath"
	"strings"
)

// TkkInit function calls an Ansible playbook that creates a cluster config
// snapshot and generates Terraform main script.
func TkkInit() error {

	extravars := []string{
		"tkk_home=" + env.ProjectHomePath,
		"tkk_cluster_action=" + env.ClusterAction,
		"tkk_cluster_name=" + env.ClusterName,
		"tkk_cluster_path=" + env.ClusterPath,
	}

	// Set custom config path if provided.
	if env.IsCustomConfig {
		extravars = append(extravars, "config_path="+env.ConfigPath)
	}

	err := helpers.ExecAnsiblePlaybookLocal(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
		Venv:         helpers.MainVenv,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/tkk/init.yaml"),
		Extravars:    extravars,
	})

	if err != nil {
		return err
	}

	return nil
}

// TkkKubespraySetup functions calls an Ansible playbook that prepares Kubespray
// configuration files (all.yaml, k8s_cluster.yaml, ...) and clones Kubespray
// git project.
func TkkKubespraySetup() error {

	extravars := []string{
		"tkk_cluster_path=" + env.ClusterPath,
	}

	err := helpers.ExecAnsiblePlaybookLocal(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
		Venv:         helpers.MainVenv,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/tkk/kubespray-setup.yaml"),
		Extravars:    extravars,
	})

	if err != nil {
		return err
	}

	return nil
}

// TkkFinalize function calls an Ansible playbook that finalizes Kubernetes
// cluster installation.
func TkkFinalize(sshUser string, sshPKey string) error {

	extravars := []string{
		"tkk_cluster_path=" + env.ClusterPath,
	}

	err := helpers.ExecAnsiblePlaybook(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
		Venv:         helpers.MainVenv,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/tkk/finalize.yaml"),
		Inventory:    filepath.Join(env.ClusterPath, "config/hosts.ini"),
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

// HAProxyCreate function calls an Ansible playbook that configures HAProxy
// load balancers.
func HAProxyCreate(sshUser string, sshPKey string) error {

	err := helpers.ExecAnsiblePlaybook(env.ClusterPath, &helpers.AnsiblePlaybookCmd{
		Venv:         helpers.MainVenv,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/haproxy/haproxy.yaml"),
		Inventory:    filepath.Join(env.ClusterPath, "config/hosts.ini"),
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
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
		Venv:         helpers.KubesprayVenv,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/kubespray/cluster.yml"),
		Inventory:    filepath.Join(env.ClusterPath, "config/hosts.ini"),
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
		Venv:         helpers.KubesprayVenv,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/kubespray/upgrade-cluster.yml"),
		Inventory:    filepath.Join(env.ClusterPath, "config/hosts.ini"),
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
		Venv:         helpers.KubesprayVenv,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/kubespray/scale.yml"),
		Inventory:    filepath.Join(env.ClusterPath, "config/hosts.ini"),
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
		Venv:         helpers.KubesprayVenv,
		PlaybookFile: filepath.Join(env.ClusterPath, "ansible/kubespray/remove-node.yml"),
		Inventory:    filepath.Join(env.ClusterPath, "config/hosts.ini"),
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
