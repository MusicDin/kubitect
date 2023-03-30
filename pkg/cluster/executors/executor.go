package executors

import "github.com/MusicDin/kubitect/pkg/cluster/event"

type Executor interface {
	Init() error
	Sync() error
	Create() error
	Upgrade() error
	ScaleUp(event.Events) error
	ScaleDown(event.Events) error
}
