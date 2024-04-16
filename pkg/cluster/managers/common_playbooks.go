package managers

import (
	"path/filepath"

	"github.com/MusicDin/kubitect/pkg/tools/ansible"
)

// haproxy calls playbook that configures external HAProxy load balancers.
func (e common) HAProxy() error {
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

// finalize calls playbook that finalizes Kubernetes cluster installation.
// This includes exp
func (e common) Finalize() error {
	vars := map[string]string{
		"bin_dir": e.SharedDir,
	}

	pb := ansible.Playbook{
		Path:       filepath.Join(e.ClusterPath, "ansible/kubitect/finalize.yaml"),
		Inventory:  filepath.Join(e.ClusterPath, "config/nodes.yaml"),
		Become:     true,
		User:       e.SshUser(),
		PrivateKey: e.SshPKey(),
		Timeout:    3000,
		ExtraVars:  vars,
	}

	return e.Ansible.Exec(pb)
}
