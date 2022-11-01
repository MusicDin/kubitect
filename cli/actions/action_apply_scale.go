package actions

import (
	"cli/tools/ansible"
	"cli/tools/terraform"
)

func scale(c Cluster, events []*OnChangeEvent) error {

	// TODO: Remove nodes
	// 1.) Extract remove nodes from events
	// 2.) Ask user permission
	// 3.) Remove nodes
	// 3.1) Call Kubespray remove for that node
	// 3.2) Remove that node from the config
	// 4.) TF Apply to remove VMs
	// 5.) Run scale playbook

	if err := terraform.Apply(c.Path); err != nil {
		return err
	}

	if err := c.Sync(); err != nil {
		return err
	}

	sshUser := string(*c.InfraCfg.Cluster.NodeTemplate.User)
	sshPKey := string(*c.InfraCfg.Cluster.NodeTemplate.SSH.PrivateKeyPath)

	k8sVersion := string(*c.NewCfg.Kubernetes.Version)

	if err := ansible.HAProxy(c.Path, sshUser, sshPKey); err != nil {
		return err
	}

	if err := ansible.KubesprayScale(c.Path, sshUser, sshPKey, k8sVersion); err != nil {
		return err
	}

	return nil
}
