package terraform

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

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

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}

	cmd.Stderr = ui.Streams().Err().File()
	if showOutput || ui.Debug() {
		cmd.Stdout = ui.Streams().Out().File()
	}

	cmd.Env = []string{fmt.Sprintf("PATH=%s", os.Getenv("PATH"))}
	if ui.Debug() {
		cmd.Env = append(cmd.Env, "TF_LOG=INFO")
	}
	if os.Getenv("HTTPS_PROXY") != "" || os.Getenv("https_proxy") != "" {
		proxyValue := os.Getenv("HTTPS_PROXY")
		if proxyValue == "" {
			proxyValue = os.Getenv("https_proxy")
		}
		cmd.Env = append(cmd.Env, fmt.Sprintf("HTTPS_PROXY=%s", proxyValue))
	}
	err := cmd.Run()
	exitCode := cmd.ProcessState.ExitCode()

	if err != nil {
		err = fmt.Errorf("terraform %s failed: %v", action, err)
	}

	return exitCode, err
}
