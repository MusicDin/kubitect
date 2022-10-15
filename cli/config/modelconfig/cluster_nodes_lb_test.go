package modelconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLBPortForwardTarget(t *testing.T) {
	assert.Error(t, LBPortForwardTarget("wrong").Validate())
	assert.NoError(t, LBPortForwardTarget("all").Validate())
	assert.NoError(t, LBPortForwardTarget("workers").Validate())
	assert.NoError(t, LBPortForwardTarget("masters").Validate())
	assert.NoError(t, ALL.Validate())
}

func TestLBPortForward(t *testing.T) {
	name := "http"
	port := Port(80)

	fp := &LBPortForward{}
	fp1 := &LBPortForward{Port: &port}
	fp2 := &LBPortForward{Name: &name}
	fp3 := &LBPortForward{Name: &name, Port: &port}

	assert.ErrorContains(t, fp.Validate(), "Field 'name' is required.")
	assert.ErrorContains(t, fp.Validate(), "Field 'port' is required.")
	assert.ErrorContains(t, fp1.Validate(), "Field 'name' is required.")
	assert.ErrorContains(t, fp2.Validate(), "Field 'port' is required.")
	assert.NoError(t, fp3.Validate())
}

func TestLBDefault(t *testing.T) {
	size := GB(5)
	cpu := VCpu(5)

	def := LBDefault{
		CPU:          &cpu,
		RAM:          &size,
		MainDiskSize: &size,
	}

	assert.NoError(t, def.Validate())
	assert.NoError(t, LBDefault{}.Validate())
}

func TestLB_Minimal(t *testing.T) {
	id := "id1"

	lb := LB{
		Instances: &[]LBInstance{
			{
				Id: &id,
			},
		},
	}

	assert.NoError(t, lb.Validate())
	assert.NoError(t, LB{}.Validate())
}

func TestLB_MissingID(t *testing.T) {
	lb := LB{
		Instances: &[]LBInstance{{}},
	}

	assert.ErrorContains(t, lb.Validate(), "Field 'id' is required.")
}

func TestLB_UniqueID(t *testing.T) {
	id := "id1"

	lb := LB{
		Instances: &[]LBInstance{
			{
				Id: &id,
			},
			{
				Id: &id,
			},
		},
	}

	assert.ErrorContains(t, lb.Validate(), "Field 'Id' must be unique for each element in 'instances'.")
}

func TestLB_VIP(t *testing.T) {
	id1 := "id1"
	id2 := "id2"
	ip := IPv4("192.168.113.13")

	lb := LB{
		VIP: &ip,
		Instances: &[]LBInstance{
			{
				Id: &id1,
			},
			{
				Id: &id2,
			},
		},
	}

	assert.NoError(t, lb.Validate())
}

func TestLB_MissingVIP(t *testing.T) {
	id1 := "id1"
	id2 := "id2"

	lb := LB{
		Instances: &[]LBInstance{
			{
				Id: &id1,
			},
			{
				Id: &id2,
			},
		},
	}

	assert.ErrorContains(t, lb.Validate(), "Virtual IP (VIP) is required when multiple load balancer instances are configured.")
}
