package modelinfra

import (
	c "cli/config/modelconfig"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Empty(t *testing.T) {
	assert.NoError(t, Config{}.Validate())
}

func TestConfig_DuplicateIP(t *testing.T) {
	name := "test"
	ip := c.IPv4("192.168.113.13")

	cls := Config{}
	cls.Nodes = c.Nodes{
		Master: c.Master{
			Instances: []c.MasterInstance{
				{
					Id: &name,
					IP: &ip,
				},
			},
		},
		Worker: c.Worker{
			Instances: []c.WorkerInstance{
				{
					Id: &name,
					IP: &ip,
				},
			},
		},
	}

	assert.EqualError(t, cls.Validate(), "Duplicate IPs detected in the provisioned infrastructure. (duplicates: [192.168.113.13])")
}

func TestConfig_DuplicateMAC(t *testing.T) {
	id := "test"
	mac := c.MAC("AA:BB:CC:DD:EE:FF")

	cls := Config{}
	cls.Nodes = c.Nodes{
		Master: c.Master{
			Instances: []c.MasterInstance{
				{
					Id:  &id,
					MAC: &mac,
				},
			},
		},
		Worker: c.Worker{
			Instances: []c.WorkerInstance{
				{
					Id:  &id,
					MAC: &mac,
				},
			},
		},
	}

	assert.EqualError(t, cls.Validate(), "Duplicate MAC addresses detected in the provisioned infrastructure. (duplicates: [AA:BB:CC:DD:EE:FF])")
}

func TestConfig_Complete(t *testing.T) {
	id1 := "id1"
	id2 := "id2"
	id3 := "id3"
	ip := c.IPv4("192.168.113.13")

	cls := Config{}
	cls.Nodes = c.Nodes{
		LoadBalancer: c.LB{
			VIP: &ip,
			Instances: []c.LBInstance{
				{
					Id: &id1,
				},
				{
					Id: &id2,
				},
				{
					Id: &id3,
				},
			},
		},
		Master: c.Master{
			Instances: []c.MasterInstance{
				{
					Id: &id1,
				},
				{
					Id: &id2,
				},
				{
					Id: &id3,
				},
			},
		},
		Worker: c.Worker{
			Instances: []c.WorkerInstance{
				{
					Id: &id1,
				},
				{
					Id: &id2,
				},
				{
					Id: &id3,
				},
			},
		},
	}

	assert.NoError(t, cls.Validate())
}
