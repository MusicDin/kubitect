package terraform

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/MusicDin/kubitect/pkg/cluster/event"
	"github.com/MusicDin/kubitect/pkg/config/modelconfig"
	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/MusicDin/kubitect/pkg/ui"
	"github.com/MusicDin/kubitect/pkg/utils/cmp"

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
	hosts := []modelconfig.Host{
		modelconfig.MockLocalHost(t, "test1", false),
		modelconfig.MockLocalHost(t, "test2", true),
		modelconfig.MockRemoteHost(t, "test3", false, false),
	}

	cfg := &modelconfig.Config{Hosts: hosts}

	prov := NewTerraformProvisioner(clsPath(t), "shared/path", true, cfg)
	assert.NoError(t, prov.Init(nil))
}

func TestNewTerraformProvisioner_InvalidHosts(t *testing.T) {
	cfg := &modelconfig.Config{}

	prov := NewTerraformProvisioner(clsPath(t), "shared/path", true, cfg)
	assert.ErrorContains(t, prov.Init(nil), "hosts list is empty")
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

func TestExtractRemovedHosts(t *testing.T) {
	host := modelconfig.Host{
		Name: "test",
	}

	change := []cmp.Change{
		{
			Type:   reflect.TypeOf(host),
			Before: host,
			After:  nil,
		},
	}

	events := []event.Event{
		event.MockEvent(t, event.OK, cmp.DELETE, change),
		event.MockEvent(t, event.OK, cmp.CREATE, change),
	}

	hosts := extractRemovedHosts(events)
	assert.Len(t, hosts, 1)
}
