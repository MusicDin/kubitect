package actions

import (
	"cli/tools/playbook"
	"cli/tools/terraform"
)

func create(c *Cluster) error {
	t, err := terraform.NewTerraform(c.Ctx, c.Path)

	if err != nil {
		return err
	}

	if err := t.Apply(); err != nil {
		return err
	}

	if err := c.Sync(); err != nil {
		return err
	}

	sshUser := string(*c.InfraConfig.Cluster.NodeTemplate.User)
	sshPKey := string(*c.InfraConfig.Cluster.NodeTemplate.SSH.PrivateKeyPath)

	k8sVersion := string(*c.NewConfig.Kubernetes.Version)

	if err := playbook.KubitectInit(playbook.TAG_KUBESPRAY, playbook.TAG_GEN_NODES); err != nil {
		return err
	}

	if err := playbook.HAProxy(sshUser, sshPKey); err != nil {
		return err
	}

	if err := playbook.KubesprayCreate(sshUser, sshPKey, k8sVersion); err != nil {
		return err
	}

	return playbook.KubitectFinalize(sshUser, sshPKey)
}
