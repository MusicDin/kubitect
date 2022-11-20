package playbook

import (
	"cli/tools/virtualenv"
	"strings"
)

type KubitectInitTag string

const (
	TAG_INIT      KubitectInitTag = "init"
	TAG_KUBESPRAY KubitectInitTag = "kubespray"
	TAG_GEN_NODES KubitectInitTag = "gen_nodes"
)

// KubitectInit function calls Ansible playbook that is responsible for initializing
// the cluster. Different workflow is executed based on the provided tags.
func KubitectInit(tags ...KubitectInitTag) error {
	var sTags []string

	for _, s := range tags {
		sTags = append(sTags, string(s))
	}

	pb := Playbook{
		PlaybookFile: "ansible/kubitect/init.yaml",
		Tags:         sTags,
		Local:        true,
	}

	return pb.Exec(virtualenv.MAIN)
}

// KubitectHostsSetup function calls an Ansible playbook that ensures Kubitect target
// hosts meet all the requirements before cluster is created.
func KubitectHostsSetup() error {
	pb := Playbook{
		PlaybookFile: "ansible/kubitect/hosts-setup.yaml",
		Inventory:    "config/hosts.yaml",
		Local:        true,
	}

	return pb.Exec(virtualenv.MAIN)
}

// KubitectFinalize function calls an Ansible playbook that finalizes Kubernetes
// cluster installation.
func KubitectFinalize(sshUser, sshPKey string) error {
	pb := Playbook{
		PlaybookFile: "ansible/kubitect/finalize.yaml",
		Inventory:    "config/nodes.yaml",
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
	}

	return pb.Exec(virtualenv.MAIN)
}

// HAProxy function calls an Ansible playbook that configures HAProxy
// load balancers.
func HAProxy(sshUser, sshPKey string) error {
	pb := Playbook{
		PlaybookFile: "ansible/kubitect/haproxy.yaml",
		Inventory:    "config/nodes.yaml",
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
	}

	return pb.Exec(virtualenv.MAIN)
}

// KubesprayCreate function calls an Ansible playbook that configures Kubernetes
// cluster.
func KubesprayCreate(sshUser, sshPKey, k8sVersion string) error {
	pb := Playbook{
		PlaybookFile: "ansible/kubespray/cluster.yml",
		Inventory:    "config/nodes.yaml",
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
		ExtraVars: []string{
			"kube_version=" + k8sVersion,
		},
	}

	return pb.Exec(virtualenv.KUBESPRAY)
}

// KubesprayUpgrade function calls an Ansible playbook that upgrades Kubernetes
// nodes to a newer version.
func KubesprayUpgrade(sshUser, sshPKey, k8sVersion string) error {
	pb := Playbook{
		PlaybookFile: "ansible/kubespray/upgrade-cluster.yml",
		Inventory:    "config/nodes.yaml",
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
		ExtraVars: []string{
			"kube_version=" + k8sVersion,
		},
	}

	return pb.Exec(virtualenv.KUBESPRAY)
}

// KubesprayScale function calls an Ansible playbook that configures virtual machines
// that are freshly added to the cluster.
func KubesprayScale(sshUser, sshPKey, k8sVersion string) error {
	pb := Playbook{
		PlaybookFile: "ansible/kubespray/scale.yml",
		Inventory:    "config/nodes.yaml",
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
		ExtraVars: []string{
			"kube_version=" + k8sVersion,
		},
	}

	return pb.Exec(virtualenv.KUBESPRAY)
}

// KubesprayRemoveNodes function calls an Ansible playbook that removes the nodes with
// the provided names.
func KubesprayRemoveNodes(sshUser, sshPKey string, removedNodeNames []string) error {
	pb := Playbook{
		PlaybookFile: "ansible/kubespray/remove-node.yml",
		Inventory:    "config/nodes.yaml",
		Become:       true,
		User:         sshUser,
		PrivateKey:   sshPKey,
		Timeout:      3000,
		ExtraVars: []string{
			"skip_confirmation=yes",
			"delete_nodes_confirmation=yes",
			"node=" + strings.Join(removedNodeNames, ","),
		},
	}

	return pb.Exec(virtualenv.KUBESPRAY)
}
