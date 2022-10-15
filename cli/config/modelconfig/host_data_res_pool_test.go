package modelconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataResPool(t *testing.T) {
	name := "test"

	drp1 := DataResourcePool{
		Name: &name,
	}

	assert.NoError(t, drp1.Validate())
	assert.ErrorContains(t, DataResourcePool{}.Validate(), "Field 'name' is required.")
}
