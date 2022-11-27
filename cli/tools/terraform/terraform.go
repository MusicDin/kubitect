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
	Ctx *env.Context

	version    string // Terraform version
	binPath    string // Terraform binary path
	projectDir string // Terraform project dir
}

// NewTerraform returns Terraform object with initialized Terraform project.
// Terraform is also installed, if binary file is not found locally.
func NewTerraform(ctx *env.Context, clusterPath string) (*Terraform, error) {
	var err error

	t := &Terraform{
		Ctx:        ctx,
		version:    env.ConstTerraformVersion,
		projectDir: filepath.Join(clusterPath, "terraform"),
	}

	binDir := filepath.Join(t.Ctx.ShareDir(), "terraform", t.version)

	fmt.Printf("Ensuring Terraform %s is installed...\n", t.version)

	fs := &fs.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion(t.version)),
		ExtraPaths: []string{binDir},
	}

	// Search for local Terraform installation before installing it.
	t.binPath, err = fs.Find(context.Background())

	if err == nil {
		fmt.Printf("Terraform %s found locally (%s).\n", t.version, t.binPath)
		return t, t.init()
	}

	fmt.Printf("Terraform %s could not be found locally.\n", t.version)
	fmt.Printf("Installing Terraform %s in '%s'...\n", t.version, binDir)

	if err := os.MkdirAll(binDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create Terraform install directory: %v", err)
	}

	installer := &releases.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion(t.version)),
		InstallDir: binDir,
	}

	// Install specific version of Terraform into shared directory.
	t.binPath, err = installer.Install(context.Background())

	if err != nil {
		return nil, fmt.Errorf("failed to install Terraform: %v", err)
	}

	return t, t.init()
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

	cmd.HasOutput(env.Debug || t.Ctx.ShowTerraformPlan())

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
	if changes && t.Ctx.ShowTerraformPlan() {
		if err := ui.GlobalUi().Ask("Proceed with terraform apply?"); err != nil {
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
