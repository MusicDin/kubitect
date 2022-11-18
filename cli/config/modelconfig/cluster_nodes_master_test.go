package modelconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMasterDefault(t *testing.T) {
	size := GB(5)
	cpu := VCpu(5)

	def := MasterDefault{
		CPU:          &cpu,
		RAM:          &size,
		MainDiskSize: &size,
	}

	assert.NoError(t, def.Validate())
	assert.NoError(t, MasterDefault{}.Validate())
}

func TestMaster_Type(t *testing.T) {
	assert.Equal(t, MasterInstance{}.GetTypeName(), "master")
}

func TestMaster_Minimal(t *testing.T) {
	id := "id"

	m := Master{
		Instances: []MasterInstance{
			{
				Id: &id,
			},
		},
	}

	assert.NoError(t, m.Validate())
	assert.EqualError(t, Master{}.Validate(), "At least one master instance must be configured.")
}

func TestMaster_MissingID(t *testing.T) {
	m := Master{
		Instances: []MasterInstance{{}},
	}

	assert.Nil(t, m.Instances[0].GetID())
	assert.EqualError(t, m.Validate(), "Field 'id' is required.")
}

func TestMaster_UniqueID(t *testing.T) {
	id := "id"

	m := Master{
		Instances: []MasterInstance{
			{
				Id: &id,
			},
			{
				Id: &id,
			},
			{
				Id: &id,
			},
		},
	}

	assert.EqualError(t, m.Validate(), "Field 'Id' must be unique for each element in 'instances'.")
}

func TestMaster_OddNumberOfInstances(t *testing.T) {
	id1 := "id"
	id2 := "id2"

	m := Master{
		Instances: []MasterInstance{
			{
				Id: &id1,
			},
			{
				Id: &id2,
			},
		},
	}

	assert.EqualError(t, m.Validate(), "Number of master instances must be odd (1, 3, 5 etc.).")
}

func TestMaster_DataDisk(t *testing.T) {
	name := "id"
	size := GB(42)

	m := Master{
		Instances: []MasterInstance{
			{
				Id: &name,
				DataDisks: []DataDisk{
					{
						Name: &name,
						Size: &size,
					},
				},
			},
		},
	}

	assert.NoError(t, m.Validate())
}

func TestMaster_DataDiskUniqueName(t *testing.T) {
	name := "id"
	size := GB(42)

	m := Master{
		Instances: []MasterInstance{
			{
				Id: &name,
				DataDisks: []DataDisk{
					{
						Name: &name,
						Size: &size,
					},
					{
						Name: &name,
						Size: &size,
					},
				},
			},
		},
	}

	assert.EqualError(t, m.Validate(), "Field 'Name' must be unique for each element in 'dataDisks'.")
}
