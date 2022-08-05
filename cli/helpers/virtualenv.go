package helpers

import (
	"cli/env"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
		RequirementsPath: "ansible/kubitect/requirements.txt",
	},
	Kubespray: VirtualEnvironment{
		Name:             "kubespray",
		RequirementsPath: "ansible/kubespray/requirements.txt",
	},
}

// SetupVirtualEnvironment creates virtual environment in the cluster path
// and installs required pip3 and ansible dependencies.
func SetupVirtualEnvironment(clusterPath string, venv VirtualEnvironment) error {

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

	venvPath := filepath.Join(env.ConstVenvBinDir, venvName)

	cmd := exec.Command("virtualenv", "-p", "python3", venvPath)
	cmd.Dir = clusterPath

	if env.DebugMode {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to create virtual environment: %w", err)
	}

	return nil
}

// installPipRequirements installs pip3 requirements into virtual envrionment.
func installPipRequirements(clusterPath string, venv VirtualEnvironment) error {

	fmt.Println("Installing pip3 dependencies...")
	fmt.Println("This can take up to a minute when the virtual environment is initialized for the first time...")

	cmd := exec.Command("pip3", "install", "-r", venv.RequirementsPath)
	cmd.Path = filepath.Join(clusterPath, env.ConstVenvBinDir, venv.Name, "bin", "pip3")
	cmd.Dir = clusterPath

	if env.DebugMode {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to install pip3 requirements: %w", err)
	}

	return nil
}
