package modelconfig

import (
	"testing"

	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/MusicDin/kubitect/pkg/utils/defaults"

	"github.com/stretchr/testify/assert"
)

func TestOSDistro(t *testing.T) {
	assert.NoError(t, OSDistro(UBUNTU).Validate())
	assert.NoError(t, OSDistro("ubuntu").Validate())
	assert.NoError(t, OSDistro("debian").Validate())
	assert.Error(t, OSDistro("wrong").Validate())
}

func TestOSDistro_Presets(t *testing.T) {
	for k, v := range env.ProjectOsPresets {
		assert.NoErrorf(t, OSDistro(k).Validate(), "OS preset %s represents an unknown OS distro", k)
		assert.NoErrorf(t, URL(v).Validate(), "%s preset URL (%s) is invalid!", k, v)
	}
}

func TestOSNetworkInterface(t *testing.T) {
	assert.EqualError(t, OSNetworkInterface("").Validate(), "Field can contain only alphanumeric characters. (actual: )")
	assert.EqualError(t, OSNetworkInterface("1234567890abcdefg").Validate(), "Maximum length of the field is 16 (actual: 1234567890abcdefg)")
	assert.NoError(t, OSNetworkInterface("ens3").Validate())
}

func TestOS_Empty(t *testing.T) {
	assert.ErrorContains(t, OS{}.Validate(), "Field 'distro' must be one of the following values: [ubuntu|ubuntu20|ubuntu22|debian|debian11]")
	assert.ErrorContains(t, OS{}.Validate(), "Field 'networkInterface' can contain only alphanumeric characters.")
}

func TestOS_Defaults(t *testing.T) {
	os1 := OS{Distro: UBUNTU}
	os2 := OS{Source: OSSource("./cluster_node_template_test.go")}

	assert.NoError(t, defaults.Assign(&OS{}).Validate())
	assert.NoError(t, defaults.Assign(&os1).Validate())
	assert.NoError(t, defaults.Assign(&os2).Validate())

	assert.Equal(t, UBUNTU, os1.Distro)
	assert.Equal(t, UBUNTU, os2.Distro)
	assert.Equal(t, "ens3", string(os1.NetworkInterface))
	assert.Equal(t, "ens3", string(os2.NetworkInterface))
	assert.Equal(t, env.ProjectOsPresets["ubuntu"], string(os1.Source))
	assert.Equal(t, "./cluster_node_template_test.go", string(os2.Source))
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
	// assert.EqualError(t, nts2.Validate(), "Field 'privateKeyPath' must be a valid file path that points to an existing file. (actual: ./non-existing)")
	assert.NoError(t, nts2.Validate())
}

func TestCpuMode(t *testing.T) {
	assert.NoError(t, CpuMode(PASSTHROUGH).Validate())
	assert.NoError(t, CpuMode("custom").Validate())
}
