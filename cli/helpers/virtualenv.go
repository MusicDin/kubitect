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

// PrepareVirtualEnironment creates virtual environment in the cluster path
// and installs required pip3 and ansible dependencies.
func PrepareVirtualEnironment(clusterPath string, venvName string, requirementsFile string) error {

	var err error

	err = createVirtualEnvironment(clusterPath, venvName)
	if err != nil {
		return err
	}

	err = installPipRequirements(clusterPath, venvName, requirementsFile)
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
func installPipRequirements(clusterPath string, venvName string, requirementsFile string) error {

	fmt.Println("Installing pip3 dependencies...")
	fmt.Println("This may take up to a minute if the virtual environment is initialized for the first time...")

	cmd := exec.Command("pip3", "install", "-r", requirementsFile)
	cmd.Path = filepath.Join(clusterPath, venvBinDir, venvName, "bin", "pip3")
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
