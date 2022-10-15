package modelconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOSDistro(t *testing.T) {
	assert.NoError(t, OSDistro(UBUNTU).Validate())
	assert.NoError(t, OSDistro("ubuntu").Validate())
	assert.NoError(t, OSDistro("debian").Validate())
	assert.Error(t, OSDistro("wrong").Validate())
}

func TestOSNetworkInterface(t *testing.T) {
	assert.Error(t, OSNetworkInterface("").Validate())
	assert.Error(t, OSNetworkInterface("1234567890abcdefg").Validate()) // longer than 16 chars
	assert.NoError(t, OSNetworkInterface("ens3").Validate())
}

func TestOS(t *testing.T) {
	distro := UBUNTU

	os1 := OS{
		Distro: &distro,
	}

	assert.NoError(t, OS{}.Validate())
	assert.NoError(t, os1.Validate())
}

func TestNodeTemplateSSH(t *testing.T) {
	file1 := File("./cluster_node_template_test.go")
	file2 := File("./non-existing")

	nts1 := NodeTemplateSSH{
		PrivateKeyPath: &file1,
	}

	nts2 := NodeTemplateSSH{
		PrivateKeyPath: &file2,
	}

	assert.NoError(t, NodeTemplateSSH{}.Validate())
	assert.NoError(t, nts1.Validate())
	assert.Error(t, nts2.Validate())
}
