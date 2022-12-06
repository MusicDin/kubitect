package virtualenv

import (
	"cli/ui"
	"fmt"
	"os/exec"
	"path/filepath"
)

type VirtualEnv struct {
	Name             string
	Path             string
	WorkingDir       string
	RequirementsPath string
	initialized      bool
	Ui               *ui.Ui
}

// Init creates virtual environment in the cluster path
// and installs required pip3 and ansible dependencies.
func (e *VirtualEnv) Init() error {
	if e.initialized {
		return nil
	}

	e.Ui.Printf(ui.INFO, "Setting up '%s' virtual environment...\n", e.Name)

	if err := e.create(); err != nil {
		return err
	}

	if err := e.installPipReq(); err != nil {
		return err
	}

	e.initialized = true

	return nil
}

// create creates virtual environment if it does not yet exist.
func (e *VirtualEnv) create() error {
	e.Ui.Println(ui.INFO, "Creating virtual environment...")

	cmd := exec.Command("virtualenv", "-p", "python3", e.Path)
	cmd.Dir = e.WorkingDir

	if e.Ui.Debug {
		cmd.Stdout = e.Ui.Streams.Out.File
		cmd.Stderr = e.Ui.Streams.Err.File
	}

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("Failed to create virtual environment: %v", err)
	}

	return nil
}

// installPipReq installs pip3 requirements into virtual environment.
func (e *VirtualEnv) installPipReq() error {
	e.Ui.Println(ui.INFO, "Installing pip3 dependencies...")
	e.Ui.Println(ui.INFO, "This can take up to a minute when the virtual environment is initialized for the first time...")

	cmd := exec.Command("pip3", "install", "-r", e.RequirementsPath)
	cmd.Path = filepath.Join(e.Path, "bin", "pip3")
	cmd.Dir = e.WorkingDir

	if e.Ui.Debug {
		cmd.Stdout = e.Ui.Streams.Out.File
		cmd.Stderr = e.Ui.Streams.Err.File
	}

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("Failed to install pip3 requirements: %v", err)
	}

	return nil
}
