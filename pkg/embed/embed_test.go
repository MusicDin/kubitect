package embed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	p, err := GetPreset("minimal")
	assert.NoError(t, err)
	assert.Equal(t, "minimal", p.Name)
	assert.Equal(t, "presets/minimal.yaml", p.Path)
}

func TestGetPreset_NotFound(t *testing.T) {
	p, err := GetPreset("invalid")
	assert.Nil(t, p)
	assert.EqualError(t, err, "get preset: preset invalid not found")
}
