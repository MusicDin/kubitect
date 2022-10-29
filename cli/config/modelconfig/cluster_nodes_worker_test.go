package modelconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkerDefault(t *testing.T) {
	size := GB(42)
	cpu := VCpu(24)

	def := WorkerDefault{
		CPU:          &cpu,
		RAM:          &size,
		MainDiskSize: &size,
	}

	assert.NoError(t, def.Validate())
	assert.NoError(t, WorkerDefault{}.Validate())
}

func TestWorker_Minimal(t *testing.T) {
	id := "id"

	w := Worker{
		Instances: []WorkerInstance{
			{
				Id: &id,
			},
		},
	}

	assert.NoError(t, Worker{}.Validate())
	assert.NoError(t, w.Validate())
}

func TestWorker_MissingID(t *testing.T) {
	w := Worker{
		Instances: []WorkerInstance{{}},
	}

	assert.EqualError(t, w.Validate(), "Field 'id' is required.")
}

func TestWorker_UniqueID(t *testing.T) {
	id := "id"

	w := Worker{
		Instances: []WorkerInstance{
			{
				Id: &id,
			},
			{
				Id: &id,
			},
		},
	}

	assert.EqualError(t, w.Validate(), "Field 'Id' must be unique for each element in 'instances'.")
}

func TestWorker_DataDisk(t *testing.T) {
	name := "id"
	size := GB(42)

	w := Worker{
		Instances: []WorkerInstance{
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

	assert.NoError(t, w.Validate())
}

func TestWorker_DataDiskUniqueName(t *testing.T) {
	name := "id"
	size := GB(42)

	w := Worker{
		Instances: []WorkerInstance{
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

	assert.EqualError(t, w.Validate(), "Field 'Name' must be unique for each element in 'dataDisks'.")
}
