package cmp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SimpleStruct struct {
	Value interface{}
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
			Path:        "Value",
			StructPath:  "Value",
			GenericPath: "Value",
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
			Path:        "Value",
			StructPath:  "Value",
			GenericPath: "Value",
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
			Path:        "Value",
			StructPath:  "Value",
			GenericPath: "Value",
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
		Value: []SimpleStruct{
			{"42"},
			{24},
		},
	}

	d, _ := Compare(true, true)
	assert.Equal(t, "true", d.ToYaml())

	d, _ = Compare(true, false)
	assert.Equal(t, "true -> false", d.ToYaml())

	expect := "Value: \n  - Value: \"42\"\n  - Value: 24"

	d, _ = Compare(s, nil)
	assert.Equal(t, expect, d.ToYaml())

	d, _ = Compare(nil, s)
	assert.Equal(t, expect, d.ToYaml())
}

func TestOutput_YamlDiff(t *testing.T) {
	s := SimpleStruct{
		Value: []SimpleStruct{
			{"42"},
			{24},
		},
	}

	d, _ := Compare(true, true)
	assert.Equal(t, "", d.ToYamlDiff())

	d, _ = Compare(true, false)
	assert.Equal(t, "true -> false", d.ToYamlDiff())

	expect := "Value: \n  - Value: \"42\"\n  - Value: 24"

	d, _ = Compare(s, nil)
	assert.Equal(t, expect, d.ToYamlDiff())

	d, _ = Compare(nil, s)
	assert.Equal(t, expect, d.ToYamlDiff())
}

func TestOutput_YamlDiffComplex(t *testing.T) {
	type SimpleStruct struct {
		Id   string `cmp:",id"`
		List []SimpleStruct
	}

	s1 := SimpleStruct{
		List: []SimpleStruct{
			{Id: "42"},
		},
	}

	s2 := SimpleStruct{
		List: []SimpleStruct{
			{
				Id: "42",
				List: []SimpleStruct{
					{Id: "24"},
				},
			},
		},
	}

	expect := "List: \n  - Id: \"42\"\n    List: \n      - Id: \"24\"\n        List: <nil>"

	d, _ := Compare(s1, s2)
	assert.Equal(t, expect, d.ToYamlDiff())
}

func TestOutput_IgnoreEmptyChanges(t *testing.T) {
	type SimpleStruct struct {
		Id   *string `cmp:",id"`
		List []SimpleStruct
	}

	id := "42"

	cmp := NewComparator()
	cmp.IgnoreEmptyChanges = true

	d, _ := cmp.Compare(SimpleStruct{Id: &id}, nil)
	assert.Equal(t, "Id: \"42\"", d.ToYamlDiff())

	d, _ = cmp.Compare(SimpleStruct{}, nil)
	assert.Equal(t, "", d.ToYamlDiff())
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
		Value: []SimpleStruct{
			{42},
			{24},
		},
	}

	events := []*TestEvent{
		{
			action: DELETE,
			path:   "Value.*",
		},
	}

	d, _ := Compare(s, nil)
	triggered := TriggerEvents(d, events)
	assert.Equal(t, events[0].Paths(), triggered[0].Paths())
	assert.Equal(t, []string{"Value.0", "Value.1"}, triggered[0].triggerPaths)

	d, _ = Compare(s, s)
	triggered = TriggerEvents(d, events)
	assert.Empty(t, triggered)
}

func TestEvents_Changes(t *testing.T) {
	s1 := SimpleStruct{
		Value: []SimpleStruct{
			{42},
			{24},
		},
	}

	s2 := SimpleStruct{
		Value: []SimpleStruct{
			{24},
		},
	}

	events := []*TestEvent{
		{
			action: DELETE,
			path:   "Value",
		},
	}

	expect := Changes{
		{
			Path:        "Value.0.Value",
			StructPath:  "Value.0.Value",
			GenericPath: "Value.*.Value",
			Before:      42,
			After:       nil,
			Action:      DELETE,
		},
		{
			Path:        "Value.1.Value",
			StructPath:  "Value.1.Value",
			GenericPath: "Value.*.Value",
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
