package ansible

import (
	"cli/tools/virtualenv"
	"path/filepath"
	"strings"
)

type KubitectInitTag string

const (
	INIT      KubitectInitTag = "init"
	KUBESPRAY KubitectInitTag = "kubespray"
	GEN_NODES KubitectInitTag = "gen_nodes"
)

// KubitectInit function calls Ansible playbook that is responsible for initializing
// the cluster. Different workflow is executed based on the provided tags.
func KubitectInit(clusterPath string, tags ...KubitectInitTag) error {
	var sTags []string

	for _, s := range tags {
		sTags = append(sTags, string(s))
	}

	pb := Playbook{
		VenvPath:     virtualenv.Main.Path,
		PlaybookFile: filepath.Join(clusterPath, "ansible/kubitect/init.yaml"),
		Tags:         sTags,
		Local:        true,
	}

	return pb.Exec()
}

// kubitectHostsSetup function calls an Ansible playbook that ensures Kubitect target
// hosts meet all the requirements before cluster is created.
func KubitectHostsSetup(clusterPath string) error {
	pb := Playbook{
		VenvPath:     virtualenv.Main.Path,
		PlaybookFile: filepath.Join(clusterPath, "ansible/kubitect/hosts-setup.yaml"),
		Inventory:    filepath.Join(clusterPath, "config/hosts.yaml"),
		Local:        true,
	}

	return pb.Exec()
}

// kubitectKubespraySetup functions calls an Ansible playbook that prepares Kubespray
// configuration files (all.yaml, k8s_cluster.yaml, ...) and clones Kubespray
// git project.
// func kubitectKubespraySetup(clusterPath string, tags ...string) error {
// 	ksSetupPB := AnsiblePlaybookCmd{
// 		VenvPath:     helpers.Venvs.Main.Path,
// 		PlaybookFile: filepath.Join(clusterPath, "ansible/kubitect/kubespray-setup.yaml"),
// 		Tags:         tags,
// 	}

// 	return ExecAnsiblePlaybookLocal(clusterPath, &ksSetupPB)
// }

// KubitectFinalize function calls an Ansible playbook that finalizes Kubernetes
// cluster installation.
func KubitectFinalize(clusterPath, sshUser, sshPKey string) error {
	pb := Playbook{
		VenvPath:     virtualenv.Main.Path,
		PlaybookFile: filepath.Join(clusterPath, "ansible/kubitect/finalize.yaml"),
		Inventory:    filepath.Join(clusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
	}

	return pb.Exec()
}

// HAProxySetup function calls an Ansible playbook that configures HAProxy
// load balancers.
func HAProxy(clusterPath, sshUser, sshPKey string) error {
	pb := Playbook{
		VenvPath:     virtualenv.Main.Path,
		PlaybookFile: filepath.Join(clusterPath, "ansible/kubitect/haproxy.yaml"),
		Inventory:    filepath.Join(clusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
	}

	return pb.Exec()
}

// KubesprayCreate function calls an Ansible playbook that configures Kubernetes
// cluster.
func KubesprayCreate(clusterPath, sshUser, sshPKey, k8sVersion string) error {
	extraVars := []string{
		"kube_version=" + k8sVersion,
	}

	pb := Playbook{
		VenvPath:     virtualenv.Kubespray.Path,
		PlaybookFile: filepath.Join(clusterPath, "ansible/kubespray/cluster.yml"),
		Inventory:    filepath.Join(clusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
		ExtraVars:    extraVars,
	}

	return pb.Exec()
}

// KubesprayUpgrade function calls an Ansible playbook that upgrades Kubernetes
// nodes to a newer version.
func KubesprayUpgrade(clusterPath, sshUser, sshPKey, k8sVersion string) error {
	extraVars := []string{
		"kube_version=" + k8sVersion,
	}

	pb := Playbook{
		VenvPath:     virtualenv.Kubespray.Path,
		PlaybookFile: filepath.Join(clusterPath, "ansible/kubespray/upgrade-cluster.yml"),
		Inventory:    filepath.Join(clusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
		ExtraVars:    extraVars,
	}

	return pb.Exec()
}

// KubesprayScale function calls an Ansible playbook that configures virtual machines
// that are freshly added to the cluster.
func KubesprayScale(clusterPath, sshUser, sshPKey, k8sVersion string) error {
	extraVars := []string{
		"kube_version=" + k8sVersion,
	}

	pb := Playbook{
		VenvPath:     virtualenv.Kubespray.Path,
		PlaybookFile: filepath.Join(clusterPath, "ansible/kubespray/scale.yml"),
		Inventory:    filepath.Join(clusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
		ExtraVars:    extraVars,
	}

	return pb.Exec()
}

// KubesprayRemoveNodes function calls an Ansible playbook that removes the nodes with
// the provided names.
func KubesprayRemoveNodes(clusterPath, sshUser, sshPKey string, removedNodeNames []string) error {
	extraVars := []string{
		"skip_confirmation=yes",
		"delete_nodes_confirmation=yes",
		"node=" + strings.Join(removedNodeNames, ","),
	}

	pb := Playbook{
		VenvPath:     virtualenv.Kubespray.Path,
		PlaybookFile: filepath.Join(clusterPath, "ansible/kubespray/remove-node.yml"),
		Inventory:    filepath.Join(clusterPath, "config/nodes.yaml"),
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
		ExtraVars:    extraVars,
	}

	return pb.Exec()
}
