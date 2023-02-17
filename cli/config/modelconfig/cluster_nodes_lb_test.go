package modelconfig

import (
	"testing"

	"github.com/MusicDin/kubitect/cli/utils/defaults"

	"github.com/stretchr/testify/assert"
)

func TestLBPortForwardTarget(t *testing.T) {
	assert.NoError(t, ALL.Validate())
	assert.NoError(t, LBPortForwardTarget("all").Validate())
	assert.NoError(t, LBPortForwardTarget("workers").Validate())
	assert.NoError(t, LBPortForwardTarget("masters").Validate())
	assert.EqualError(t, LBPortForwardTarget("wrong").Validate(), "Field must be one of the following values: [workers|masters|all] (actual: wrong).")
}

func TestLBPortForward(t *testing.T) {
	fp := LBPortForward{}
	fp1 := LBPortForward{Port: 80}
	fp2 := LBPortForward{Name: "http"}
	fp3 := LBPortForward{Name: "http", Port: 80}

	assert.ErrorContains(t, fp.Validate(), "Field 'name' is required and cannot be empty.")
	assert.ErrorContains(t, fp.Validate(), "Field 'port' is required and cannot be empty.")
	assert.ErrorContains(t, fp1.Validate(), "Field 'name' is required and cannot be empty.")
	assert.ErrorContains(t, fp2.Validate(), "Field 'port' is required and cannot be empty.")
	assert.EqualError(t, fp3.Validate(), "Minimum value for field 'targetPort' is 1 (actual: 0).")
}

func TestLBPortForward_Default(t *testing.T) {
	fp := defaults.Assign(&LBPortForward{Name: "http", Port: 80})

	assert.NoError(t, fp.Validate())
	assert.Equal(t, Port(80), fp.TargetPort)
	assert.Equal(t, WORKERS, fp.Target)
}

func TestLBDefault(t *testing.T) {
	def := LBDefault{
		CPU:          VCpu(5),
		RAM:          GB(5),
		MainDiskSize: GB(5),
	}

	assert.NoError(t, def.Validate())
	assert.NoError(t, def.Validate())
	assert.NoError(t, defaults.Assign(&LBDefault{}).Validate())
	assert.ErrorContains(t, LBDefault{}.Validate(), "Minimum value for field 'cpu' is 1 (actual: 0).")
	assert.ErrorContains(t, LBDefault{}.Validate(), "Minimum value for field 'ram' is 1 (actual: 0).")
	assert.ErrorContains(t, LBDefault{}.Validate(), "Minimum value for field 'mainDiskSize' is 1 (actual: 0).")
}

func TestLB_Type(t *testing.T) {
	assert.Equal(t, LBInstance{}.GetTypeName(), "lb")
}

func TestLB_Minimal(t *testing.T) {
	lb := LB{
		Instances: []LBInstance{
			{Id: "id"},
		},
	}

	assert.NoError(t, defaults.Assign(&LBDefault{}).Validate())
	assert.NoError(t, defaults.Assign(&lb).Validate())
}

func TestLB_Defaults(t *testing.T) {
	lb1 := LB{
		Instances: []LBInstance{
			{Id: "id"},
		},
	}

	lb2 := LB{
		Instances: []LBInstance{
			{Id: "id1"},
			{Id: "id2"},
		},
	}

	defaults.Assign(&lb1)
	defaults.Assign(&lb2)

	assert.Nil(t, lb1.VirtualRouterId, "LB VRID is set even if only one instance is configured!")
	assert.Nil(t, lb1.Instances[0].Priority, "LB instance priority is set even if only one instance is configured!")
	assert.Equal(t, &defaultVRID, lb2.VirtualRouterId, "Default LB VRID is not set when multiple instances are configured!")
	assert.Equal(t, &defaultPriority, lb2.Instances[0].Priority, "Default LB instance priority is not set when multiple instances are configured!")
}

func TestLB_MissingID(t *testing.T) {
	lb := LB{
		Instances: []LBInstance{{}},
	}

	assert.Empty(t, lb.Instances[0].GetID())
	assert.EqualError(t, defaults.Assign(&lb).Validate(), "Field 'id' is required and cannot be empty.")
}

func TestLB_UniqueID(t *testing.T) {
	lb := LB{
		VIP: IPv4("192.168.113.13"),
		Instances: []LBInstance{
			{Id: "id"},
			{Id: "id"},
		},
	}

	assert.EqualError(t, defaults.Assign(&lb).Validate(), "Field 'Id' must be unique for each element in 'instances'.")
}

func TestLB_VIP(t *testing.T) {
	lb := LB{
		VIP: IPv4("192.168.113.13"),
		Instances: []LBInstance{
			{Id: "id1"},
			{Id: "id2"},
		},
	}

	assert.NoError(t, defaults.Assign(&lb).Validate())
}

func TestLB_MissingVIP(t *testing.T) {
	lb := LB{
		Instances: []LBInstance{
			{Id: "id1"},
			{Id: "id2"},
		},
	}

	assert.EqualError(t, defaults.Assign(&lb).Validate(), "Virtual IP (VIP) is required when multiple load balancer instances are configured.")
}
