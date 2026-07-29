package terraform

import (
	"fmt"
	"os"
)

func terraformCmdEnv(debug bool) []string {
	envMap := make(map[string]string)

	for _, env := range os.Environ() {
		key, value, ok := splitEnvVar(env)
		if !ok {
			continue
		}

		envMap[key] = value
	}

	envMap["PATH"] = os.Getenv("PATH")

	if debug {
		envMap["TF_LOG"] = "INFO"
	}

	env := make([]string, 0, len(envMap))
	for key, value := range envMap {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	return env
}

func splitEnvVar(env string) (string, string, bool) {
	for i := 0; i < len(env); i++ {
		if env[i] == '=' {
			return env[:i], env[i+1:], true
		}
	}

	return "", "", false
}
