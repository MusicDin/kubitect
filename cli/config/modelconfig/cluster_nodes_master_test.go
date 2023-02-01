package modelconfig

import (
	"cli/utils/defaults"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMasterDefault(t *testing.T) {
	def := MasterDefault{
		CPU:          VCpu(5),
		RAM:          GB(5),
		MainDiskSize: GB(5),
	}

	assert.NoError(t, def.Validate())
	assert.ErrorContains(t, MasterDefault{}.Validate(), "Minimum value for field 'cpu' is 1 (actual: 0).")
	assert.ErrorContains(t, MasterDefault{}.Validate(), "Minimum value for field 'ram' is 1 (actual: 0).")
	assert.ErrorContains(t, MasterDefault{}.Validate(), "Minimum value for field 'mainDiskSize' is 1 (actual: 0).")
}

func TestMaster_Type(t *testing.T) {
	assert.Equal(t, MasterInstance{}.GetTypeName(), "master")
}

func TestMaster_Minimal(t *testing.T) {
	m := Master{
		Instances: []MasterInstance{
			{Id: "id"},
		},
	}

	assert.NoError(t, defaults.Assign(&m).Validate())
	assert.EqualError(t, defaults.Assign(&Master{}).Validate(), "At least one master instance must be configured.")
}

func TestMaster_MissingID(t *testing.T) {
	m := Master{
		Instances: []MasterInstance{{}},
	}

	assert.Empty(t, m.Instances[0].GetID())
	assert.EqualError(t, defaults.Assign(&m).Validate(), "Field 'id' is required and cannot be empty.")
}

func TestMaster_UniqueID(t *testing.T) {
	m := Master{
		Instances: []MasterInstance{
			{Id: "id"},
			{Id: "id"},
			{Id: "id"},
		},
	}

	assert.EqualError(t, defaults.Assign(&m).Validate(), "Field 'Id' must be unique for each element in 'instances'.")
}

func TestMaster_OddNumberOfInstances(t *testing.T) {
	m := Master{
		Instances: []MasterInstance{
			{Id: "id1"},
			{Id: "id2"},
		},
	}

	assert.EqualError(t, defaults.Assign(&m).Validate(), "Number of master instances must be odd (1, 3, 5 etc.).")
}

func TestMaster_DataDisk(t *testing.T) {
	m := Master{
		Instances: []MasterInstance{
			{
				Id: "id",
				DataDisks: []DataDisk{
					{
						Name: "disk",
						Size: GB(42),
					},
				},
			},
		},
	}

	assert.NoError(t, defaults.Assign(&m).Validate())
}

func TestMaster_DataDiskUniqueName(t *testing.T) {
	m := Master{
		Instances: []MasterInstance{
			{
				Id: "id",
				DataDisks: []DataDisk{
					{Name: "disk1", Size: GB(42)},
					{Name: "disk1", Size: GB(42)},
				},
			},
		},
	}

	assert.EqualError(t, defaults.Assign(&m).Validate(), "Field 'Name' must be unique for each element in 'dataDisks'.")
}
