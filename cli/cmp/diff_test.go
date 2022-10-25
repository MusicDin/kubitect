package cmp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SimpleStruct struct {
	value interface{}
}

func TestChanges_Basic(t *testing.T) {
	d, _ := Compare(nil, nil)
	assert.Empty(t, d.Changes())

	d, _ = Compare(true, true)
	assert.Empty(t, d.Changes())

	d, _ = Compare(true, false)
	assert.Equal(t, "(modify) : true -> false", fmt.Sprint(d.Changes()))
}

func TestChanges_StructCreate(t *testing.T) {
	s1 := SimpleStruct{"24"}

	expect := Changes{
		{
			Path:        "value",
			GenericPath: "value",
			Before:      nil,
			After:       "24",
			Action:      CREATE,
		},
	}

	d, _ := Compare(nil, s1)
	assert.Equal(t, expect, d.Changes())
}

func TestChanges_StructDelete(t *testing.T) {
	s1 := SimpleStruct{"24"}

	expect := Changes{
		{
			Path:        "value",
			GenericPath: "value",
			Before:      "24",
			After:       nil,
			Action:      DELETE,
		},
	}

	d, _ := Compare(s1, nil)
	assert.Equal(t, expect, d.Changes())
}

func TestChanges_StructModify(t *testing.T) {
	s1 := SimpleStruct{"24"}
	s2 := SimpleStruct{"42"}

	expect := Changes{
		{
			Path:        "value",
			GenericPath: "value",
			Before:      "24",
			After:       "42",
			Action:      MODIFY,
		},
	}

	d, _ := Compare(s1, s2)
	assert.Equal(t, expect, d.Changes())
}

func TestOutput_Yaml(t *testing.T) {
	s := SimpleStruct{
		value: []SimpleStruct{
			{
				42,
			},
			{
				24,
			},
		},
	}

	d, _ := Compare(true, true)
	assert.Equal(t, "true", d.ToYaml())

	d, _ = Compare(true, false)
	assert.Equal(t, "true -> false", d.ToYaml())

	expect := "value: \n  - value: 42\n  - value: 24"

	d, _ = Compare(s, nil)
	assert.Equal(t, expect, d.ToYaml())

	d, _ = Compare(nil, s)
	assert.Equal(t, expect, d.ToYaml())
}

func TestOutput_Json(t *testing.T) {
	s := SimpleStruct{
		value: []SimpleStruct{
			{
				42,
			},
			{
				24,
			},
		},
	}

	d, _ := Compare(true, true)
	assert.Equal(t, "true", d.ToJson())

	d, _ = Compare(true, false)
	assert.Equal(t, "true -> false", d.ToJson())

	expect := "{\n  value: [\n    {\n      value: 42,\n    },\n    {\n      value: 24,\n    },\n  ],\n}"

	d, _ = Compare(s, nil)
	assert.Equal(t, expect, d.ToJson())

	d, _ = Compare(nil, s)
	assert.Equal(t, expect, d.ToJson())
}

type TestEvent struct {
	action       ActionType
	path         string
	triggerPaths []string
}

func (e TestEvent) Action() ActionType {
	return e.action
}

func (e TestEvent) Paths() []string {
	return []string{e.path}
}

func (e *TestEvent) TriggerPath(path string) {
	e.triggerPaths = append(e.triggerPaths, path)
}

func TestEvents_Triggered(t *testing.T) {
	s := SimpleStruct{
		value: []SimpleStruct{
			{
				42,
			},
			{
				24,
			},
		},
	}

	events := []*TestEvent{
		{
			action: DELETE,
			path:   "value",
		},
	}

	d, _ := Compare(s, nil)
	triggered := TriggerEvents(d, events)
	assert.Equal(t, events[0].Paths(), triggered[0].Paths())
	assert.Equal(t, events[0].Action(), triggered[0].Action())
	assert.Equal(t, []string{"value"}, triggered[0].triggerPaths)

	d, _ = Compare(s, s)
	triggered = TriggerEvents(d, events)
	assert.Empty(t, triggered)
}

func TestEvents_Changes(t *testing.T) {
	s1 := SimpleStruct{
		value: []SimpleStruct{
			{
				42,
			},
			{
				24,
			},
		},
	}

	s2 := SimpleStruct{
		value: []SimpleStruct{
			{
				24,
			},
		},
	}

	events := []*TestEvent{
		{
			action: DELETE,
			path:   "value",
		},
	}

	expect := Changes{
		{
			Path:        "value.[0].value",
			GenericPath: "value.[*].value",
			Before:      42,
			After:       nil,
			Action:      DELETE,
		},
		{
			Path:        "value.[1].value",
			GenericPath: "value.[*].value",
			Before:      24,
			After:       nil,
			Action:      DELETE,
		},
	}

	d, _ := Compare(s1, nil)
	assert.ElementsMatch(t, expect, MatchingChanges(d, events))
	assert.Empty(t, NonMatchingChanges(d, events))

	d, _ = Compare(s1, s2)
	assert.Empty(t, MatchingChanges(d, events))
	assert.ElementsMatch(t, Changes{expect[0]}, NonMatchingChanges(d, events))

	d, _ = Compare(s1, s1)
	assert.Empty(t, MatchingChanges(d, events))
	assert.Empty(t, NonMatchingChanges(d, events))
}
