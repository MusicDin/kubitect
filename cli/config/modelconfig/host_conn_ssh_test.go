package modelconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnSSH(t *testing.T) {
	pk := File("host_conn_ssh_test.go")

	ssh := ConnectionSSH{
		Keyfile: &pk,
	}

	assert.NoError(t, ssh.Validate())
	assert.ErrorContains(t, ConnectionSSH{}.Validate(), "Path to password-less private key of the remote host is required.")
}
