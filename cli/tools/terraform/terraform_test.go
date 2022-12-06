package terraform

import (
	"cli/env"
	"cli/ui"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerraformInit(t *testing.T) {
	tmpDir := t.TempDir()

	tf := &Terraform{
		Version:    env.ConstTerraformVersion,
		BinDir:     tmpDir,
		WorkingDir: tmpDir,
		Ui:         ui.MockUi(t),
	}

	err := tf.init()
	assert.NoError(t, err)

	// Find Terraform locally
	err = tf.init()
	assert.NoError(t, err)
}
