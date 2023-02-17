package executors

import "github.com/MusicDin/kubitect/cli/cluster/event"

type Executor interface {
	Init() error
	Create() error
	Upgrade() error
	ScaleUp(event.Events) error
	ScaleDown(event.Events) error
}
