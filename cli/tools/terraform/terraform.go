package terraform

import (
	"cli/env"
	"cli/ui"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
)

type Terraform struct {
	Version    string // Terraform version
	BinPath    string // Terraform binary path
	WorkingDir string // Terraform project dir
	ShowPlan   bool
}

// NewTerraform returns Terraform object with initialized Terraform project.
// Terraform is also installed, if binary file is not found locally.
func NewTerraform(ctx *env.Context, clusterPath string) (*Terraform, error) {
	ver := env.ConstTerraformVersion

	binDir := filepath.Join(ctx.ShareDir(), "terraform", ver)
	binPath, err := findOrInstall(ver, binDir)

	if err != nil {
		return nil, err
	}

	tf := &Terraform{
		Version:    ver,
		BinPath:    binPath,
		WorkingDir: filepath.Join(clusterPath, "terraform"),
	}

	return tf, tf.init()
}

// init initializes a Terraform project.
func (t *Terraform) init() error {
	cmd := t.NewCmd("init")

	cmd.AddArg("force-copy")
	cmd.AddArg("input", false)
	cmd.AddArg("get", true)

	_, err := cmd.Run()

	return err
}

// Plan shows Terraform project changes (plan).
func (t *Terraform) Plan() (bool, error) {
	cmd := t.NewCmd("plan")

	cmd.ShowOutput(env.Debug || t.ShowPlan)

	cmd.AddArg("detailed-exitcode")
	cmd.AddArg("input", false)
	cmd.AddArg("lock", true)
	cmd.AddArg("lock-timeout", "0s")
	cmd.AddArg("parallelism", 10)
	cmd.AddArg("refresh", true)

	exitCode, err := cmd.Run()

	if err != nil && exitCode == 2 {
		return true, nil
	}

	return false, err
}

// Apply applies new Terraform configurations.
func (t *Terraform) Apply() error {
	changes, err := t.Plan()

	if err != nil {
		return err
	}

	// Ask user for permission if there are any changes
	if changes && t.ShowPlan {
		err := ui.GlobalUi().Ask("Proceed with terraform apply?")

		if err != nil {
			return err
		}
	}

	cmd := t.NewCmd("apply")

	cmd.AddArg("auto-approve")
	cmd.AddArg("input", false)
	cmd.AddArg("lock", true)
	cmd.AddArg("lock-timeout", "0s")
	cmd.AddArg("parallelism", 10)
	cmd.AddArg("refresh", true)

	_, err = cmd.Run()

	return err
}

// Destroy destroys the Terraform project.
func (t *Terraform) Destroy() error {
	cmd := t.NewCmd("destroy")

	cmd.AddArg("auto-approve")
	cmd.AddArg("input", false)
	cmd.AddArg("lock", true)
	cmd.AddArg("lock-timeout", "0s")
	cmd.AddArg("parallelism", 10)
	cmd.AddArg("refresh", true)

	_, err := cmd.Run()

	return err
}

// findOrInstall first searches for Terraform binary locally and if
// binary is not found, it is installed in given binDir.
func findOrInstall(ver, binDir string) (string, error) {
	var binPath string
	var err error

	fmt.Printf("Ensuring Terraform %s is installed...\n", ver)

	fs := &fs.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion(ver)),
		ExtraPaths: []string{binDir},
	}

	// Search for local Terraform installation before installing it.
	binPath, err = fs.Find(context.Background())

	if err == nil {
		fmt.Printf("Terraform %s found locally (%s).\n", ver, binPath)
		return binPath, nil
	}

	fmt.Printf("Terraform %s could not be found locally.\n", ver)
	fmt.Printf("Installing Terraform %s in '%s'...\n", ver, binDir)

	if err := os.MkdirAll(binDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create Terraform install directory: %v", err)
	}

	installer := &releases.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion(ver)),
		InstallDir: binDir,
	}

	// Install specific version of Terraform into shared directory.
	binPath, err = installer.Install(context.Background())

	if err != nil {
		return "", fmt.Errorf("failed to install Terraform: %v", err)
	}

	return binPath, nil
}
