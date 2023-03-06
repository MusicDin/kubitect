package modelconfig

import (
	"github.com/MusicDin/kubitect/cli/pkg/utils/defaults"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodes_Empty(t *testing.T) {
	assert.ErrorContains(t, Nodes{}.Validate(), "At least one master instance must be configured.")
	assert.ErrorContains(t, defaults.Assign(&Nodes{}).Validate(), "At least one master instance must be configured.")
}

func TestNodes_Mock(t *testing.T) {
	assert.NoError(t, MockNodes(t).Validate())
}

func TestCluster_NilInstances(t *testing.T) {
	n := Nodes{
		LoadBalancer: LB{
			Instances: nil,
		},
		Master: Master{
			Instances: nil,
		},
		Worker: Worker{
			Instances: nil,
		},
	}

	assert.EqualError(t, defaults.Assign(&n).Validate(), "At least one master instance must be configured.")
}

func TestNodes_Minimal(t *testing.T) {
	n := Nodes{
		LoadBalancer: LB{},
		Master: Master{
			Instances: []MasterInstance{
				{
					Id: "id",
				},
			},
		},
	}

	assert.NoError(t, defaults.Assign(&n).Validate())
}

func TestNodes_MissingMaster(t *testing.T) {
	n := Nodes{
		LoadBalancer: LB{},
		Master:       Master{},
		Worker:       Worker{},
	}

	assert.EqualError(t, defaults.Assign(&n).Validate(), "At least one master instance must be configured.")
}

func TestNodes_MissingLB(t *testing.T) {
	n := Nodes{
		Master: Master{
			Instances: []MasterInstance{
				{Id: "id1"},
				{Id: "id2"},
				{Id: "id3"},
			},
		},
	}

	assert.EqualError(t, defaults.Assign(&n).Validate(), "At least one load balancer instance is required when multiple master instances are configured.")
}

func TestNodes_SingleLB(t *testing.T) {
	n := Nodes{
		LoadBalancer: LB{
			Instances: []LBInstance{
				{Id: "lb"},
			},
		},
		Master: Master{
			Instances: []MasterInstance{
				{Id: "m1"},
				{Id: "m2"},
				{Id: "m3"},
			},
		},
	}

	assert.NoError(t, defaults.Assign(&n).Validate())
}

func TestNodes_MultiLB(t *testing.T) {
	n := Nodes{
		LoadBalancer: LB{
			VIP: IPv4("192.168.113.13"),
			Instances: []LBInstance{
				{Id: "lb1"},
				{Id: "lb2"},
				{Id: "lb3"},
			},
		},
		Master: Master{
			Instances: []MasterInstance{
				{Id: "m1"},
				{Id: "m2"},
				{Id: "m3"},
			},
		},
	}

	assert.NoError(t, defaults.Assign(&n).Validate())
}
