package event

import (
	"testing"

	"github.com/MusicDin/kubitect/pkg/utils/cmp"
)

func MockEvent(t *testing.T, eventType EventType, action cmp.ActionType, changes []cmp.Change) Event {
	t.Helper()

	return Event{
		eType:   eventType,
		changes: changes,
		path:    t.TempDir(),
		msg:     "mock event",
		action:  action,
	}
}
