package cluster

import (
	"fmt"
	"os"

	"github.com/MusicDin/kubitect/pkg/ui"
)

// Destroy destroys an active cluster and removes cluster's
// directory. If cluster does not exist or does not contain
// a terraform state file (is inactive), an error is returned.
func (c *ClusterMeta) Destroy() error {
	if !c.ContainsTfStateConfig() {
		return fmt.Errorf("cluster '%s' is already destroyed (or not yet initialized).", c.Name)
	}

	ui.Printf(ui.INFO, "Cluster '%s' will be destroyed.\n", c.Name)
	if err := ui.Ask(); err != nil {
		return err
	}

	ui.Println(ui.INFO, "Destroying cluster...")
	if err := c.Provisioner().Destroy(); err != nil {
		return err
	}

	ui.Println(ui.INFO, "Cleaning up cluster directory...", c.Name)
	if err := os.RemoveAll(c.Path); err != nil {
		return fmt.Errorf("failed to remove directory of the cluster '%s': %v", c.Name, err)
	}

	ui.Printf(ui.INFO, "Cluster '%s' has been successfully destroyed.\n", c.Name)
	return nil
}
