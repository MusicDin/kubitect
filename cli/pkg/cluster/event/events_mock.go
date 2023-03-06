package event

import (
	cmp2 "github.com/MusicDin/kubitect/cli/pkg/utils/cmp"
	"testing"
)

func MockEvent(t *testing.T, eventType EventType, changes []cmp2.Change) Event {
	t.Helper()

	return Event{
		eType:   eventType,
		changes: changes,
		path:    t.TempDir(),
		msg:     "mock event",
		action:  cmp2.MODIFY,
	}
}
