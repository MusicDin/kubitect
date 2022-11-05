package modelconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sample_default = true
	sample_name1   = "test"
	sample_name2   = "test2"
	sample_dd_size = GB(5)
	sample_ip      = IPv4("192.168.114.13")

	net_cidr = CIDRv4("192.168.113.0/24")
	net      = Network{
		CIDR: &net_cidr,
	}

	localhost_type = LOCAL
	localhost      = Host{
		Name:    &sample_name1,
		Default: &sample_default,
		Connection: Connection{
			Type: &localhost_type,
		},
		DataResourcePools: []DataResourcePool{
			{
				Name: &sample_name1,
			},
			{
				Name: &sample_name2,
			},
		},
	}

	remotehost_keyfile = File("./config_test.go")
	remotehost_user    = User("user")
	remotehost_type    = REMOTE
	remotehost         = Host{
		Name: &sample_name2,
		Connection: Connection{
			Type: &remotehost_type,
			IP:   &sample_ip,
			User: &remotehost_user,
			SSH: ConnectionSSH{
				Keyfile: &remotehost_keyfile,
			},
		},
	}

	k8s_version    = Version("v1.2.3")
	k8s_ks_version = MasterVersion(k8s_version)
	k8s            = Kubernetes{
		Version: &k8s_version,
		Kubespray: Kubespray{
			Version: &k8s_ks_version,
		},
	}

	cluster = Cluster{
		Name:    &sample_name1,
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
	assert.ErrorContains(t, Config{}.Validate(), "At least 1 hosts must be configured.")
	assert.ErrorContains(t, Config{}.Validate(), "Configuration must contain 'cluster' section.")
	assert.ErrorContains(t, Config{}.Validate(), "Configuration must contain 'kubernetes' section.")
}

func TestConfig_Valid(t *testing.T) {
	cls := cluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id:   &sample_name1,
					Host: &sample_name1,
					DataDisks: []DataDisk{
						// Correct pool reference
						{
							Name: &sample_name1,
							Pool: &sample_name1,
							Size: &sample_dd_size,
						},
						// Main pool reference (no direct pool ref)
						{
							Name: &sample_name2,
							Size: &sample_dd_size,
						},
					},
				},
			},
		},
	}

	cfg := config
	cfg.Cluster = cls

	assert.NoError(t, cfg.Validate())
}

func TestConfig_InvalidIP(t *testing.T) {
	cls := cluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id:   &sample_name1,
					Host: &sample_name1,
					IP:   &sample_ip,
				},
			},
		},
	}

	cfg := config
	cfg.Cluster = cls

	assert.EqualError(t, cfg.Validate(), "Field 'ip' must be a valid IP address within '192.168.113.0/24' subnet. (actual: 192.168.114.13)")
}

func TestConfig_MultipleDefaultHosts(t *testing.T) {
	cls := cluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id:   &sample_name1,
					Host: &sample_name1,
				},
			},
		},
	}

	rhDef := remotehost
	rhDef.Default = &sample_default

	cfg := config
	cfg.Cluster = cls
	cfg.Hosts = []Host{
		localhost,
		rhDef,
	}

	assert.EqualError(t, cfg.Validate(), "Only one host can be configured as default.")
}

func TestConfig_InvalidHostRef(t *testing.T) {
	host := "wrong"

	cls := cluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id:   &sample_name1,
					Host: &host,
				},
			},
		},
	}

	cfg := config
	cfg.Cluster = cls

	assert.EqualError(t, cfg.Validate(), "Field 'host' must point to one of the configured hosts: [test] (actual: wrong)")
}

func TestConfig_InvalidPoolHostRef(t *testing.T) {
	cls := cluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id:   &sample_name1,
					Host: &sample_name2,
					DataDisks: []DataDisk{
						{
							Name: &sample_name1,
							Pool: &sample_name1,
							Size: &sample_dd_size,
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

	assert.EqualError(t, cfg.Validate(), "Field 'pool' points to a data resource pool, but matching host 'test2' has none configured.")
}

func TestConfig_InvalidPoolRef(t *testing.T) {
	pool := "wrong"

	cls := cluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id: &sample_name1,
					DataDisks: []DataDisk{
						{
							Name: &sample_name1,
							Pool: &pool,
							Size: &sample_dd_size,
						},
					},
				},
			},
		},
	}

	cfg := config
	cfg.Cluster = cls

	assert.EqualError(t, cfg.Validate(), "Field 'pool' must point to one of the pools configured on a matching host 'test': [test|test2] (actual: wrong)")
}

func TestConfig_MainPoolRef(t *testing.T) {
	pool := "main"

	cls := cluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id: &sample_name1,
					DataDisks: []DataDisk{
						{
							Name: &sample_name1,
							Pool: &pool,
							Size: &sample_dd_size,
						},
					},
				},
			},
		},
	}

	cfg := config
	cfg.Cluster = cls

	assert.NoError(t, cfg.Validate())
}

// When disk pool matches a host with no hostname, no extra error should
// be thrown beside missing host name error.
func TestConfig_MissingHostName(t *testing.T) {
	cls := cluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id: &sample_name1,
					DataDisks: []DataDisk{
						{
							Name: &sample_name1,
							Size: &sample_dd_size,
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
				Type: &localhost_type,
			},
		},
	}

	assert.EqualError(t, cfg.Validate(), "Field 'name' is required.")
}
