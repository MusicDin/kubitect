package modelconfig

import (
	"testing"

	"github.com/MusicDin/kubitect/pkg/utils/defaults"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Mock(t *testing.T) {
	assert.NoError(t, MockConfig(t).Validate())
}

func TestConfig_Empty(t *testing.T) {
	assert.ErrorContains(t, Config{}.Validate(), "At least 1 host must be configured.")
	assert.ErrorContains(t, Config{}.Validate(), "Configuration must contain 'cluster' section.")
	assert.ErrorContains(t, Config{}.Validate(), "Configuration must contain 'kubernetes' section.")
}

func TestConfig_Valid(t *testing.T) {
	cfg := MockConfig(t)
	cfg.Cluster.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id:   "id",
					Host: "local",
					DataDisks: []DataDisk{
						// Correct pool reference
						{
							Name: "disk1",
							Pool: "pool1",
							Size: GB(5),
						},
						// Main pool reference (no direct pool ref)
						{
							Name: "disk2",
							Size: GB(5),
						},
					},
				},
			},
		},
	}

	assert.NoError(t, defaults.Assign(&cfg).Validate())
}

func TestConfig_InvalidIP(t *testing.T) {
	cfg := MockConfig(t)
	cfg.Cluster.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id:   "id",
					Host: "local",
					IP:   "192.168.114.13",
				},
			},
		},
	}

	assert.EqualError(t, defaults.Assign(&cfg).Validate(), "Field 'ip' must be a valid IP address within '192.168.113.0/24' subnet. (actual: 192.168.114.13)")
}

func TestConfig_MultipleDefaultHosts(t *testing.T) {
	cfg := MockConfig(t)
	cfg.Hosts = []Host{
		MockLocalHost(t, "local", true),
		MockRemoteHost(t, "remote", true, false),
	}

	assert.EqualError(t, defaults.Assign(&cfg).Validate(), "Only one host can be configured as default.")
}

func TestConfig_InvalidHostRef(t *testing.T) {
	cfg := MockConfig(t)
	cfg.Cluster.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id:   "id",
					Host: "wrong",
				},
			},
		},
	}

	assert.EqualError(t, defaults.Assign(&cfg).Validate(), "Field 'host' must point to one of the configured hosts: [local|remote] (actual: wrong)")
}

func TestConfig_InvalidPoolHostRef(t *testing.T) {
	cfg := MockConfig(t)
	cfg.Cluster.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id:   "id",
					Host: "remote",
					DataDisks: []DataDisk{
						{
							Name: "disk",
							Pool: "pool1",
							Size: GB(5),
						},
					},
				},
			},
		},
	}

	assert.EqualError(t, defaults.Assign(&cfg).Validate(), "Field 'pool' points to a data resource pool, but matching host 'remote' has none configured.")
}

func TestConfig_InvalidPoolRef(t *testing.T) {
	cfg := MockConfig(t)
	cfg.Cluster.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id: "id",
					DataDisks: []DataDisk{
						{
							Name: "disk",
							Pool: "wrong",
							Size: GB(5),
						},
					},
				},
			},
		},
	}

	assert.EqualError(t, defaults.Assign(&cfg).Validate(), "Field 'pool' must point to one of the pools configured on a matching host 'local': [pool1|pool2] (actual: wrong)")
}

func TestConfig_MainPoolRef(t *testing.T) {
	cfg := MockConfig(t)
	cfg.Cluster.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id: "id",
					DataDisks: []DataDisk{
						{
							Name: "disk1",
							Pool: "main",
							Size: GB(5),
						},
					},
				},
			},
		},
	}

	assert.NoError(t, defaults.Assign(&cfg).Validate())
}

// When disk pool matches a host with no hostname, no extra error should
// be thrown beside missing host name error.
func TestConfig_MissingHostName(t *testing.T) {
	cfg := MockConfig(t)
	cfg.Hosts = []Host{
		{
			Connection: Connection{
				Type: LOCAL,
			},
		},
	}

	assert.EqualError(t, defaults.Assign(&cfg).Validate(), "Field 'name' is required and cannot be empty.")
}
