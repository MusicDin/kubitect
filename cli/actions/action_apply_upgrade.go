package actions

import "cli/tools/playbook"

func upgrade(c *Cluster) error {
	sshUser := string(*c.InfraConfig.Cluster.NodeTemplate.User)
	sshPKey := string(*c.InfraConfig.Cluster.NodeTemplate.SSH.PrivateKeyPath)

	k8sVersion := string(*c.NewConfig.Kubernetes.Version)

	if err := playbook.KubitectInit(playbook.TAG_INIT, playbook.TAG_KUBESPRAY, playbook.TAG_GEN_NODES); err != nil {
		return err
	}

	if err := playbook.KubitectHostsSetup(); err != nil {
		return err
	}

	if err := playbook.KubesprayUpgrade(sshUser, sshPKey, k8sVersion); err != nil {
		return err
	}

	return playbook.KubitectFinalize(sshUser, sshPKey)
}
