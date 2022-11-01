package actions

import "cli/tools/ansible"

func upgrade(c Cluster) error {
	sshUser := string(*c.InfraCfg.Cluster.NodeTemplate.User)
	sshPKey := string(*c.InfraCfg.Cluster.NodeTemplate.SSH.PrivateKeyPath)

	k8sVersion := string(*c.NewCfg.Kubernetes.Version)

	if err := ansible.KubitectInit(c.Path, ansible.KUBESPRAY, ansible.GEN_NODES); err != nil {
		return err
	}

	if err := c.SetupKubesprayVE(); err != nil {
		return err
	}

	if err := ansible.KubesprayUpgrade(c.Path, sshUser, sshPKey, k8sVersion); err != nil {
		return err
	}

	return ansible.KubitectFinalize(c.Path, sshUser, sshPKey)
}
