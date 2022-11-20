package virtualenv

import (
	"cli/env"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type VirtualEnvType string

const (
	MAIN      VirtualEnvType = "main"
	KUBESPRAY VirtualEnvType = "kubespray"
)

type VirtualEnv struct {
	Name             string
	Path             string
	ClusterPath      string
	RequirementsPath string
	initialized      bool
}

var virtualEnvs map[VirtualEnvType]*VirtualEnv

func Set(t VirtualEnvType, v *VirtualEnv) {
	if virtualEnvs == nil {
		virtualEnvs = make(map[VirtualEnvType]*VirtualEnv)
	}

	virtualEnvs[t] = v
}

func Get(t VirtualEnvType) (*VirtualEnv, error) {
	v, ok := virtualEnvs[t]

	if ok {
		return v, v.setup()
	}

	return nil, fmt.Errorf("Virtual environment %v does not exist!", t)
}

// setup creates virtual environment in the cluster path
// and installs required pip3 and ansible dependencies.
func (ve *VirtualEnv) setup() error {
	if ve.initialized {
		return nil
	}

	fmt.Printf("Setting up '%s' virtual environment...\n", ve.Name)

	if err := ve.create(); err != nil {
		return err
	}

	if err := ve.installPipReq(); err != nil {
		return err
	}

	ve.initialized = true

	return nil
}

// create creates virtual environment if it does not yet exist.
func (ve *VirtualEnv) create() error {
	fmt.Println("Creating virtual environment...")

	cmd := exec.Command("virtualenv", "-p", "python3", ve.Path)
	cmd.Dir = ve.ClusterPath

	if env.Debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("Failed to create virtual environment: %v", err)
	}

	return nil
}

// installPipReq installs pip3 requirements into virtual environment.
func (ve *VirtualEnv) installPipReq() error {
	fmt.Println("Installing pip3 dependencies...")
	fmt.Println("This can take up to a minute when the virtual environment is initialized for the first time...")

	cmd := exec.Command("pip3", "install", "-r", ve.RequirementsPath)
	cmd.Path = filepath.Join(ve.Path, "bin", "pip3")
	cmd.Dir = ve.ClusterPath

	if env.Debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("Failed to install pip3 requirements: %v", err)
	}

	return nil
}
