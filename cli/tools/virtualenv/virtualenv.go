package virtualenv

import (
	"cli/ui"
	"fmt"
	"os/exec"
	"path/filepath"
)

type (
	VirtualEnv interface {
		Init() error
		Path() string
	}

	virtualEnv struct {
		name             string
		path             string
		workingDir       string
		requirementsPath string
		initialized      bool
	}
)

func NewVirtualEnv(name, path, workingDir, reqPath string) VirtualEnv {
	return &virtualEnv{
		name:             name,
		path:             path,
		workingDir:       workingDir,
		requirementsPath: reqPath,
	}
}

func (e *virtualEnv) Path() string {
	return e.path
}

// Init creates virtual environment in the cluster path
// and installs required pip3 and ansible dependencies.
func (e *virtualEnv) Init() error {
	if e.initialized {
		return nil
	}

	ui.Printf(ui.INFO, "Setting up '%s' virtual environment...\n", e.name)
	ui.Println(ui.INFO, "Creating virtual environment...")

	if err := e.create(); err != nil {
		return err
	}

	ui.Println(ui.INFO, "Installing pip3 dependencies...")
	ui.Println(ui.INFO, "This can take up to a minute when the virtual environment is initialized for the first time...")

	if err := e.installPipReq(); err != nil {
		return err
	}

	e.initialized = true

	return nil
}

// create creates virtual environment if it does not yet exist.
func (e *virtualEnv) create() error {
	cmd := exec.Command("virtualenv", "-p", "python3", e.path)
	cmd.Dir = e.workingDir

	if ui.Debug() {
		cmd.Stdout = ui.Streams().Out().File()
		cmd.Stderr = ui.Streams().Err().File()
	}

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to create virtual environment: %v", err)
	}

	return nil
}

// installPipReq installs pip3 requirements into virtual environment.
func (e *virtualEnv) installPipReq() error {
	cmd := exec.Command("pip3", "install", "-r", e.requirementsPath)
	cmd.Path = filepath.Join(e.path, "bin", "pip3")
	cmd.Dir = e.workingDir

	if ui.Debug() {
		cmd.Stdout = ui.Streams().Out().File()
		cmd.Stderr = ui.Streams().Err().File()
	}

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to install pip3 requirements: %v", err)
	}

	return nil
}
