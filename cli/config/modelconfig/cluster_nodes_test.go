package modelconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodes_Empty(t *testing.T) {
	assert.EqualError(t, Nodes{}.Validate(), "At least one master instance must be configured.")
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

	assert.EqualError(t, n.Validate(), "At least one master instance must be configured.")
}

func TestNodes_Minimal(t *testing.T) {
	id := "id"

	n := Nodes{
		LoadBalancer: LB{},
		Master: Master{
			Instances: []MasterInstance{
				{
					Id: &id,
				},
			},
		},
	}

	assert.NoError(t, n.Validate())
}

func TestNodes_MissingMaster(t *testing.T) {
	n := Nodes{
		LoadBalancer: LB{},
		Master:       Master{},
		Worker:       Worker{},
	}

	assert.EqualError(t, n.Validate(), "At least one master instance must be configured.")
}

func TestNodes_MissingLB(t *testing.T) {
	id1 := "id1"
	id2 := "id2"
	id3 := "id3"

	n := Nodes{
		Master: Master{
			Instances: []MasterInstance{
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

	assert.EqualError(t, n.Validate(), "At least one load balancer instance is required when multiple master instances are configured.")
}

func TestNodes_SingleLB(t *testing.T) {
	id1 := "id1"
	id2 := "id2"
	id3 := "id3"

	n := Nodes{
		LoadBalancer: LB{
			Instances: []LBInstance{
				{
					Id: &id1,
				},
			},
		},
		Master: Master{
			Instances: []MasterInstance{
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

	assert.NoError(t, n.Validate())
}

func TestNodes_MultiLB(t *testing.T) {
	id1 := "id1"
	id2 := "id2"
	id3 := "id3"
	ip := IPv4("192.168.113.13")

	n := Nodes{
		LoadBalancer: LB{
			VIP: &ip,
			Instances: []LBInstance{
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
		Master: Master{
			Instances: []MasterInstance{
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

	assert.NoError(t, n.Validate())
}
