package modelconfig

import (
	"testing"

	"github.com/MusicDin/kubitect/cli/utils/defaults"

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

func TestWorker_Defaults(t *testing.T) {
	w := Worker{
		Default: WorkerDefault{
			CPU:          VCpu(2),
			RAM:          GB(4),
			MainDiskSize: GB(256),
		},
		Instances: []WorkerInstance{
			{Id: "id1", CPU: VCpu(4)},
			{Id: "id2", RAM: GB(8)},
			{Id: "id3", MainDiskSize: (512)},
		},
	}

	defaults.Assign(&w)
	assert.Equal(t, VCpu(4), w.Instances[0].CPU)
	assert.Equal(t, VCpu(2), w.Instances[1].CPU)
	assert.Equal(t, VCpu(2), w.Instances[2].CPU)
	assert.Equal(t, GB(4), w.Instances[0].RAM)
	assert.Equal(t, GB(8), w.Instances[1].RAM)
	assert.Equal(t, GB(4), w.Instances[2].RAM)
	assert.Equal(t, GB(256), w.Instances[0].MainDiskSize)
	assert.Equal(t, GB(256), w.Instances[1].MainDiskSize)
	assert.Equal(t, GB(512), w.Instances[2].MainDiskSize)
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

func TestWorker_DefaultDataDisks(t *testing.T) {
	defDisks := []DataDisk{
		{Name: "def-disk1", Size: GB(42)},
		{Name: "def-disk2", Size: GB(42)},
	}

	addDisks := []DataDisk{
		{Name: "disk1", Size: GB(42)},
		{Name: "disk2", Size: GB(42)},
	}

	w := Worker{
		Default: WorkerDefault{
			DataDisks: defDisks,
		},
		Instances: []WorkerInstance{
			{Id: "id1", DataDisks: addDisks},
			{Id: "id2"},
			{Id: "id3"},
		},
	}

	assert.NoError(t, defaults.Assign(&w).Validate())
	assert.Equal(t, append(defDisks, addDisks...), w.Instances[0].DataDisks)
	assert.Equal(t, w.Default.DataDisks, w.Instances[1].DataDisks)
	assert.Equal(t, w.Default.DataDisks, w.Instances[2].DataDisks)
}
