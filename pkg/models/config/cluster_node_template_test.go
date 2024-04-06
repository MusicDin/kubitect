package config

import (
	"testing"

	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/MusicDin/kubitect/pkg/utils/defaults"

	"github.com/stretchr/testify/assert"
)

func TestOSDistro(t *testing.T) {
	assert.NoError(t, OSDistro(UBUNTU20).Validate())
	assert.NoError(t, OSDistro("ubuntu22").Validate())
	assert.NoError(t, OSDistro("debian12").Validate())
	assert.Error(t, OSDistro("invalid").Validate())
}

func TestOSDistro_Presets(t *testing.T) {
	for k, v := range env.ProjectOsPresets {
		assert.NoErrorf(t, OSDistro(k).Validate(), "OS preset %s represents an unknown OS distro", k)
		assert.NoErrorf(t, URL(v.Source).Validate(), "%s preset URL (%s) is invalid!", k, v)
	}
}

func TestOSNetworkInterface(t *testing.T) {
	assert.EqualError(t, OSNetworkInterface("").Validate(), "Field can contain only alphanumeric characters. (actual: )")
	assert.EqualError(t, OSNetworkInterface("1234567890abcdefg").Validate(), "Maximum length of the field is 16 (actual: 1234567890abcdefg)")
	assert.NoError(t, OSNetworkInterface("ens3").Validate())
}

func TestOS_Empty(t *testing.T) {
	assert.ErrorContains(t, OS{}.Validate(), "Field 'distro' must be one of the following values: [ubuntu20|")
	assert.ErrorContains(t, OS{}.Validate(), "Field 'networkInterface' can contain only alphanumeric characters.")
}

func TestOS_Defaults(t *testing.T) {
	os1 := OS{Distro: CENTOS9}
	os2 := OS{Distro: ROCKY9}
	os3 := OS{Source: OSSource("./cluster_node_template_test.go")}

	assert.NoError(t, defaults.Assign(&OS{}).Validate())
	assert.NoError(t, defaults.Assign(&os1).Validate())
	assert.NoError(t, defaults.Assign(&os2).Validate())
	assert.NoError(t, defaults.Assign(&os3).Validate())

	assert.Equal(t, CENTOS9, os1.Distro)
	assert.Equal(t, ROCKY9, os2.Distro)
	assert.Equal(t, UBUNTU22, os3.Distro)
	assert.Equal(t, "eth0", string(os1.NetworkInterface))
	assert.Equal(t, "eth0", string(os2.NetworkInterface))
	assert.Equal(t, "ens3", string(os3.NetworkInterface))
	assert.Equal(t, env.ProjectOsPresets["centos9"].Source, string(os1.Source))
	assert.Equal(t, env.ProjectOsPresets["rocky9"].Source, string(os2.Source))
	assert.Equal(t, "./cluster_node_template_test.go", string(os3.Source))
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
	assert.NoError(t, nts2.Validate())
}

func TestCpuMode(t *testing.T) {
	assert.NoError(t, CpuMode("custom").Validate())
	assert.NoError(t, CpuMode(HOST_PASSTHROUGH).Validate())
	assert.NoError(t, CpuMode("host-model").Validate())
	assert.NoError(t, CpuMode("maximum").Validate())
}
