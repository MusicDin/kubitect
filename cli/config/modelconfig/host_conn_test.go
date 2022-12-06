package modelconfig

import (
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
	local := LOCAL
	remote := REMOTE
	ip := IPv4("192.168.113.13")
	user := User("user")
	pk := File("host_conn_test.go")

	c1 := Connection{
		Type: &local,
	}

	c2 := Connection{
		Type: &remote,
		IP:   &ip,
		User: &user,
		SSH: ConnectionSSH{
			Keyfile: &pk,
		},
	}

	c3 := Connection{
		Type: &local,
	}

	c4 := Connection{
		Type: &remote,
	}

	assert.NoError(t, c1.Validate())
	assert.NoError(t, c2.Validate())
	assert.NoError(t, c3.Validate())
	assert.ErrorContains(t, c4.Validate(), "Field 'ip' is required when connection type is set to 'remote'.")
	assert.ErrorContains(t, c4.Validate(), "Field 'user' is required when connection type is set to 'remote'.")
	assert.ErrorContains(t, c4.Validate(), "Field 'ssh' is required when connection type is set to 'remote'.")
	assert.EqualError(t, Connection{}.Validate(), "Field 'type' is required.")
}
