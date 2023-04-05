package config

import (
	"testing"

	"github.com/MusicDin/kubitect/pkg/utils/defaults"

	"github.com/stretchr/testify/assert"
)

func TestNetBridge(t *testing.T) {
	assert.Error(t, NetworkBridge(" ").Validate())
	assert.Error(t, NetworkBridge("br-0").Validate())
	assert.Error(t, NetworkBridge("br_0").Validate())
	assert.Error(t, NetworkBridge("1234567890abcdefg").Validate()) // longer than 16 chars
	assert.NoError(t, NetworkBridge("").Validate())
	assert.NoError(t, NetworkBridge("br0").Validate())
}

func TestNetMode(t *testing.T) {
	assert.Error(t, NetworkMode("wrong").Validate())
	assert.NoError(t, NetworkMode("nat").Validate())
	assert.NoError(t, NetworkMode("bridge").Validate())
	assert.NoError(t, NetworkMode("route").Validate())
	assert.NoError(t, NetworkMode(NAT).Validate())
}

func TestNetwork(t *testing.T) {
	br := NetworkBridge("br0")
	cidr := CIDRv4("192.168.113.0/20")
	mode := BRIDGE

	net1 := Network{
		CIDR: cidr,
	}

	net2 := Network{
		CIDR:   cidr,
		Mode:   mode,
		Bridge: br,
	}

	net3 := Network{
		CIDR: cidr,
		Mode: mode,
	}

	assert.NoError(t, defaults.Assign(&net1).Validate())
	assert.NoError(t, defaults.Assign(&net2).Validate())
	assert.EqualError(t, defaults.Assign(&net3).Validate(), "Field 'bridge' is required when network mode is set to 'bridge'.")
	assert.ErrorContains(t, Network{}.Validate(), "Field 'cidr' is required and cannot be empty.")
	assert.ErrorContains(t, Network{}.Validate(), "Field 'mode' must be one of the following values")
}
