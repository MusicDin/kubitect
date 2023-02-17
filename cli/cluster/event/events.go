package event

import (
	"github.com/MusicDin/kubitect/cli/utils/cmp"
)

type EventType string

const (
	OK         EventType = "ok"
	SCALE_UP   EventType = "scale_up"
	SCALE_DOWN EventType = "scale_down"

	// WARN change requires user permission to continue.
	WARN EventType = "warn"

	// BLOCK change prevents further actions on the cluster.
	BLOCK EventType = "block"
)

type Event struct {
	eType   EventType
	msg     string
	path    string
	paths   []string
	changes []cmp.Change
	action  cmp.ActionType
}

func (e Event) Action() cmp.ActionType {
	return e.action
}

func (e Event) Paths() []string {
	if len(e.path) > 0 {
		return []string{e.path}
	}

	return e.paths
}

func (e Event) Changes() []cmp.Change {
	return e.changes
}

func (e Event) Error() error {
	var paths []string

	for _, c := range e.changes {
		paths = append(paths, c.Path)
	}

	switch e.eType {
	case WARN:
		return NewConfigChangeWarning(e.msg, paths...)
	case BLOCK:
		return NewConfigChangeError(e.msg, paths...)
	}

	return nil
}

type Events []Event

// Blocking returns events of type BLOCK.
func (es Events) Blocking() Events {
	return es.OfType(BLOCK)
}

// Warns returns events of type WARN.
func (es Events) Warns() Events {
	return es.OfType(WARN)
}

// OfType returns events matching the given type.
func (es Events) OfType(t EventType) Events {
	var events Events

	for _, e := range es {
		if e.eType == t {
			events = append(events, e)
		}
	}

	return events
}

// Errors converts events to the utils.Errors.
func (es Events) Errors() []error {
	var err []error

	for _, e := range es {
		err = append(err, e.Error())
	}

	return err
}

// add adds the event with the corresponding change to the list.
// If an event with a matching action and path already exists in
// the list then the change is appended to the existing event.
func (es *Events) add(event Event, c cmp.Change) {
	for i, e := range *es {
		if e.action == event.action && e.path == event.path {
			(*es)[i].changes = append((*es)[i].changes, c)
			return
		}
	}

	event.changes = []cmp.Change{c}
	*es = append(*es, event)
}

// triggerEvents returns triggered events of the corresponding action.
func TriggerEvents(diff *cmp.DiffNode, events Events) Events {
	var trig Events

	cmp.TriggerEventsF(diff, events, trig.add)
	cc := cmp.ConflictingChanges(diff, events)

	if len(cc) > 0 {
		trig = append(trig, Event{
			eType:   BLOCK,
			msg:     "Disallowed changes.",
			changes: cc,
		})
	}

	return trig
}
