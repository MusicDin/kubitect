package modelinfra

import (
	"github.com/MusicDin/kubitect/cli/pkg/config/modelconfig"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Empty(t *testing.T) {
	assert.NoError(t, Config{}.Validate(), "Terraform produced invalid output.")
}

func TestConfig(t *testing.T) {
	cfg := Config{}
	cfg.Nodes = modelconfig.Nodes{
		LoadBalancer: modelconfig.LB{
			VIP: modelconfig.IPv4("192.168.113.13"),
			Instances: []modelconfig.LBInstance{
				{Id: "id1"},
				{Id: "id2"},
				{Id: "id3"},
			},
		},
		Master: modelconfig.Master{
			Instances: []modelconfig.MasterInstance{
				{Id: "id1"},
				{Id: "id2"},
				{Id: "id3"},
			},
		},
		Worker: modelconfig.Worker{
			Instances: []modelconfig.WorkerInstance{
				{Id: "id1"},
				{Id: "id2"},
				{Id: "id3"},
			},
		},
	}

	assert.NoError(t, cfg.Validate())
}

func TestConfig_DuplicateIP(t *testing.T) {
	ip := modelconfig.IPv4("192.168.113.13")

	cfg := Config{}
	cfg.Nodes = modelconfig.Nodes{
		Master: modelconfig.Master{
			Instances: []modelconfig.MasterInstance{
				{
					Id: "id",
					IP: ip,
				},
			},
		},
		Worker: modelconfig.Worker{
			Instances: []modelconfig.WorkerInstance{
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
	mac := modelconfig.MAC("AA:BB:CC:DD:EE:FF")

	cfg := Config{}
	cfg.Nodes = modelconfig.Nodes{
		Master: modelconfig.Master{
			Instances: []modelconfig.MasterInstance{
				{
					Id:  "id",
					MAC: mac,
				},
			},
		},
		Worker: modelconfig.Worker{
			Instances: []modelconfig.WorkerInstance{
				{
					Id:  "id",
					MAC: mac,
				},
			},
		},
	}

	assert.EqualError(t, cfg.Validate(), "Duplicate MAC addresses detected in the provisioned infrastructure. (duplicates: [AA:BB:CC:DD:EE:FF])")
}
