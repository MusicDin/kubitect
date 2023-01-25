package cluster

import (
	"cli/ui"
	"fmt"
	"os"
)

func (c *ClusterMeta) Destroy() error {
	if !c.ContainsTfStateConfig() {
		return fmt.Errorf("cluster '%s' is already destroyed (or not yet initialized).", c.Name)
	}

	c.Ui().Printf(ui.INFO, "Cluster '%s' will be destroyed.\n", c.Name)

	if err := c.Ui().Ask(); err != nil {
		return err
	}

	c.Ui().Printf(ui.INFO, "Destroying cluster '%s'...\n", c.Name)

	if err := c.Provisioner().Destroy(); err != nil {
		return err
	}

	os.Remove(c.TfStatePath())
	os.Remove(c.KubeconfigPath())

	return nil
}
