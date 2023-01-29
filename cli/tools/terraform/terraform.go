package terraform

import (
	"cli/ui"
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
)

type (
	Terraform interface {
		Plan() (exitCode bool, err error)
		Apply() error
		Destroy() error
	}

	terraform struct {
		// Required Terraform version
		version string

		// Path where terraform binary will be installed
		// if it is not found locally.
		binDir string

		// Evaluated during findAndInstall.
		binPath string

		// Dir where main.tf is located (root Terraform dir).
		projectDir string

		// If true, Terraform plan will be shown.
		showPlan bool

		// Indicates that terraform project has
		// been already initialized.
		initialized bool
	}
)

func NewTerraform(
	version,
	// clusterPath,
	binDir,
	workingDir string,
	// hosts []modelconfig.Host,
	showPlan bool,
) Terraform {
	return &terraform{
		version:    version,
		binDir:     binDir,
		projectDir: workingDir,
		showPlan:   showPlan,
	}
}

// init initializes a Terraform project.
func (t *terraform) init() error {
	if t.initialized {
		return nil
	}

	binPath, err := t.findOrInstall()
	if err != nil {
		return err
	}

	t.binPath = binPath

	args := []string{
		flag("force-copy"),
		flag("input", false),
		flag("get", true),
	}

	_, err = t.runCmd("init", args, true)

	if err == nil {
		t.initialized = true
	}

	return err
}

// Plan shows Terraform project changes (plan).
// It returns a potential error and whether there
// are changes or not.
func (t *terraform) Plan() (bool, error) {
	if err := t.init(); err != nil {
		return false, err
	}

	args := []string{
		flag("detailed-exitcode"),
		flag("input", false),
		flag("lock", true),
		flag("lock-timeout", "0s"),
		flag("parallelism", 10),
		flag("refresh", true),
	}

	exitCode, err := t.runCmd("plan", args, t.showPlan)

	// "exitCode 2" indicates terraform plan changes
	if err != nil && exitCode == 2 {
		return true, nil
	}

	return false, err
}

// Apply applies new Terraform configurations. In case any
// changes are detected, user confirmation is required.
func (t *terraform) Apply() error {
	changes, err := t.Plan()

	if err != nil {
		return err
	}

	// Ask user for permission if there are any changes
	if changes && t.showPlan {
		err := ui.Ask("Proceed with terraform apply?")

		if err != nil {
			return err
		}
	}

	args := []string{
		flag("auto-approve"),
		flag("input", false),
		flag("lock", true),
		flag("lock-timeout", "0s"),
		flag("parallelism", 10),
		flag("refresh", true),
	}

	_, err = t.runCmd("apply", args, true)
	return err
}

// Destroy destroys the Terraform project.
func (t *terraform) Destroy() error {
	err := t.init()
	if err != nil {
		return err
	}

	args := []string{
		flag("auto-approve"),
		flag("input", false),
		flag("lock", true),
		flag("lock-timeout", "0s"),
		flag("parallelism", 10),
		flag("refresh", true),
	}

	_, err = t.runCmd("destroy", args, true)

	return err
}

// findOrInstall first searches for Terraform binary locally and
// if binary is not found, it is installed in given binDir.
func (t *terraform) findOrInstall() (string, error) {
	var binPath string
	var err error

	ui.Printf(ui.INFO, "Ensuring Terraform %s is installed...\n", t.version)

	binPath, err = findTerraform(t.version, t.binDir)

	if err == nil {
		ui.Printf(ui.INFO, "Terraform %s found locally (%s).\n", t.version, binPath)
		return binPath, nil
	}

	ui.Printf(ui.INFO, "Terraform %s could not be found locally.\n", t.version)
	ui.Printf(ui.INFO, "Installing Terraform %s in '%s'...\n", t.version, t.binDir)

	binPath, err = installTerraform(t.version, t.binDir)

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

// installTerraform installs Terraform binary of the provided
// version in a given directory.
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

// runCmd runs terraform command and returns exit code with
// a potential error.
func (t *terraform) runCmd(action string, args []string, showOutput bool) (int, error) {
	args = append([]string{action}, args...)

	if !ui.HasColor() {
		args = append(args, flag("no-color"))
	}

	cmd := exec.Command(t.binPath, args...)
	cmd.Dir = t.projectDir

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}

	if showOutput || ui.Debug() {
		cmd.Stdout = ui.Streams().Out().File()
		cmd.Stderr = ui.Streams().Err().File()
	}

	err := cmd.Run()
	exitCode := cmd.ProcessState.ExitCode()

	if err != nil {
		err = fmt.Errorf("terraform %s failed: %v", action, err)
	}

	return exitCode, err
}

// Flag concatenates key and value with "=" if value is provided.
func flag(key string, value ...interface{}) string {
	if len(value) > 0 && value[0] != nil {
		return fmt.Sprintf("-%s=%v", key, value[0])
	}

	return fmt.Sprintf("-%s", key)
}
