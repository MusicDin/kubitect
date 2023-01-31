package provisioner

import "testing"

type provisionerMock struct{}

func (m provisionerMock) Init() error         { return nil }
func (m provisionerMock) Plan() (bool, error) { return true, nil }
func (m provisionerMock) Apply() error        { return nil }
func (m provisionerMock) Destroy() error      { return nil }

func MockProvisioner(t *testing.T) Provisioner {
	return provisionerMock{}
}
