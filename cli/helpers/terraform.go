package helpers

import (
	"cli/env"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

const (
	terraformBinDir     = "bin/terraform"
	terraformProjectDir = "terraform"
)

// TerraformApply prepares Terraform project and applies the configuration.
func TerraformApply(clusterPath string) error {

	tf, err := terraformInit(clusterPath)
	if err != nil {
		return err
	}

	err = tf.Apply(context.Background())
	if err != nil {
		return fmt.Errorf("Error running Terraform apply: %w", err)
	}

	return nil
}

// TerraformDestroy destroys the Terraform project on the provided path.
func TerraformDestroy(clusterPath string) error {

	tf, err := terraformInit(clusterPath)
	if err != nil {
		return err
	}

	err = tf.Destroy(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to destroy Terraform project: %w", err)
	}

	return nil
}

// terraformInit installs terraform with the appropriate version
// into the bin directory. Afterwards it initializes the project
// and returns Terraform object.
func terraformInit(clusterPath string) (*tfexec.Terraform, error) {

	fmt.Println("Installing Terraform...")

	tfInstallDir := filepath.Join(env.ProjectHomePath, terraformBinDir, env.ConstTerraformVersion)
	tfProjectDir := filepath.Join(clusterPath, terraformProjectDir)

	// Make sure terraform install directory exists
	err := os.MkdirAll(tfInstallDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Failed creating Terraform install directory: %w", err)
	}

	installer := &releases.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion(env.ConstTerraformVersion)),
		InstallDir: tfInstallDir,
	}

	// Install specific version of Terraform into Terraform install directory.
	execPath, err := installer.Install(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Error installing Terraform: %w", err)
	}

	tf, err := tfexec.NewTerraform(tfProjectDir, execPath)
	if err != nil {
		return nil, fmt.Errorf("Error running NewTerraform: %w", err)
	}

	tf.SetStdout(os.Stdout)
	// tf.SetColor(true)

	fmt.Println("Initializing Terraform project...")

	err = tf.Init(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize Terraform project: %w", err)
	}

	return tf, nil
}
