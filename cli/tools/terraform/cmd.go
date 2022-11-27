package terraform

import (
	"cli/env"
	"cli/ui"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

type TerraformCmd struct {
	showOutput bool

	path       string
	projectDir string
	action     string
	args       map[string]string
}

func (t *Terraform) NewCmd(action string) *TerraformCmd {
	return &TerraformCmd{
		showOutput: true,
		path:       t.binPath,
		projectDir: t.projectDir,
		action:     action,
		args:       make(map[string]string),
	}
}

func (c *TerraformCmd) HasOutput(b bool) {
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
	if env.NoColor {
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
	cmd := exec.Command(c.path, c.Args()...)
	cmd.Dir = c.projectDir

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}

	if c.showOutput {
		cmd.Stdout = ui.GlobalUi().Streams.Out.File
		cmd.Stderr = ui.GlobalUi().Streams.Err.File
	}

	err := cmd.Run()

	code := cmd.ProcessState.ExitCode()

	if err != nil {
		return code, fmt.Errorf("terraform: failed to %v", c.action)
	}

	return code, nil
}
