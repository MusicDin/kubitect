package event

import (
	"cli/lib/cmp"
	"testing"
)

func MockEvent(t *testing.T, eventType EventType, changes []cmp.Change) Event {
	t.Helper()

	return Event{
		eType:   eventType,
		changes: changes,
		path:    t.TempDir(),
		msg:     "mock event",
		action:  cmp.MODIFY,
	}
}
