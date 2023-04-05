package infra

import (
	"testing"

	"github.com/MusicDin/kubitect/pkg/models/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Empty(t *testing.T) {
	assert.NoError(t, Config{}.Validate(), "Terraform produced invalid output.")
}

func TestConfig(t *testing.T) {
	cfg := Config{}
	cfg.Nodes = config.Nodes{
		LoadBalancer: config.LB{
			VIP: config.IPv4("192.168.113.13"),
			Instances: []config.LBInstance{
				{Id: "id1"},
				{Id: "id2"},
				{Id: "id3"},
			},
		},
		Master: config.Master{
			Instances: []config.MasterInstance{
				{Id: "id1"},
				{Id: "id2"},
				{Id: "id3"},
			},
		},
		Worker: config.Worker{
			Instances: []config.WorkerInstance{
				{Id: "id1"},
				{Id: "id2"},
				{Id: "id3"},
			},
		},
	}

	assert.NoError(t, cfg.Validate())
}

func TestConfig_DuplicateIP(t *testing.T) {
	ip := config.IPv4("192.168.113.13")

	cfg := Config{}
	cfg.Nodes = config.Nodes{
		Master: config.Master{
			Instances: []config.MasterInstance{
				{
					Id: "id",
					IP: ip,
				},
			},
		},
		Worker: config.Worker{
			Instances: []config.WorkerInstance{
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
	mac := config.MAC("AA:BB:CC:DD:EE:FF")

	cfg := Config{}
	cfg.Nodes = config.Nodes{
		Master: config.Master{
			Instances: []config.MasterInstance{
				{
					Id:  "id",
					MAC: mac,
				},
			},
		},
		Worker: config.Worker{
			Instances: []config.WorkerInstance{
				{
					Id:  "id",
					MAC: mac,
				},
			},
		},
	}

	assert.EqualError(t, cfg.Validate(), "Duplicate MAC addresses detected in the provisioned infrastructure. (duplicates: [AA:BB:CC:DD:EE:FF])")
}
