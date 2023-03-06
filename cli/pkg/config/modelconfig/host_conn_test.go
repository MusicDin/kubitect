package modelconfig

import (
	"github.com/MusicDin/kubitect/cli/pkg/utils/defaults"
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

func TestConnSSH_Empty(t *testing.T) {
	assert.ErrorContains(t, ConnectionSSH{}.Validate(), "Path to password-less private key of the remote host is required.")
	assert.ErrorContains(t, ConnectionSSH{}.Validate(), "Minimum value for field 'port' is 1 (actual: 0).")
}

func TestConnSSH_Default(t *testing.T) {
	assert.EqualError(t, defaults.Assign(&ConnectionSSH{}).Validate(), "Path to password-less private key of the remote host is required.")
}

func TestConnSSH(t *testing.T) {
	ssh := ConnectionSSH{
		Keyfile: File("host_conn_test.go"),
	}

	assert.NoError(t, defaults.Assign(&ssh).Validate())
}

func TestConn_Empty(t *testing.T) {
	assert.EqualError(t, Connection{}.Validate(), "Field 'type' is required and cannot be empty.")
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
