package terraform

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraformCmdEnv_PreservesHomeAndAddsTFLog(t *testing.T) {
	oldHome := os.Getenv("HOME")
	oldPath := os.Getenv("PATH")
	defer func() {
		require.NoError(t, os.Setenv("HOME", oldHome))
		require.NoError(t, os.Setenv("PATH", oldPath))
	}()

	require.NoError(t, os.Setenv("HOME", "/tmp/test-home"))
	require.NoError(t, os.Setenv("PATH", "/tmp/test-bin"))

	env := terraformCmdEnv(true)
	envMap := envSliceToMap(env)

	assert.Equal(t, "/tmp/test-home", envMap["HOME"])
	assert.Equal(t, "/tmp/test-bin", envMap["PATH"])
	assert.Equal(t, "INFO", envMap["TF_LOG"])
}

func TestTerraformCmdEnv_OmitsTFLogWhenDebugDisabled(t *testing.T) {
	env := terraformCmdEnv(false)
	envMap := envSliceToMap(env)

	_, ok := envMap["TF_LOG"]
	assert.False(t, ok)
}

func envSliceToMap(env []string) map[string]string {
	result := make(map[string]string, len(env))

	for _, item := range env {
		key, value, ok := splitEnvVar(item)
		if !ok {
			continue
		}

		result[key] = value
	}

	return result
}
