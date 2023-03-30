package provisioner

import (
	"testing"

	"github.com/MusicDin/kubitect/pkg/cluster/event"
)

type provisionerMock struct{}

func (m provisionerMock) Init(event.Events) error { return nil }
func (m provisionerMock) Plan() (bool, error)     { return true, nil }
func (m provisionerMock) Apply() error            { return nil }
func (m provisionerMock) Destroy() error          { return nil }

func MockProvisioner(t *testing.T) Provisioner {
	return provisionerMock{}
}
