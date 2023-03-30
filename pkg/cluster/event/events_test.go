package event

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/MusicDin/kubitect/pkg/config/modelconfig"
	"github.com/MusicDin/kubitect/pkg/utils/cmp"

	"github.com/stretchr/testify/assert"
)

func TestMockEvent(t *testing.T) {
	assert.NotNil(t, MockEvent(t, OK, cmp.MODIFY, nil))
}

func TestEvent_Path(t *testing.T) {
	var e Event

	e.paths = []string{"test", "123"}
	assert.Equal(t, e.paths, e.Paths())

	// if e.path is set, e.paths should be ignored
	e.path = "test_path"
	assert.Equal(t, []string{"test_path"}, e.Paths())
}

func TestEvent(t *testing.T) {
	changes := []cmp.Change{
		{
			Action: cmp.CREATE,
			Path:   "test",
		},
	}

	event := Event{
		eType:   OK,
		changes: changes,
		msg:     "Event msg",
		action:  cmp.CREATE,
	}

	assert.Equal(t, cmp.CREATE, event.Action())
	assert.Equal(t, changes, event.Changes())
	assert.Equal(t, nil, event.Error())
}

func TestEvent_ConfigWarning(t *testing.T) {
	event := Event{
		eType: WARN,
		msg:   "event msg",
		changes: []cmp.Change{
			{Path: "test"},
		},
	}

	assert.Equal(t, NewConfigChangeWarning("event msg", "test"), event.Error())
}

func TestEvent_ConfigError(t *testing.T) {
	event := Event{
		eType: BLOCK,
		changes: []cmp.Change{
			{Path: "test"},
		},
	}

	assert.NotEqual(t, NewConfigChangeError("event msg", "test"), event.Error())
}

func TestEvents_Type(t *testing.T) {
	eventOk := Event{eType: OK}
	eventWarn := Event{eType: WARN}
	eventBlock := Event{eType: BLOCK}

	events := Events{eventOk, eventWarn, eventBlock}

	assert.Equal(t, Events{eventOk}, events.OfType(OK))
	assert.Equal(t, Events{eventWarn}, events.OfType(WARN))
	assert.Equal(t, Events{eventWarn}, events.Warns())
	assert.Equal(t, Events{eventBlock}, events.OfType(BLOCK))
	assert.Equal(t, Events{eventBlock}, events.Blocking())
}

func TestEvents_Errors(t *testing.T) {
	events := Events{
		Event{eType: OK, msg: "ok"},
		Event{eType: WARN, msg: "warn"},
		Event{eType: BLOCK, msg: "block"},
	}

	expect := []error{
		nil,
		NewConfigChangeWarning("warn"),
		NewConfigChangeError("block"),
	}

	assert.Equal(t, expect, events.Errors())
}

func TestEvents_Add(t *testing.T) {
	var events Events

	e1 := Event{path: "path.test", msg: "event1", action: cmp.CREATE}
	e2 := Event{path: "path.test", msg: "event2", action: cmp.DELETE}

	change := cmp.Change{}

	events.add(e1, change)
	assert.Len(t, events, 1)
	assert.Len(t, events[0].changes, 1)

	// Add change to the already existing event
	events.add(e1, change)
	assert.Len(t, events, 1)
	assert.Len(t, events[0].changes, 2)

	// Add new event
	events.add(e2, change)
	assert.Len(t, events, 2)
	assert.Len(t, events[0].changes, 2)
	assert.Len(t, events[1].changes, 1)
}

func TestTriggerEvents(t *testing.T) {
	type S struct {
		Value string
	}

	s1 := S{"test1"}
	s2 := S{"test2"}

	events := Events{Event{path: "Value", eType: BLOCK}}

	diff, err := cmp.Compare(s1, s2)
	assert.NoError(t, err)

	actual := TriggerEvents(diff, events).Errors()
	expect := []error{NewConfigChangeError("", "Value")}

	assert.Equal(t, expect, actual)
}

func TestTriggerEvents_Conflicting(t *testing.T) {
	type S struct {
		Value string
		Next  *S
	}

	s1 := S{Value: "test1"}
	s2 := S{Value: "test2", Next: &s1}

	events := Events{Event{path: "Value", eType: BLOCK}}

	diff, err := cmp.Compare(s1, s2)
	assert.NoError(t, err)

	trig := TriggerEvents(diff, events)
	assert.Equal(t, BLOCK, trig[1].eType)
	assert.Equal(t, "Disallowed changes.", trig[1].msg)
	assert.Equal(t, "Next.Value", trig[1].changes[0].Path)
}

// TestEventPaths verifies that paths of the events (from events_list.go)
// are valid.
func TestEventPaths(t *testing.T) {
	var events []Event

	events = append(events, ModifyEvents...)
	events = append(events, ScaleEvents...)
	events = append(events, UpgradeEvents...)

	for _, e := range events {
		for _, p := range e.Paths() {
			validateConfigPath(t, p)
		}
	}
}

func validateConfigPath(t *testing.T, path string) {
	paths := strings.Split(path, ".")
	cType := reflect.TypeOf(modelconfig.Config{})

	pass := typePathExists(cType, paths...)

	if !pass {
		err := fmt.Sprintf("Change event path '%s' does not represent any '%s' object!", path, cType)
		assert.Fail(t, err)
	}
}

func typePathExists(t reflect.Type, path ...string) bool {
	if len(path) == 0 {
		return true
	}

	x := reflect.New(t)

	if x.Kind() == reflect.Pointer {
		x = reflect.Indirect(x)
	}

	switch x.Kind() {
	case reflect.Struct:
		for i := 0; i < x.NumField(); i++ {
			f := x.Type().Field(i)

			if f.Name != path[0] {
				continue
			}

			if len(path) == 1 {
				return true
			}

			if !f.IsExported() {
				return false
			}

			return typePathExists(f.Type, path[1:]...)
		}
	case reflect.Slice, reflect.Array:
		if path[0] == "*" {
			return typePathExists(x.Type().Elem(), path[1:]...)
		}
	}

	return false
}
