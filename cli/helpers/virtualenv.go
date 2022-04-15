package helpers

import (
	"cli/env"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	venvBinDir = "bin/venvs"
)

type VirtualEnvironment struct {
	Name             string
	RequirementsPath string
}

var Venvs = struct {
	Main      VirtualEnvironment
	Kubespray VirtualEnvironment
}{
	Main: VirtualEnvironment{
		Name:             "main",
		RequirementsPath: "requirements.txt",
	},
	Kubespray: VirtualEnvironment{
		Name:             "kubespray",
		RequirementsPath: "ansible/kubespray/requirements.txt",
	},
}

// var (
// 	MainVenv = &VirtualEnvironment{
// 		Name:             "main",
// 		RequirementsPath: "requirements.txt",
// 	}
// 	KubesprayVenv = &VirtualEnvironment{
// 		Name:             "kubespray",
// 		RequirementsPath: "ansible/kubespray/requirements.txt",
// 	}
// )

// setupVirtualEnironment creates virtual environment in the cluster path
// and installs required pip3 and ansible dependencies.
func SetupVirtualEnironment(clusterPath string, venv VirtualEnvironment) error {

	fmt.Printf("Setting up '%s' virtual environment...\n", venv.Name)

	var err error

	err = createVirtualEnvironment(clusterPath, venv.Name)
	if err != nil {
		return err
	}

	err = installPipRequirements(clusterPath, venv)
	if err != nil {
		return err
	}

	return nil
}

// createVirtualEnvironment creates virtual environment if it does not yet exist.
func createVirtualEnvironment(clusterPath string, venvName string) error {

	fmt.Println("Creating virtual environment...")

	venvPath := filepath.Join(venvBinDir, venvName)

	cmd := exec.Command("virtualenv", "-p", "python3", venvPath)
	cmd.Dir = clusterPath

	if env.DebugMode {
		cmd.Stdout = os.Stdout
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to create virutal environment: %w", err)
	}

	return nil
}

// installPipRequirements installs pip3 requirements into virtual envrionment.
func installPipRequirements(clusterPath string, venv VirtualEnvironment) error {

	fmt.Println("Installing pip3 dependencies...")
	fmt.Println("This can take up to a minute when the virtual environment is initialized for the first time...")

	cmd := exec.Command("pip3", "install", "-r", venv.RequirementsPath)
	cmd.Path = filepath.Join(clusterPath, venvBinDir, venv.Name, "bin", "pip3")
	cmd.Dir = clusterPath

	if env.DebugMode {
		cmd.Stdout = os.Stdout
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to install pip3 requirements: %w", err)
	}

	return nil
}
