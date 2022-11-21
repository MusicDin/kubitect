package actions

import (
	"cli/env"
	"cli/tools/terraform"
	"cli/ui"
	"fmt"
	"os"
)

func Destroy(ctx *env.Context, clusterName string) error {
	if clusterName == "" {
		return fmt.Errorf("a valid (non-empty) cluster name must be provided")
	}

	clusters, err := Clusters(ctx)

	if err != nil {
		return err
	}

	c := clusters.FindByName(clusterName)

	if c == nil {
		return fmt.Errorf("cluster '%s' not found.", clusterName)
	}

	count := clusters.CountByName(clusterName)

	if count > 1 {
		return fmt.Errorf("cannot destroy the cluster: multiple clusters (%d) have been found with the name '%s'", count, clusterName)
	}

	if !c.ContainsAppliedConfig() {
		return fmt.Errorf("cluster '%s' is already destroyed (or not yet initialized).", clusterName)
	}

	fmt.Printf("Cluster '%s' will be destroyed.\n", clusterName)

	if err := ui.Ask(); err != nil {
		return err
	}

	fmt.Printf("Destroying cluster '%s'...\n", clusterName)

	t, err := terraform.NewTerraform(ctx, c.Path)

	if err != nil {
		return err
	}

	if err := t.Destroy(); err != nil {
		return err
	}

	if err := os.Remove(c.TfStatePath()); err != nil {
		return err
	}

	return os.Remove(c.KubeconfigPath())
}
