//go:build !linux
// +build !linux

package terraform

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/MusicDin/kubitect/pkg/ui"
)

// runCmd runs terraform command and returns exit code with
// a potential error.
func (t *terraform) runCmd(action string, args []string, showOutput bool) (int, error) {
	args = append([]string{action}, args...)

	if !ui.HasColor() {
		args = append(args, flag("no-color"))
	}

	cmd := exec.Command(t.binPath, args...)
	cmd.Dir = t.projectDir

	cmd.Stderr = ui.Streams().Err().File()
	if showOutput || ui.Debug() {
		cmd.Stdout = ui.Streams().Out().File()
	}

	cmd.Env = []string{fmt.Sprintf("PATH=%s", os.Getenv("PATH"))}
	if ui.Debug() {
		cmd.Env = append(cmd.Env, "TF_LOG=INFO")
	}

	err := cmd.Run()
	exitCode := cmd.ProcessState.ExitCode()

	if err != nil {
		err = fmt.Errorf("terraform %s failed: %v", action, err)
	}

	return exitCode, err
}
