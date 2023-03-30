package provisioner

import "github.com/MusicDin/kubitect/pkg/cluster/event"

type Provisioner interface {
	Init(events event.Events) error
	Plan() (bool, error)
	Apply() error
	Destroy() error
}
