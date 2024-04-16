package virtualenv

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/MusicDin/kubitect/pkg/ui"
)

type VirtualEnv struct {
	path             string
	requirementsPath string
	initialized      bool
}

// NewVirtualEnv returns new virtual environment (VE). It expects VE path
// to point on the directory where VE is created and 'reqPath' to the
// requirements.txt file.
func NewVirtualEnv(virtualEnvPath, reqPath string) *VirtualEnv {
	return &VirtualEnv{
		path:             virtualEnvPath,
		requirementsPath: reqPath,
	}
}

// Init creates virtual environment in the cluster path
// and installs required pip3 and ansible dependencies.
func (e *VirtualEnv) Init() error {
	if e.initialized {
		return nil
	}

	ui.Println(ui.INFO, "Setting up virtual environment...")

	err := e.create()
	if err != nil {
		return err
	}

	ui.Println(ui.INFO, "Installing pip3 dependencies...")
	ui.Println(ui.INFO, "This can take up to a minute when the virtual environment is initialized for the first time...")

	err = e.installPipReq()
	if err != nil {
		return err
	}

	e.initialized = true

	return nil
}

// create creates virtual environment if it does not yet exist.
func (e *VirtualEnv) create() error {
	wd := path.Dir(e.path)

	err := os.MkdirAll(wd, os.ModePerm)
	if err != nil {
		return err
	}

	cmd := exec.Command("virtualenv", "-p", "python3", e.path)
	cmd.Dir = wd

	if ui.Debug() {
		cmd.Stdout = ui.Streams().Out().File()
		cmd.Stderr = ui.Streams().Err().File()
	}

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create virtual environment: %v", err)
	}

	return nil
}

// installPipReq installs pip3 requirements into virtual environment.
func (e *VirtualEnv) installPipReq() error {
	cmd := exec.Command("pip3", "install", "-r", e.requirementsPath)
	cmd.Path = filepath.Join(e.path, "bin", "pip3")
	cmd.Dir = filepath.Dir(e.path)

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
