package modelconfig

import (
	"cli/utils/defaults"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnSSH_Default(t *testing.T) {
	assert.ErrorContains(t, ConnectionSSH{}.Validate(), "Path to password-less private key of the remote host is required.")
	assert.ErrorContains(t, ConnectionSSH{}.Validate(), "Minimum value for field 'port' is 1 (actual: 0).")
	assert.EqualError(t, defaults.Assign(&ConnectionSSH{}).Validate(), "Path to password-less private key of the remote host is required.")
}

func TestConnSSH(t *testing.T) {
	ssh := ConnectionSSH{
		Keyfile: File("host_conn_ssh_test.go"),
	}

	assert.NoError(t, defaults.Assign(&ssh).Validate())
}
