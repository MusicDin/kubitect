package cluster

import (
	"fmt"
	"os"

	"github.com/MusicDin/kubitect/pkg/ui"
	"github.com/MusicDin/kubitect/pkg/utils/file"
)

// Destroy destroys the cluster and removes cluster's directory.
// Terraform resources are wiped if terraform state file is found.
func (c *ClusterMeta) Destroy() error {
	if !file.Exists(c.Path) {
		return fmt.Errorf("cluster %q does not exist", c.Name)
	}

	ui.Printf(ui.INFO, "Cluster %q will be destroyed.\n", c.Name)
	if err := ui.Ask(); err != nil {
		return err
	}

	if c.ContainsTfStateConfig() {
		ui.Println(ui.INFO, "Removing cluster resources...")
		if err := c.Provisioner().Destroy(); err != nil {
			return err
		}
	}

	ui.Println(ui.INFO, "Removing cluster cache...")
	_ = os.RemoveAll(c.CacheDir())

	ui.Println(ui.INFO, "Removing cluster directory...")
	if err := os.RemoveAll(c.Path); err != nil {
		return fmt.Errorf("failed to remove directory of the cluster %q: %v", c.Name, err)
	}

	ui.Printf(ui.INFO, "Cluster %q has been successfully destroyed.\n", c.Name)
	return nil
}
