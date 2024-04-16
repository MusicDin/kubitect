package interfaces

import (
	"testing"

	"github.com/MusicDin/kubitect/pkg/cluster/event"
)

type managerMock struct{}

func (m managerMock) Init() error                  { return nil }
func (m managerMock) Sync() error                  { return nil }
func (m managerMock) Create() error                { return nil }
func (m managerMock) Upgrade() error               { return nil }
func (m managerMock) ScaleDown(event.Events) error { return nil }
func (m managerMock) ScaleUp(event.Events) error   { return nil }

func MockManager(t *testing.T) Manager {
	return managerMock{}
}
