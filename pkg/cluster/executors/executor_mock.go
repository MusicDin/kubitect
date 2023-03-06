package executors

import (
	"testing"

	"github.com/MusicDin/kubitect/pkg/cluster/event"
)

type executorMock struct{}

func (m executorMock) Init() error                  { return nil }
func (m executorMock) Create() error                { return nil }
func (m executorMock) Upgrade() error               { return nil }
func (m executorMock) ScaleDown(event.Events) error { return nil }
func (m executorMock) ScaleUp(event.Events) error   { return nil }

func MockExecutor(t *testing.T) Executor {
	return executorMock{}
}
