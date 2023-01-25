package terraform

import (
	"cli/ui"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

type TerraformCmd struct {
	binPath    string
	workingDir string
	action     string
	args       map[string]string

	showOutput bool
	ui         *ui.Ui
}

func (t *terraform) NewCmd(action string) *TerraformCmd {
	return &TerraformCmd{
		binPath:    t.binPath,
		workingDir: t.WorkingDir,
		action:     action,
		args:       make(map[string]string),
		showOutput: true,
		ui:         t.Ui,
	}
}

func (c *TerraformCmd) ShowOutput(b bool) {
	c.showOutput = b
}

func (c *TerraformCmd) AddArg(key string, value ...interface{}) {
	if key == "" {
		return
	}

	var v string

	if len(value) > 0 {
		v = fmt.Sprintf("%v", value[0])
	}

	c.args[key] = v
}

func (c *TerraformCmd) Args() []string {
	if c.ui.NoColor {
		c.AddArg("no-color")
	}

	args := []string{c.action}

	for k, v := range c.args {
		if !strings.HasPrefix(k, "-") {
			k = "-" + k
		}

		if v == "" {
			args = append(args, k)
		} else {
			args = append(args, fmt.Sprintf("%s=%v", k, v))
		}
	}

	return args
}

func (c *TerraformCmd) Run() (int, error) {
	cmd := exec.Command(c.binPath, c.Args()...)
	cmd.Dir = c.workingDir

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}

	if c.showOutput {
		cmd.Stdout = c.ui.Streams.Out.File
		cmd.Stderr = c.ui.Streams.Err.File
	}

	err := cmd.Run()

	code := cmd.ProcessState.ExitCode()

	if err != nil {
		return code, fmt.Errorf("terraform: failed to %v", c.action)
	}

	return code, nil
}
