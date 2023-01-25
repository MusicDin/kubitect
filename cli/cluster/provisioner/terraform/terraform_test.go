package terraform

import (
	"cli/env"
	"cli/ui"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockTerraform(t *testing.T) *terraform {
	tmpDir := t.TempDir()

	return &terraform{
		Version:    env.ConstTerraformVersion,
		BinDir:     tmpDir,
		WorkingDir: tmpDir,
		Ui:         ui.MockUi(t),
	}
}

func TestTerraformInit(t *testing.T) {
	tf := MockTerraform(t)

	err := tf.init()
	assert.NoError(t, err)

	// Find Terraform locally
	err = tf.init()
	assert.NoError(t, err)
}
