package modelconfig

import (
	"github.com/MusicDin/kubitect/cli/pkg/utils/defaults"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataResPool_Empty(t *testing.T) {
	assert.ErrorContains(t, DataResourcePool{}.Validate(), "Field 'name' is required and cannot be empty.")
	assert.ErrorContains(t, DataResourcePool{}.Validate(), "Field 'path' is required and cannot be empty.")
}

func TestDataResPool_Default(t *testing.T) {
	assert.EqualError(t, defaults.Assign(&DataResourcePool{}).Validate(), "Field 'name' is required and cannot be empty.")
}

func TestDataResPool(t *testing.T) {
	drp1 := DataResourcePool{
		Name: "test",
		Path: "/path",
	}

	assert.NoError(t, drp1.Validate())
}

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

func TestHost_Mock(t *testing.T) {
	assert.NoError(t, MockLocalHost(t, "test", true).Validate())
	assert.NoError(t, MockRemoteHost(t, "test", true, false).Validate())
}
