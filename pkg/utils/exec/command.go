package exec

import (
	"strings"
)

type Command struct {
	command    string
	args       []string
	envs       map[string]string
	workingDir string
}

func NewCommand(command string, args ...string) Command {
	if len(args) == 0 {
		// Split command by spaces, and treat it as command
		// with arguments. This prevents spaces in commands
		// but allows passing commands as a single string.
		split := strings.Split(command, " ")
		command = split[0]
		if len(split) > 1 {
			args = split[1:]
		}
	}

	return Command{
		command: command,
		args:    args,
	}
}

func (c Command) WithWorkingDir(workingDir string) Command {
	c.workingDir = workingDir
	return c
}

func (c Command) WithEnv(key string, value string) Command {
	if c.envs == nil {
		c.envs = make(map[string]string)
	}

	c.envs[key] = value
	return c
}

func (c Command) WithEnvMap(envs map[string]string) Command {
	if c.envs == nil {
		c.envs = make(map[string]string)
	}

	for k, v := range envs {
		c.envs[k] = v
	}

	return c
}
