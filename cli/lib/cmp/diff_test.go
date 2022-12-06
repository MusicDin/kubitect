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
	changesEqual(t, expect, d.Changes())
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
	changesEqual(t, expect, d.Changes())
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
	changesEqual(t, expect, d.Changes())
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

func (e *TestEvent) Trigger(c Change) {
	e.triggerPaths = append(e.triggerPaths, c.Path)
}

func TestEvents_Trigger(t *testing.T) {
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
	TriggerEvents(d, events)
	assert.Equal(t, events[0].Paths(), events[0].Paths())
	assert.Equal(t, []string{"Value.0", "Value.1"}, events[0].triggerPaths)

	events[0].triggerPaths = []string{}

	d, _ = Compare(s, s)
	TriggerEvents(d, events)
	assert.Empty(t, events[0].triggerPaths)
}

func TestEvents_TriggerF(t *testing.T) {
	s := SimpleStruct{
		Value: []SimpleStruct{
			{42},
			{24},
		},
	}

	events := []TestEvent{
		{
			action: DELETE,
			path:   "Value.*",
		},
	}

	var matchedPaths []string

	mf := func(e TestEvent, c Change) {
		matchedPaths = append(matchedPaths, c.Path)
	}

	d, _ := Compare(s, nil)
	TriggerEventsF(d, events, mf)
	assert.Equal(t, []string{"Value.0", "Value.1"}, matchedPaths)

	matchedPaths = []string{}

	d, _ = Compare(s, s)
	TriggerEventsF(d, events, mf)
	assert.Empty(t, matchedPaths)
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
	mc := MatchingChanges(d, events)
	assert.Empty(t, ConflictingChanges(d, events))
	changesEqual(t, expect, mc)

	d, _ = Compare(s1, s2)
	cc := ConflictingChanges(d, events)
	assert.Empty(t, MatchingChanges(d, events))
	changeEquals(t, expect[0], cc[0])

	d, _ = Compare(s1, s1)
	assert.Empty(t, MatchingChanges(d, events))
	assert.Empty(t, ConflictingChanges(d, events))
}

func changesEqual(t *testing.T, a Changes, b Changes) {

	if len(a) != len(b) {
		assert.Fail(t, "Changes length differs.", "expected: %v\ngot: %v", len(a), len(b))
	}

	for i := range a {
		e := diffChanges(a[i], b[i], i)

		if len(e) > 0 {
			assert.Fail(t, "Changes differ!", e)
		}
	}
}

func changeEquals(t *testing.T, a Change, b Change) {
	e := diffChanges(a, b, -1)

	if len(e) > 0 {
		assert.Fail(t, "Change differs!", e)
	}
}

func diffChanges(a Change, b Change, i int) string {

	type diff struct {
		key      string
		expected interface{}
		got      interface{}
	}

	var diffs []diff

	if a.Path != b.Path {
		diffs = append(diffs, diff{
			key:      "Path",
			expected: a.Path,
			got:      b.Path,
		})
	}

	if a.StructPath != b.StructPath {
		diffs = append(diffs, diff{
			key:      "StructPath",
			expected: a.StructPath,
			got:      b.StructPath,
		})
	}

	if a.GenericPath != b.GenericPath {
		diffs = append(diffs, diff{
			key:      "GenericPath",
			expected: a.GenericPath,
			got:      b.GenericPath,
		})
	}

	if a.Before != b.Before {
		diffs = append(diffs, diff{
			key:      "Before",
			expected: a.Before,
			got:      b.Before,
		})
	}

	if a.After != b.After {
		diffs = append(diffs, diff{
			key:      "After",
			expected: a.After,
			got:      b.After,
		})
	}

	if a.Action != b.Action {
		diffs = append(diffs, diff{
			key:      "Action",
			expected: a.Action,
			got:      b.Action,
		})
	}

	if len(diffs) == 0 {
		return ""
	}

	var e string

	if i > 0 {
		e = fmt.Sprintf("Changes[%d]\n\n", i)
	}

	for _, d := range diffs {
		e += fmt.Sprintf("(%s)\n", d.key)
		e += fmt.Sprintf("expected: %v\n", d.expected)
		e += fmt.Sprintf("got:      %v\n\n", d.got)
	}

	return e
}
