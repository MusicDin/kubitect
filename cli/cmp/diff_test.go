package cmp

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SimpleStruct struct {
	value interface{}
}

var s = SimpleStruct{
	value: []SimpleStruct{
		{
			42,
		},
		{
			24,
		},
	},
}

func TestChanges(t *testing.T) {

	d, _ := Compare(true, true)
	assert.Empty(t, d.Changes())

	d, _ = Compare(true, false)
	assert.Equal(t, "(modify) : true -> false", fmt.Sprint(d.Changes()))

	expect := Changes{
		{
			Path:   []string{"value", "[0]", "value"},
			Before: 42,
			After:  nil,
			Action: DELETE,
		},
		{
			Path:   []string{"value", "[1]", "value"},
			Before: 24,
			After:  nil,
			Action: DELETE,
		},
	}

	d, _ = Compare(s, nil)
	assert.True(t, reflect.DeepEqual(expect, d.Changes()))
}

func TestOutputYaml(t *testing.T) {

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

func TestOutputJson(t *testing.T) {
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
	Path   string
	Action ActionType
}

func (e TestEvent) GetPath() string {
	return e.Path
}

func (e TestEvent) GetAction() ActionType {
	return e.Action
}

func TestEvents(t *testing.T) {

	events := []TestEvent{
		{
			Path:   "value",
			Action: DELETE,
		},
	}

	d, _ := Compare(s, nil)
	triggered := TriggerEvents(d, events)
	assert.Equal(t, events, triggered)

	d, _ = Compare(s, s)
	triggered = TriggerEvents(d, events)
	assert.Empty(t, triggered)
}
