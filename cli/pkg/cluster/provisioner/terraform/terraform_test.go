package terraform

import (
	modelconfig2 "github.com/MusicDin/kubitect/cli/pkg/config/modelconfig"
	"github.com/MusicDin/kubitect/cli/pkg/env"
	ui2 "github.com/MusicDin/kubitect/cli/pkg/ui"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockTerraform(t *testing.T) *terraform {
	ui2.MockGlobalUi(t, ui2.UiOptions{NoColor: true, Debug: true})

	tmpDir := t.TempDir()

	binDir := path.Join(tmpDir, env.ConstTerraformVersion)
	projDir := tmpDir

	// Create sample terraform main.tf file
	maintf := "output \"test\" { value = \"test\" }"
	maintfPath := path.Join(projDir, "main.tf")
	err := ioutil.WriteFile(maintfPath, []byte(maintf), 0777)
	assert.NoError(t, err)

	return &terraform{
		version:    env.ConstTerraformVersion,
		binDir:     binDir,
		projectDir: projDir,
	}
}

func MockMissingTerraform(t *testing.T) *terraform {
	tf := MockTerraform(t)
	tf.version = "1.0.0"
	return tf
}

func MockInvalidTerraform(t *testing.T) *terraform {
	tf := MockMissingTerraform(t)

	// Create a file on binDir path
	tf.binDir = path.Join(tf.projectDir, "invalid")
	_, err := os.Create(tf.binDir)
	assert.NoError(t, err)

	return tf
}

func TestNewTerraformProvisioner(t *testing.T) {
	hosts := []modelconfig2.Host{
		modelconfig2.MockLocalHost(t, "test1", false),
		modelconfig2.MockLocalHost(t, "test2", true),
		modelconfig2.MockRemoteHost(t, "test3", false, false),
	}

	cfg := modelconfig2.Config{Hosts: hosts}

	prov := NewTerraformProvisioner(clsPath(t), "shared/path", true, &cfg)
	assert.NoError(t, prov.Init())
}

func TestNewTerraformProvisioner_InvalidHosts(t *testing.T) {
	cfg := modelconfig2.Config{}

	prov := NewTerraformProvisioner(clsPath(t), "shared/path", true, &cfg)
	assert.ErrorContains(t, prov.Init(), "hosts list is empty")
}

func TestTerraform_init(t *testing.T) {
	tf := MockMissingTerraform(t)
	tfPath := path.Join(tf.binDir, "terraform")

	assert.NoError(t, tf.init())
	assert.Equal(t, tfPath, tf.binPath)
	assert.Equal(t, true, tf.initialized)

	// tf.init() should quit immediately if initialized == true
	assert.NoError(t, tf.init())

	// Set initialized to false to call findAndInstall again.
	// Since tf is already installed, tf must be found locally
	tf.initialized = false
	assert.NoError(t, tf.init())
	assert.Equal(t, tfPath, tf.binPath)
	assert.Equal(t, true, tf.initialized)
}

func TestTerraform_init_InvalidBinDir(t *testing.T) {
	tf := MockInvalidTerraform(t)
	assert.ErrorContains(t, tf.init(), "not a directory")
}

func TestTerraform_Actions(t *testing.T) {
	tf := MockTerraform(t)

	_, err := tf.Plan()
	assert.NoError(t, err)
	assert.NoError(t, tf.Apply())
	assert.NoError(t, tf.Destroy())
}

func TestTerraform_Actions_Error(t *testing.T) {
	tf := MockInvalidTerraform(t)

	_, err := tf.Plan()
	assert.ErrorContains(t, err, "not a directory")
	assert.ErrorContains(t, tf.Apply(), "not a directory")
	assert.ErrorContains(t, tf.Destroy(), "not a directory")
}
