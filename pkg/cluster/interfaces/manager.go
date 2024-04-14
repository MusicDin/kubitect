package interfaces

import "github.com/MusicDin/kubitect/pkg/cluster/event"

type Manager interface {
	Init() error
	Sync() error
	Create() error
	Upgrade() error
	ScaleUp(event.Events) error
	ScaleDown(event.Events) error
}
