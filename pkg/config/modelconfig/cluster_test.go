package modelconfig

import (
	"github.com/MusicDin/kubitect/pkg/utils/defaults"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sampleCluster = Cluster{
		Name: "test",
		Network: Network{
			CIDR: CIDRv4("192.168.113.13/24"),
		},
	}
)

func TestCluster_Empty(t *testing.T) {
	cls := defaults.Assign(&Cluster{})
	assert.ErrorContains(t, cls.Validate(), "Field 'name' is required and cannot be empty")
	assert.ErrorContains(t, cls.Validate(), "Field 'cidr' is required and cannot be empty.")
	assert.ErrorContains(t, cls.Validate(), "At least one master instance must be configured.")
}

func TestCluster_Minimal(t *testing.T) {
	cls := sampleCluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{Id: "id"},
			},
		},
	}

	assert.NoError(t, defaults.Assign(&cls).Validate())
}

func TestCluster_DuplicateIP(t *testing.T) {
	ip := IPv4("192.168.113.13")

	cls := sampleCluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id: "id",
					IP: ip,
				},
			},
		},
		Worker: Worker{
			Instances: []WorkerInstance{
				{
					Id: "id",
					IP: ip,
				},
			},
		},
	}

	assert.EqualError(t, defaults.Assign(&cls).Validate(), "IP address of each node instance (including VIP) must be unique. (duplicates: [192.168.113.13])")
}

func TestCluster_DuplicateMAC(t *testing.T) {
	mac := MAC("AA:BB:CC:DD:EE:FF")

	cls := sampleCluster
	cls.Nodes = Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{
					Id:  "id",
					MAC: mac,
				},
			},
		},
		Worker: Worker{
			Instances: []WorkerInstance{
				{
					Id:  "id",
					MAC: mac,
				},
			},
		},
	}

	assert.EqualError(t, defaults.Assign(&cls).Validate(), "MAC address of each node instance must be unique. (duplicates: [AA:BB:CC:DD:EE:FF])")
}

func TestCluster_Complete(t *testing.T) {
	cls := sampleCluster
	cls.Nodes = Nodes{
		LoadBalancer: LB{
			VIP: IPv4("192.168.113.13"),
			Instances: []LBInstance{
				{Id: "id1"},
				{Id: "id2"},
				{Id: "id3"},
			},
		},
		Master: Master{
			Instances: []MasterInstance{
				{Id: "id1"},
				{Id: "id2"},
				{Id: "id3"},
			},
		},
		Worker: Worker{
			Instances: []WorkerInstance{
				{Id: "id1"},
				{Id: "id2"},
				{Id: "id3"},
			},
		},
	}

	assert.NoError(t, defaults.Assign(&cls).Validate())
}
