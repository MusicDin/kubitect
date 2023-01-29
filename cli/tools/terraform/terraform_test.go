package terraform

import (
	"cli/env"
	"cli/ui"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockTerraform(t *testing.T) *terraform {
	ui.MockGlobalUi(t, ui.UiOptions{NoColor: true, Debug: true})

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
		env:        []string{binDir},
	}
}

func MockInvalidTerraform(t *testing.T) *terraform {
	tf := MockTerraform(t)

	// Create a file on binDir path
	tf.binDir = path.Join(tf.projectDir, "invalid")
	_, err := os.Create(tf.binDir)
	assert.NoError(t, err)

	return tf
}

func TestNewTerraform(t *testing.T) {
	tf := NewTerraform("v1.0.0", "/tmp", "/tmp", true)
	assert.NotNil(t, tf)
}

func TestTerraform_Init(t *testing.T) {
	tf := MockTerraform(t)
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

func TestTerraform_Init_InvalidBinDir(t *testing.T) {
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
