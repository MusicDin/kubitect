package modelconfig

import (
	"cli/utils/defaults"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sample_default = true
	sample_name1   = "test"
	sample_name2   = "test2"
	sample_dd_size = GB(5)
	sample_ip      = IPv4("192.168.114.13")

	net = Network{
		CIDR: CIDRv4("192.168.113.0/24"),
	}

	localhost = Host{
		Name:    "localhost",
		Default: true,
		Connection: Connection{
			Type: LOCAL,
		},
		DataResourcePools: []DataResourcePool{
			{
				Name: "pool1",
				Path: "/path1",
			},
			{
				Name: "pool2",
				Path: "/path2",
			},
		},
	}

	remotehost = Host{
		Name: "remotehost",
		Connection: Connection{
			Type: REMOTE,
			IP:   sample_ip,
			User: User("user"),
			SSH: ConnectionSSH{
				Keyfile: File("./config_test.go"),
			},
		},
	}

	k8s_version    = Version("v1.2.3")
	k8s_ks_version = MasterVersion(k8s_version)
	k8s            = Kubernetes{
		Version: k8s_version,
		Kubespray: Kubespray{
			Version: &k8s_ks_version,
		},
	}

	cluster = Cluster{
		Name:    sample_name1,
		Network: net,
	}

	config = Config{
		Hosts: []Host{
			localhost,
		},
		Kubernetes: k8s,
	}
)

func TestConfig_Empty(t *testing.T) {
	assert.ErrorContains(t, Config{}.Validate(), "At least 1 host must be configured.")
	assert.ErrorContains(t, Config{}.Validate(), "Configuration must contain 'cluster' section.")
	assert.ErrorContains(t, Config{}.Validate(), "Configuration must contain 'kubernetes' section.")
}

func TestConfig_Valid(t *testing.T) {
	cls := cluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id:   "id",
					Host: "localhost",
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

	cfg := config
	cfg.Cluster = cls

	assert.NoError(t, defaults.Assign(&cfg).Validate())
}

func TestConfig_InvalidIP(t *testing.T) {
	cls := cluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id:   "id",
					Host: "localhost",
					IP:   sample_ip,
				},
			},
		},
	}

	cfg := config
	cfg.Cluster = cls

	assert.EqualError(t, defaults.Assign(&cfg).Validate(), "Field 'ip' must be a valid IP address within '192.168.113.0/24' subnet. (actual: 192.168.114.13)")
}

func TestConfig_MultipleDefaultHosts(t *testing.T) {
	cls := cluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id:   "id",
					Host: "localhost",
				},
			},
		},
	}

	rhDef := remotehost
	rhDef.Default = true

	cfg := config
	cfg.Cluster = cls
	cfg.Hosts = []Host{
		localhost,
		rhDef,
	}

	assert.EqualError(t, defaults.Assign(&cfg).Validate(), "Only one host can be configured as default.")
}

func TestConfig_InvalidHostRef(t *testing.T) {
	cls := cluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id:   "id",
					Host: "wrong",
				},
			},
		},
	}

	cfg := config
	cfg.Cluster = cls

	assert.EqualError(t, defaults.Assign(&cfg).Validate(), "Field 'host' must point to one of the configured hosts: [localhost] (actual: wrong)")
}

func TestConfig_InvalidPoolHostRef(t *testing.T) {
	cls := cluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id:   "id",
					Host: "remotehost",
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

	cfg := config
	cfg.Cluster = cls
	cfg.Hosts = []Host{
		localhost,
		remotehost,
	}

	assert.EqualError(t, defaults.Assign(&cfg).Validate(), "Field 'pool' points to a data resource pool, but matching host 'remotehost' has none configured.")
}

func TestConfig_InvalidPoolRef(t *testing.T) {
	cls := cluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id: sample_name1,
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

	cfg := config
	cfg.Cluster = cls

	assert.EqualError(t, defaults.Assign(&cfg).Validate(), "Field 'pool' must point to one of the pools configured on a matching host 'localhost': [pool1|pool2] (actual: wrong)")
}

func TestConfig_MainPoolRef(t *testing.T) {
	cls := cluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id: sample_name1,
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

	cfg := config
	cfg.Cluster = cls

	assert.NoError(t, defaults.Assign(&cfg).Validate())
}

// When disk pool matches a host with no hostname, no extra error should
// be thrown beside missing host name error.
func TestConfig_MissingHostName(t *testing.T) {
	cls := cluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id: sample_name1,
					DataDisks: []DataDisk{
						{
							Name: "disk",
							Size: GB(5),
						},
					},
				},
			},
		},
	}

	cfg := config
	cfg.Cluster = cls
	cfg.Hosts = []Host{
		{
			Connection: Connection{
				Type: LOCAL,
			},
		},
	}

	assert.EqualError(t, defaults.Assign(&cfg).Validate(), "Field 'name' is required and cannot be empty.")
}
