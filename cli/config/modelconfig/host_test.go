package modelconfig

import (
	"cli/utils/defaults"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHost_Empty(t *testing.T) {
	assert.ErrorContains(t, Host{}.Validate(), "Field 'name' is required and cannot be empty.")
	assert.ErrorContains(t, Host{}.Validate(), "Field 'type' is required and cannot be empty.")
}

func TestHost_Default(t *testing.T) {
	assert.ErrorContains(t, Host{}.Validate(), "Field 'name' is required and cannot be empty.")
	assert.ErrorContains(t, defaults.Assign(&Host{}).Validate(), "Field 'type' is required and cannot be empty.")
}

func TestHost(t *testing.T) {
	h1 := Host{
		Name: "host",
		Connection: Connection{
			Type: LOCAL,
		},
	}

	h2 := Host{
		Name: "host",
	}

	h3 := Host{
		Name: "host",
		DataResourcePools: []DataResourcePool{
			{
				Name: "drp1",
			},
			{
				Name: "drp1",
			},
		},
	}

	assert.NoError(t, h1.Validate())
	assert.EqualError(t, defaults.Assign(&h2).Validate(), "Field 'type' is required and cannot be empty.")
	assert.ErrorContains(t, h3.Validate(), "Field 'type' is required and cannot be empty.")
	assert.ErrorContains(t, h3.Validate(), "Field 'Name' must be unique for each element in 'dataResourcePools'.")
}
