package modelinfra

import (
	c "cli/config/modelconfig"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Empty(t *testing.T) {
	assert.NoError(t, Config{}.Validate(), "Terraform produced invalid output.")
}

func TestConfig(t *testing.T) {
	cfg := Config{}
	cfg.Nodes = c.Nodes{
		LoadBalancer: c.LB{
			VIP: c.IPv4("192.168.113.13"),
			Instances: []c.LBInstance{
				{Id: "id1"},
				{Id: "id2"},
				{Id: "id3"},
			},
		},
		Master: c.Master{
			Instances: []c.MasterInstance{
				{Id: "id1"},
				{Id: "id2"},
				{Id: "id3"},
			},
		},
		Worker: c.Worker{
			Instances: []c.WorkerInstance{
				{Id: "id1"},
				{Id: "id2"},
				{Id: "id3"},
			},
		},
	}

	assert.NoError(t, cfg.Validate())
}

func TestConfig_DuplicateIP(t *testing.T) {
	ip := c.IPv4("192.168.113.13")

	cfg := Config{}
	cfg.Nodes = c.Nodes{
		Master: c.Master{
			Instances: []c.MasterInstance{
				{
					Id: "id",
					IP: ip,
				},
			},
		},
		Worker: c.Worker{
			Instances: []c.WorkerInstance{
				{
					Id: "id",
					IP: ip,
				},
			},
		},
	}

	assert.EqualError(t, cfg.Validate(), "Duplicate IPs detected in the provisioned infrastructure. (duplicates: [192.168.113.13])")
}

func TestConfig_DuplicateMAC(t *testing.T) {
	mac := c.MAC("AA:BB:CC:DD:EE:FF")

	cfg := Config{}
	cfg.Nodes = c.Nodes{
		Master: c.Master{
			Instances: []c.MasterInstance{
				{
					Id:  "id",
					MAC: mac,
				},
			},
		},
		Worker: c.Worker{
			Instances: []c.WorkerInstance{
				{
					Id:  "id",
					MAC: mac,
				},
			},
		},
	}

	assert.EqualError(t, cfg.Validate(), "Duplicate MAC addresses detected in the provisioned infrastructure. (duplicates: [AA:BB:CC:DD:EE:FF])")
}
