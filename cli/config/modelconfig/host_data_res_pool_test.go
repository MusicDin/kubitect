package modelconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataResPool(t *testing.T) {
	drp1 := DataResourcePool{
		Name: "test",
		Path: "/path",
	}

	assert.NoError(t, drp1.Validate())
	assert.ErrorContains(t, DataResourcePool{}.Validate(), "Field 'name' is required and cannot be empty.")
	assert.ErrorContains(t, DataResourcePool{}.Validate(), "Field 'path' is required and cannot be empty.")
}
