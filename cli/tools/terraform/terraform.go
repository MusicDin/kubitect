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
	"github.com/hashicorp/terraform-exec/tfexec"
)

type Terraform struct {
	Ctx *env.Context

	tf *tfexec.Terraform

	version    string // Terraform version
	path       string // Terraform binary path
	projectDir string // Terraform project dir
}

// NewTerraform returns Terraform object with initialized Terraform project.
// Terraform is also installed, if binary file is not found locally.
func NewTerraform(ctx *env.Context, clusterPath string) (*Terraform, error) {
	t := Terraform{
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
	execPath, err := fs.Find(context.Background())

	if err != nil {
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

		// Install specific version of Terraform into Terraform install directory.
		execPath, err = installer.Install(context.Background())

		if err != nil {
			return nil, fmt.Errorf("terraform: failed to install Terraform: %v", err)
		}

	} else {
		fmt.Printf("Terraform %s found locally (%s).\n", t.version, execPath)
	}

	tf, err := tfexec.NewTerraform(t.projectDir, execPath)

	if err != nil {
		return nil, fmt.Errorf("failed to instantiate Terraform: %v", err)
	}

	tf.SetStdout(ui.GlobalUi().Streams.Out.File)
	tf.SetStderr(ui.GlobalUi().Streams.Err.File)
	// tf.SetColor(true)

	t.tf = tf
	t.path = execPath

	fmt.Println("Initializing Terraform project...")

	if err = tf.Init(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize Terraform project")
	}

	return &t, nil
}

// Plan shows Terraform project changes (plan).
func (t *Terraform) Plan() (bool, error) {
	if t.Ctx.ShowTerraformPlan() {
		t.tf.SetStdout(nil)
		t.tf.SetStderr(nil)
	}

	changes, err := t.tf.Plan(context.Background())

	t.tf.SetStdout(ui.GlobalUi().Streams.Out.File)
	t.tf.SetStderr(ui.GlobalUi().Streams.Err.File)

	if err != nil {
		err = fmt.Errorf("error running Terraform plan")
	}

	return changes, err
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

	if err := t.tf.Apply(context.Background()); err != nil {
		return fmt.Errorf("error running Terraform apply")
	}

	return nil
}

// Destroy destroys the Terraform project.
func (t *Terraform) Destroy() error {
	if err := t.tf.Destroy(context.Background()); err != nil {
		return fmt.Errorf("failed to destroy Terraform project")
	}

	return nil
}
