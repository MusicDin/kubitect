package terraform

import (
	"cli/cluster/provisioner"
	"cli/config/modelconfig"
	"cli/ui"
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
)

type terraform struct {
	binPath string

	Version    string
	BinDir     string
	WorkingDir string
	ShowPlan   bool

	clusterPath string
	hosts       []modelconfig.Host

	initialized bool

	Ui *ui.Ui
}

func NewTerraform(
	version,
	clusterPath,
	binDir,
	workingDir string,
	hosts []modelconfig.Host,
	showPlan bool,
	ui *ui.Ui,
) provisioner.Provisioner {
	return &terraform{
		Version:    version,
		BinDir:     binDir,
		WorkingDir: workingDir,
		ShowPlan:   showPlan,

		clusterPath: clusterPath,
		hosts:       hosts,

		Ui: ui,
	}
}

func (t *terraform) Init() error {
	if t.initialized {
		return nil
	}

	binPath, err := t.findOrInstall()
	if err != nil {
		return err
	}

	t.binPath = binPath

	cmd := t.NewCmd("init")

	cmd.AddArg("force-copy")
	cmd.AddArg("input", false)
	cmd.AddArg("get", true)

	_, err = cmd.Run()

	if err == nil {
		t.initialized = true
	}

	return err
}

// init initializes a Terraform project.
func (t *terraform) init() error {
	if t.binPath != "" {
		return nil
	}

	binPath, err := t.findOrInstall()
	if err != nil {
		return err
	}

	t.binPath = binPath

	cmd := t.NewCmd("init")

	cmd.AddArg("force-copy")
	cmd.AddArg("input", false)
	cmd.AddArg("get", true)

	_, err = cmd.Run()

	return err
}

// Plan shows Terraform project changes (plan).
// It returns a potential error and whether there
// are changes or not.
func (t *terraform) Plan() (bool, error) {
	if err := NewMainTemplate(t.hosts).Write(t.clusterPath); err != nil {
		return false, err
	}

	if err := t.init(); err != nil {
		return false, err
	}

	cmd := t.NewCmd("plan")

	cmd.ShowOutput(t.Ui.Debug || t.ShowPlan)

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
func (t *terraform) Apply() error {
	changes, err := t.Plan()

	if err != nil {
		return err
	}

	// Ask user for permission if there are any changes
	if changes && t.ShowPlan {
		err := t.Ui.Ask("Proceed with terraform apply?")

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
func (t *terraform) Destroy() error {
	err := t.init()
	if err != nil {
		return err
	}

	cmd := t.NewCmd("destroy")

	cmd.AddArg("auto-approve")
	cmd.AddArg("input", false)
	cmd.AddArg("lock", true)
	cmd.AddArg("lock-timeout", "0s")
	cmd.AddArg("parallelism", 10)
	cmd.AddArg("refresh", true)

	_, err = cmd.Run()
	return err
}

// findOrInstall first searches for Terraform binary locally and if
// binary is not found, it is installed in given binDir.
func (t *terraform) findOrInstall() (string, error) {
	var binPath string
	var err error

	t.Ui.Printf(ui.INFO, "Ensuring Terraform %s is installed...\n", t.Version)

	binPath, err = findTerraform(t.Version, t.BinDir)

	if err == nil {
		t.Ui.Printf(ui.INFO, "Terraform %s found locally (%s).\n", t.Version, binPath)
		return binPath, nil
	}

	t.Ui.Printf(ui.INFO, "Terraform %s could not be found locally.\n", t.Version)
	t.Ui.Printf(ui.INFO, "Installing Terraform %s in '%s'...\n", t.Version, t.BinDir)

	binPath, err = installTerraform(t.Version, t.BinDir)

	if err != nil {
		return "", fmt.Errorf("failed to install Terraform: %v", err)
	}

	return binPath, nil
}

// findTerraform searches for Terraform binary locally.
// If binary is found, its path is returned.
func findTerraform(ver, binDir string) (string, error) {
	fs := &fs.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion(ver)),
		ExtraPaths: []string{binDir},
	}

	return fs.Find(context.Background())
}

// installTerraform installs Terraform in a given directory.
func installTerraform(ver, binDir string) (string, error) {
	if err := os.MkdirAll(binDir, os.ModePerm); err != nil {
		return "", err
	}

	installer := &releases.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion(ver)),
		InstallDir: binDir,
	}

	return installer.Install(context.Background())
}
