package modelconfig

import (
	"cli/utils/defaults"
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
	source := OSSource("./cluster_node_template_test.go")

	os1 := OS{
		Distro: distro,
	}

	os2 := OS{
		Source: source,
	}

	assert.ErrorContains(t, OS{}.Validate(), "Field 'distro' must be one of the following values: [ubuntu|ubuntu20|ubuntu22|debian|debian11]")
	assert.ErrorContains(t, OS{}.Validate(), "Field 'networkInterface' can contain only alphanumeric characters.")
	assert.NoError(t, defaults.Assign(&os1).Validate())
	assert.NoError(t, defaults.Assign(&os2).Validate())
}

func TestNodeTemplateSSH(t *testing.T) {
	nts1 := NodeTemplateSSH{
		PrivateKeyPath: File("./cluster_node_template_test.go"),
	}

	nts2 := NodeTemplateSSH{
		PrivateKeyPath: File("./non-existing"),
	}

	assert.NoError(t, NodeTemplateSSH{}.Validate())
	assert.NoError(t, nts1.Validate())
	assert.EqualError(t, nts2.Validate(), "Field 'privateKeyPath' must be a valid file path that points to an existing file. (actual: ./non-existing)")
}
