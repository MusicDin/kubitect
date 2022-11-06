package actions

import (
	"cli/env"
	"cli/tools/terraform"
	"cli/utils"
	"fmt"
	"os"
	"path/filepath"
)

func Destroy(clusterName string) error {
	if clusterName == "" {
		return fmt.Errorf("A valid (non-empty) cluster name must be provided.")
	}

	clusters, err := ReadClustersInfo()

	if err != nil {
		return err
	}

	c := clusters.Find(clusterName)

	if c == nil {
		return fmt.Errorf("Cluster '%s' not found.", clusterName)
	}

	if !c.Active() {
		return fmt.Errorf("Cluster '%s' is already destroyed (or not yet initialized).", clusterName)
	}

	msg := fmt.Sprintf("The '%s' cluster will be destroyed.", clusterName)

	if err := utils.AskUserConfirmation(msg); err != nil {
		return err
	}

	fmt.Printf("Destroying '%s' cluster...\n", clusterName)

	if err := terraform.Destroy(c.Path); err != nil {
		return err
	}

	tfState := filepath.Join(c.Path, env.ConstTerraformStatePath)
	// TODO: Remove Kubeconfig

	return os.Remove(tfState)
}