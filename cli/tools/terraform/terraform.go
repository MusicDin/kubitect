package terraform

import (
	"cli/env"
	"cli/utils"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

// Apply prepares Terraform project and applies the configuration.
func Apply(clusterPath string) error {
	tf, err := new(clusterPath)

	if err != nil {
		return err
	}

	// run 'terraform plan' first to show changes
	if !env.AutoApprove {
		changes, err := tf.Plan(context.Background())

		if err != nil {
			return fmt.Errorf("Error running Terraform apply: %v", err)
		}

		// Ask user for permission if there are any changes
		if changes {
			warning := "Proceed with terraform apply?"
			err := utils.AskUserConfirmation(warning)

			if err != nil {
				return err
			}
		}
	}

	err = tf.Apply(context.Background())

	if err != nil {
		return fmt.Errorf("Error running Terraform apply: %v", err)
	}

	return nil
}

// Destroy destroys the Terraform project on the provided path.
func Destroy(clusterPath string) error {
	tf, err := new(clusterPath)

	if err != nil {
		return err
	}

	err = tf.Destroy(context.Background())

	if err != nil {
		return fmt.Errorf("Failed to destroy Terraform project: %v", err)
	}

	return nil
}

// new installs terraform with the appropriate version
// into the bin directory. Afterwards it initializes the project
// and returns Terraform object.
func new(clusterPath string) (*tfexec.Terraform, error) {
	ver := env.ConstTerraformVersion

	binDir := env.BinDirPath("terraform", ver)
	projDir := filepath.Join(clusterPath, "terraform")

	fmt.Printf("Ensuring Terraform %s is installed...\n", ver)

	fs := &fs.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion(ver)),
		ExtraPaths: []string{binDir},
	}

	// Search for local Terraform installation before installing it.
	execPath, err := fs.Find(context.Background())

	if err != nil {
		fmt.Printf("Terraform %s could not be found locally.\n", ver)
		fmt.Printf("Installing Terraform %s in '%s'...\n", ver, binDir)

		err := os.MkdirAll(binDir, os.ModePerm)

		if err != nil {
			return nil, fmt.Errorf("Failed creating Terraform install directory: %v", err)
		}

		installer := &releases.ExactVersion{
			Product:    product.Terraform,
			Version:    version.Must(version.NewVersion(ver)),
			InstallDir: binDir,
		}

		// Install specific version of Terraform into Terraform install directory.
		execPath, err = installer.Install(context.Background())

		if err != nil {
			return nil, fmt.Errorf("Error installing Terraform: %v", err)
		}

	} else {
		fmt.Printf("Terraform %s found locally (%s).\n", ver, execPath)
	}

	tf, err := tfexec.NewTerraform(projDir, execPath)

	if err != nil {
		return nil, fmt.Errorf("Error running NewTerraform: %v", err)
	}

	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)
	// tf.SetColor(true)

	fmt.Println("Initializing Terraform project...")

	err = tf.Init(context.Background())

	if err != nil {
		return nil, fmt.Errorf("Failed to initialize Terraform project: %v", err)
	}

	return tf, nil
}
