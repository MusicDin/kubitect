package actions

import (
	"cli/config/modelconfig"
	"cli/tools/playbook"
	"cli/tools/terraform"
	"cli/ui"
	"fmt"
)

func scale(c *Cluster, events Events) error {
	if err := scaleDown(c, events); err != nil {
		return err
	}

	t, err := terraform.NewTerraform(c.Ctx, c.Path)

	if err != nil {
		return err
	}

	if err := t.Apply(); err != nil {
		return err
	}

	return scaleUp(c, events)
}

// scaleUp adds new nodes to the cluster.
func scaleUp(c *Cluster, events Events) error {
	if len(events.OfType(SCALE_UP)) == 0 {
		return nil
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

	return playbook.KubesprayScale(sshUser, sshPKey, k8sVersion)
}

// scaleDown gracefully removes nodes from the cluster.
func scaleDown(c *Cluster, events Events) error {
	if len(events) == 0 {
		return nil
	}

	sshUser := string(*c.InfraConfig.Cluster.NodeTemplate.User)
	sshPKey := string(*c.InfraConfig.Cluster.NodeTemplate.SSH.PrivateKeyPath)

	rmNodes, err := extractNodes(events.OfType(SCALE_DOWN))

	if err != nil {
		return err
	}

	if len(rmNodes) == 0 {
		return nil
	}

	fmt.Println("The following nodes will get removed:")

	var names []string

	for _, n := range rmNodes {
		name := fmt.Sprintf("%s-%s-%s", c.Name, n.GetTypeName(), *n.GetID())
		names = append(names, name)

		fmt.Println("-", name)
	}

	if err := ui.GlobalUi().Ask(); err != nil {
		return err
	}

	if err := playbook.KubitectInit(playbook.TAG_KUBESPRAY); err != nil {
		return err
	}

	return playbook.KubesprayRemoveNodes(sshUser, sshPKey, names)
}

// extractNodes returns node instances from the event changes.
func extractNodes(events Events) ([]modelconfig.Instance, error) {
	var nodes []modelconfig.Instance

	for _, e := range events {
		for _, ch := range e.changes {
			if i, ok := ch.Before.(modelconfig.Instance); ok {
				nodes = append(nodes, i)
			}

			return nil, fmt.Errorf("%v cannot be scaled", ch.Type.Name())
		}
	}

	return nodes, nil
}
