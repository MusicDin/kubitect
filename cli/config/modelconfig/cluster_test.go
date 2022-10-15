package modelconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	cls_sample_name    = "test"
	cls_sample_cidr    = CIDRv4("192.168.113.13/24")
	cls_sample_cluster = &Cluster{
		Name: &cls_sample_name,
		Network: &Network{
			CIDR: &cls_sample_cidr,
		},
	}
)

func TestCluster_Empty(t *testing.T) {
	assert.ErrorContains(t, Cluster{}.Validate(), "Field 'name' is required.")
	assert.ErrorContains(t, Cluster{}.Validate(), "Field 'network' is required.")
	assert.ErrorContains(t, Cluster{}.Validate(), "Field 'nodes' is required.")
}

func TestCluster_Minimal(t *testing.T) {
	cls := cls_sample_cluster
	cls.Nodes = &Nodes{
		Master: &Master{
			Instances: &[]MasterInstance{
				{
					Id: &cls_sample_name,
				},
			},
		},
	}

	assert.NoError(t, cls.Validate())
}

func TestCluster_DuplicateIP(t *testing.T) {
	ip := IPv4("192.168.113.13")

	cls := cls_sample_cluster
	cls.Nodes = &Nodes{
		Master: &Master{
			Instances: &[]MasterInstance{
				{
					Id: &cls_sample_name,
					IP: &ip,
				},
			},
		},
		Worker: &Worker{
			Instances: &[]WorkerInstance{
				{
					Id: &cls_sample_name,
					IP: &ip,
				},
			},
		},
	}

	assert.ErrorContains(t, cls.Validate(), "IP address of each node instance (including VIP) must be unique. (duplicates: [192.168.113.13])")
}

func TestCluster_DuplicateMAC(t *testing.T) {
	mac := MAC("AA:BB:CC:DD:EE:FF")

	cls := cls_sample_cluster
	cls.Nodes = &Nodes{
		Master: &Master{
			Instances: &[]MasterInstance{
				{
					Id:  &cls_sample_name,
					MAC: &mac,
				},
			},
		},
		Worker: &Worker{
			Instances: &[]WorkerInstance{
				{
					Id:  &cls_sample_name,
					MAC: &mac,
				},
			},
		},
	}

	assert.ErrorContains(t, cls.Validate(), "MAC address of each node instance must be unique. (duplicates: [AA:BB:CC:DD:EE:FF])")
}

func TestCluster_Complete(t *testing.T) {
	id1 := "id1"
	id2 := "id2"
	id3 := "id3"
	ip := IPv4("192.168.113.13")

	cls := cls_sample_cluster
	cls.Nodes = &Nodes{
		LoadBalancer: &LB{
			VIP: &ip,
			Instances: &[]LBInstance{
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
		Master: &Master{
			Instances: &[]MasterInstance{
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
		Worker: &Worker{
			Instances: &[]WorkerInstance{
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
