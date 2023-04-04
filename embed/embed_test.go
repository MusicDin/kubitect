package embed

import (
	"testing"

	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/stretchr/testify/assert"
)

func TestGetTemplate(t *testing.T) {
	tpl, err := GetTemplate("etcd.yaml.tpl")
	assert.NoError(t, err)
	assert.NotNil(t, tpl)
}

func TestMirrorResource(t *testing.T) {
	for _, f := range env.ProjectRequiredFiles {
		err := MirrorResource(f, t.TempDir())
		assert.NoError(t, err)
	}
}

func TestMirrorResource_InvalidResourcePath(t *testing.T) {
	err := MirrorResource("invalid", t.TempDir())
	assert.ErrorContains(t, err, "resources/invalid: file does not exist")
}

func TestGetResource_File(t *testing.T) {
	resPath := "terraform/variables.tf"
	res, err := GetResource(resPath)
	assert.NoError(t, err)
	assert.Equal(t, res.Name, "variables.tf")
	assert.Equal(t, res.Path, resPath)
	assert.NotEmpty(t, res.Content)
}

func TestGetResources_Invalid(t *testing.T) {
	resPath := "terraform/modules"
	_, err := GetResource(resPath)
	assert.ErrorContains(t, err, "resources/terraform/modules: is a directory")
}

func TestPresets(t *testing.T) {
	expect := []string{
		"presets/minimal.yaml",
		"presets/getting-started.yaml",
	}

	presets, err := Presets()
	assert.NoError(t, err)

	var paths []string
	for _, p := range presets {
		paths = append(paths, p.Path)
	}

	assert.Subset(t, paths, expect)
}

func TestGetPreset(t *testing.T) {
	p, err := GetPreset("minimal.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "minimal.yaml", p.Name)
	assert.Equal(t, "minimal.yaml", p.Path)
}

func TestGetPreset_NotFound(t *testing.T) {
	p, err := GetPreset("invalid")
	assert.Nil(t, p)
	assert.ErrorContains(t, err, "presets/invalid: file does not exist")
}
