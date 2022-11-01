package virtualenv

import (
	"cli/env"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type VirtualEnvironment struct {
	Name             string
	Path             string
	RequirementsPath string
}

var Env = struct {
	Main      VirtualEnvironment
	Kubespray VirtualEnvironment
}{
	Main: VirtualEnvironment{
		Name:             "main",
		RequirementsPath: "ansible/kubitect/requirements.txt",
	},
	Kubespray: VirtualEnvironment{
		Name:             "kubespray",
		RequirementsPath: "ansible/kubespray/requirements.txt",
	},
}

// Setup creates virtual environment in the cluster path
// and installs required pip3 and ansible dependencies.
func (ve *VirtualEnvironment) Setup(clusterPath, version string) error {
	fmt.Printf("Setting up '%s' virtual environment...\n", ve.Name)

	ve.Path = env.VenvDirPath(ve.Name, version)

	err := ve.create(clusterPath)

	if err != nil {
		return err
	}

	err = ve.installPipReq(clusterPath)

	if err != nil {
		return err
	}

	return nil
}

// create creates virtual environment if it does not yet exist.
func (ve VirtualEnvironment) create(clusterPath string) error {
	fmt.Println("Creating virtual environment...")

	cmd := exec.Command("virtualenv", "-p", "python3", ve.Path)
	cmd.Dir = clusterPath

	if env.DebugMode {
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
func (ve VirtualEnvironment) installPipReq(clusterPath string) error {
	fmt.Println("Installing pip3 dependencies...")
	fmt.Println("This can take up to a minute when the virtual environment is initialized for the first time...")

	cmd := exec.Command("pip3", "install", "-r", ve.RequirementsPath)
	cmd.Path = filepath.Join(ve.Path, "bin", "pip3")
	cmd.Dir = clusterPath

	if env.DebugMode {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("Failed to install pip3 requirements: %v", err)
	}

	return nil
}
