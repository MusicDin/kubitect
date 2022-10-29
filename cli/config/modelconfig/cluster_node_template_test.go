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
	assert.EqualError(t, OSNetworkInterface("").Validate(), "Field can contain only alphanumeric characters. (actual: )")
	assert.EqualError(t, OSNetworkInterface("1234567890abcdefg").Validate(), "Maximum length of the field is 16 (actual: 1234567890abcdefg)")
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
	assert.EqualError(t, nts2.Validate(), "Field 'privateKeyPath' must be a valid file path that points to an existing file. (actual: ./non-existing)")
}
