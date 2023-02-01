package modelconfig

import (
	"cli/utils/defaults"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnType(t *testing.T) {
	assert.EqualError(t, ConnectionType("").Validate(), "Field must be one of the following values: [localhost|local|remote] (actual: ).")
	assert.EqualError(t, ConnectionType("wrong").Validate(), "Field must be one of the following values: [localhost|local|remote] (actual: wrong).")
	assert.NoError(t, ConnectionType("local").Validate())
	assert.NoError(t, ConnectionType("remote").Validate())
	assert.NoError(t, LOCALHOST.Validate())
}

func TestConn(t *testing.T) {
	c1 := Connection{
		Type: LOCAL,
	}

	c2 := Connection{
		Type: REMOTE,
		IP:   IPv4("192.168.113.13"),
		User: User("user"),
		SSH: ConnectionSSH{
			Keyfile: File("./host_conn_test.go"),
		},
	}

	c4 := Connection{
		Type: REMOTE,
	}

	assert.NoError(t, c1.Validate())
	assert.NoError(t, defaults.Assign(&c2).Validate())
	assert.ErrorContains(t, c4.Validate(), "Field 'ip' is required when connection type is set to 'remote'.")
	assert.ErrorContains(t, c4.Validate(), "Field 'user' is required when connection type is set to 'remote'.")
	assert.ErrorContains(t, c4.Validate(), "Field 'ssh' is required when connection type is set to 'remote'.")
	assert.EqualError(t, Connection{}.Validate(), "Field 'type' is required and cannot be empty.")
}
