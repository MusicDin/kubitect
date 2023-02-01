package modelconfig

import (
	"cli/utils/defaults"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkerDefault(t *testing.T) {
	def := WorkerDefault{
		CPU:          VCpu(5),
		RAM:          GB(5),
		MainDiskSize: GB(5),
	}

	assert.NoError(t, def.Validate())
	assert.NoError(t, defaults.Assign(&def).Validate())
	assert.ErrorContains(t, WorkerDefault{}.Validate(), "Minimum value for field 'cpu' is 1 (actual: 0).")
	assert.ErrorContains(t, WorkerDefault{}.Validate(), "Minimum value for field 'ram' is 1 (actual: 0).")
	assert.ErrorContains(t, WorkerDefault{}.Validate(), "Minimum value for field 'mainDiskSize' is 1 (actual: 0).")
}

func TestWorker_Type(t *testing.T) {
	assert.Equal(t, WorkerInstance{}.GetTypeName(), "worker")
}

func TestWorker_Minimal(t *testing.T) {
	w := Worker{
		Instances: []WorkerInstance{
			{Id: "id"},
		},
	}

	assert.NoError(t, defaults.Assign(&w).Validate())
	assert.NoError(t, defaults.Assign(&Worker{}).Validate())
}

func TestWorker_MissingID(t *testing.T) {
	w := Worker{
		Instances: []WorkerInstance{{}},
	}

	assert.Empty(t, w.Instances[0].GetID())
	assert.EqualError(t, defaults.Assign(&w).Validate(), "Field 'id' is required and cannot be empty.")
}

func TestWorker_UniqueID(t *testing.T) {
	w := Worker{
		Instances: []WorkerInstance{
			{Id: "id"},
			{Id: "id"},
		},
	}

	assert.EqualError(t, defaults.Assign(&w).Validate(), "Field 'Id' must be unique for each element in 'instances'.")
}

func TestWorker_DataDisk(t *testing.T) {
	w := Worker{
		Instances: []WorkerInstance{
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

	assert.NoError(t, defaults.Assign(&w).Validate())
}

func TestWorker_DataDiskUniqueName(t *testing.T) {
	w := Worker{
		Instances: []WorkerInstance{
			{
				Id: "id",
				DataDisks: []DataDisk{
					{Name: "disk1", Size: GB(42)},
					{Name: "disk1", Size: GB(42)},
				},
			},
		},
	}

	assert.EqualError(t, defaults.Assign(&w).Validate(), "Field 'Name' must be unique for each element in 'dataDisks'.")
}
