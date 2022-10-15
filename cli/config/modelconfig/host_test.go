package modelconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHost(t *testing.T) {

	connType := LOCAL
	name := "test"

	h1 := Host{
		Name: &name,
		Connection: &Connection{
			Type: &connType,
		},
	}

	h2 := Host{
		Name: &name,
	}

	h3 := Host{
		Name: &name,
		DataResourcePools: &[]DataResourcePool{
			{
				Name: &name,
			},
			{
				Name: &name,
			},
		},
	}

	assert.NoError(t, h1.Validate())
	assert.ErrorContains(t, h2.Validate(), "Field 'connection' is required.")
	assert.ErrorContains(t, h3.Validate(), "Field 'Name' must be unique for each element in 'dataResourcePools'.")
	assert.ErrorContains(t, Host{}.Validate(), "Field 'name' is required.")
	assert.ErrorContains(t, Host{}.Validate(), "Field 'connection' is required.")
}
