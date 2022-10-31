package modelinfra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Empty(t *testing.T) {
	assert.EqualError(t, Config{}.Validate(), "Terraform produced invalid output.")
}
