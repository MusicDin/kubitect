package cmp

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// changesEqual is a helper function that compares two list of changes and
// fails the test if they are not equal.
func changesEqual(t *testing.T, a []Change, b []Change) {
	if len(a) != len(b) {
		require.FailNow(t, "Changes length differs.", "expected: %v\ngot: %v", len(a), len(b))
	}

	res, err := Compare(a, b, Options{RespectSliceOrder: true, IgnoreEmptyChanges: true})
	require.NoError(t, err, "An unexpected error has occurred during change comparison.")

	if res.HasChanges() {
		diff := ""
		for _, c := range res.Changes() {
			diff += fmt.Sprintf("%s\n", c)
		}

		assert.Fail(t, "Changes differ!", diff)
	}
}

func TestChanges_Basic(t *testing.T) {
	res, err := Compare(nil, nil)
	require.NoError(t, err)
	assert.Empty(t, res.Changes())

	res, err = Compare(true, true)
	require.NoError(t, err)
	assert.Empty(t, res.Changes())

	res, err = Compare(true, false)
	require.NoError(t, err)
	assert.Equal(t, "[(modify) : true -> false]", fmt.Sprint(res.Changes()))
}

func TestChanges_StructCreate(t *testing.T) {
	type Struct struct{ Value any }

	s1 := Struct{"24"}

	expect := Change{
		Path:        "Value",
		StructPath:  "Value",
		ValueBefore: nil,
		ValueAfter:  "24",
		Type:        Create,
	}

	res, _ := Compare(nil, s1)
	changesEqual(t, []Change{expect}, res.Changes())
}

func TestChanges_Struct(t *testing.T) {
	type Struct struct{ Value any }

	s1 := Struct{"24"}
	s2 := Struct{"42"}

	expect := Change{
		Path:        "Value",
		StructPath:  "Value",
		ValueBefore: "24",
		ValueAfter:  "42",
		Type:        Modify,
	}

	res, _ := Compare(s1, s2)
	changesEqual(t, []Change{expect}, res.Changes())

	expect = Change{
		Path:        "Value",
		StructPath:  "Value",
		ValueBefore: "24",
		ValueAfter:  nil,
		Type:        Delete,
	}

	res, _ = Compare(s1, nil)
	changesEqual(t, []Change{expect}, res.Changes())

	expect = Change{
		Path:        "Value",
		StructPath:  "Value",
		ValueBefore: nil,
		ValueAfter:  "42",
		Type:        Create,
	}

	res, _ = Compare(nil, s2)
	changesEqual(t, []Change{expect}, res.Changes())
}

func TestDistinctChanges_Map(t *testing.T) {
	type Struct struct{ Value any }

	s1 := map[string]map[string]string{"A": {"a1": "abc"}}
	s2 := map[string]map[string]string{"A": {"a1": "xxx"}}
	empty := map[string]map[string]string{}

	opts := Options{PopulateAllNodes: true, IgnoreEmptyChanges: true}

	// Modification.
	expect := Change{
		Type:        Modify,
		Path:        "A.a1",
		StructPath:  "A.a1",
		ValueBefore: "abc",
		ValueAfter:  "xxx",
	}

	res, _ := Compare(s1, s2, opts)
	changesEqual(t, []Change{expect}, res.DistinctChanges())

	// Removal.
	expect = Change{
		Type:        Delete,
		Path:        "A",
		StructPath:  "A",
		ValueBefore: map[string]string{"a1": "abc"},
		ValueAfter:  nil,
	}

	res, _ = Compare(s1, empty, opts)
	changesEqual(t, []Change{expect}, res.DistinctChanges())

	// Addition.
	expect = Change{
		Type:        Create,
		Path:        "",
		StructPath:  "",
		ValueBefore: nil,
		ValueAfter:  s2,
	}

	res, _ = Compare(nil, s2, opts)
	changesEqual(t, []Change{expect}, res.DistinctChanges())
}

func TestDistinctChanges_Slice(t *testing.T) {
	type Struct struct {
		Id    string `cmp:",id"`
		Value any
	}

	s1 := []Struct{
		{"A", []Struct{{"1", "abc"}, {"2", "123"}}},
		{"B", []Struct{{"1", "abc"}, {"2", "xyz"}}},
	}

	s2 := []Struct{
		{"A", []Struct{{"2", "xyz"}, {"1", "abc"}}},
		{"X", []Struct{{"1", "abc"}, {"2", "xyz"}}},
	}

	expect := []Change{
		{
			Type:        Modify,
			Path:        "A.Value.2.Value",
			StructPath:  "A.Value.2.Value",
			ValueBefore: "123",
			ValueAfter:  "xyz",
		},
		{
			Type:        Delete,
			Path:        "B",
			StructPath:  "B",
			ValueBefore: s1[1],
			ValueAfter:  nil,
		},
		{
			Type:        Create,
			Path:        "X",
			StructPath:  "X",
			ValueBefore: nil,
			ValueAfter:  s2[1],
		},
	}

	res, _ := Compare(s1, s2, Options{PopulateAllNodes: true})
	changesEqual(t, expect, res.DistinctChanges())
}

func TestOutputYaml_Basic(t *testing.T) {
	res, _ := Compare(true, true)
	assert.Equal(t, "true", res.ToYaml())

	res, _ = Compare(true, false)
	assert.Equal(t, "true -> false", res.ToYaml())

	res, _ = Compare(false, true)
	assert.Equal(t, "false -> true", res.ToYaml())

	res, _ = Compare("123", "321")
	assert.Equal(t, `"123" -> "321"`, res.ToYaml())
}

func TestOutputYaml_Basic_DiffOnly(t *testing.T) {
	opts := FormatOptions{ShowDiffOnly: true}

	res, _ := Compare(true, true)
	assert.Equal(t, "", res.ToYaml(opts))

	res, _ = Compare(true, false)
	assert.Equal(t, "true -> false", res.ToYaml(opts))

	res, _ = Compare(false, true)
	assert.Equal(t, "false -> true", res.ToYaml(opts))

	res, _ = Compare("123", "321")
	assert.Equal(t, `"123" -> "321"`, res.ToYaml(opts))
}

func TestOutputYaml_Basic_WithTypePrefix(t *testing.T) {
	opts := FormatOptions{ShowChangeTypePrefix: true}

	res, _ := Compare(true, true)
	assert.Equal(t, "  │ true", res.ToYaml(opts))

	res, _ = Compare(true, false)
	assert.Equal(t, "~ │ true -> false", res.ToYaml(opts))

	res, _ = Compare(false, true)
	assert.Equal(t, "~ │ false -> true", res.ToYaml(opts))

	res, _ = Compare("123", "321")
	assert.Equal(t, `~ │ "123" -> "321"`, res.ToYaml(opts))
}

func TestOutputYaml_Struct(t *testing.T) {
	opts := FormatOptions{ShowDiffOnly: true, ShowChangeTypePrefix: true}

	type Struct struct{ Value any }

	s1 := Struct{Value: "42"}
	s2 := Struct{Value: "24"}

	// Modification.
	res, err := Compare(s1, s2)
	require.NoError(t, err)
	assert.Equal(t, `~ │ Value: "42" -> "24"`, res.ToYaml(opts))

	// Removal.
	res, err = Compare(s1, nil)
	require.NoError(t, err)
	assert.Equal(t, `- │ Value: "42"`, res.ToYaml(opts))

	// Addition.
	res, err = Compare(nil, s1)
	require.NoError(t, err)
	assert.Equal(t, `+ │ Value: "42"`, res.ToYaml(opts))

	// No change.
	res, err = Compare(s1, s1)
	require.NoError(t, err)
	assert.Equal(t, "", res.ToYaml(opts))
}

func TestOutputYaml_Map(t *testing.T) {
	opts := FormatOptions{ShowDiffOnly: true, ShowChangeTypePrefix: true}

	s1 := map[string]string{"a": "abc", "b": "123"}
	s2 := map[string]string{"a": "123", "b": "abc"}

	// Modification.
	expect := strings.TrimSpace(`
~ │ a: "abc" -> "123"
~ │ b: "123" -> "abc"
	`)

	res, err := Compare(s1, s2)
	require.NoError(t, err)
	assert.Equal(t, expect, res.ToYaml(opts))

	// Addition.
	expect = strings.TrimSpace(`
+ │ a: "abc"
+ │ b: "123"
	`)

	res, err = Compare(nil, s1)
	require.NoError(t, err)
	assert.Equal(t, expect, res.ToYaml(opts))

	// Removal.
	expect = strings.TrimSpace(`
- │ a: "abc"
- │ b: "123"
	`)

	res, err = Compare(s1, nil)
	require.NoError(t, err)
	assert.Equal(t, expect, res.ToYaml(opts))

	// No change.
	res, err = Compare(s1, s1)
	require.NoError(t, err)
	assert.Equal(t, "", res.ToYaml(opts))
}

func TestOutputYaml_List(t *testing.T) {
	opts := FormatOptions{ShowDiffOnly: true, ShowChangeTypePrefix: true}

	s1 := []string{"abc", "123"}
	s2 := []string{"123", "abc"}

	// Modification.
	expect := strings.TrimSpace(`
~ │ - "abc" -> "123"
~ │ - "123" -> "abc"
	`)

	res, err := Compare(s1, s2, Options{RespectSliceOrder: true})
	require.NoError(t, err)
	assert.Equal(t, expect, res.ToYaml(opts))

	// Removal.
	expect = strings.TrimSpace(`
- │ - "abc"
- │ - "123"
	`)

	res, err = Compare(s1, nil, Options{RespectSliceOrder: true})
	require.NoError(t, err)
	assert.Equal(t, expect, res.ToYaml(opts))

	// Addition.
	expect = strings.TrimSpace(`
+ │ - "abc"
+ │ - "123"
	`)

	res, err = Compare(nil, s1, Options{RespectSliceOrder: true})
	require.NoError(t, err)
	assert.Equal(t, expect, res.ToYaml(opts))

	// No change.
	res, err = Compare(s1, s1, Options{RespectSliceOrder: true})
	require.NoError(t, err)
	assert.Equal(t, "", res.ToYaml(opts))
}

func TestOutputYaml_ListOfLists(t *testing.T) {
	opts := FormatOptions{ShowDiffOnly: true, ShowChangeTypePrefix: true}

	s1 := [][]string{{"abc"}, {"123"}}
	s2 := [][]string{nil, {"abc", "cba"}}

	// Modification, removal and addition.
	expect := strings.TrimSpace(`
- │ - - "abc"
~ │ - - "123" -> "abc"
+ │   - "cba"
	`)

	res, err := Compare(s1, s2, Options{RespectSliceOrder: true})
	require.NoError(t, err)
	assert.Equal(t, expect, res.ToYaml(opts))

	// Removal.
	expect = strings.TrimSpace(`
- │ - - "abc"
- │ - - "123"
	`)

	res, err = Compare(s1, nil, Options{RespectSliceOrder: true})
	require.NoError(t, err)
	assert.Equal(t, expect, res.ToYaml(opts))

	// Addition.
	expect = strings.TrimSpace(`
+ │ - - "abc"
+ │ - - "123"
	`)

	res, err = Compare(nil, s1, Options{RespectSliceOrder: true})
	require.NoError(t, err)
	assert.Equal(t, expect, res.ToYaml(opts))

	// No change.
	res, err = Compare(s1, s1, Options{RespectSliceOrder: true})
	require.NoError(t, err)
	assert.Equal(t, "", res.ToYaml(opts))
}

func TestOutputYaml_ListOfStructs(t *testing.T) {
	opts := FormatOptions{ShowDiffOnly: true, ShowChangeTypePrefix: true}

	type Struct struct{ Value any }

	s1 := Struct{Value: []Struct{{"42"}, {24}}}
	s2 := Struct{Value: []Struct{{"24"}, {42}}}

	// Modification.
	expect := strings.TrimSpace(`
~ │ Value:
- │   - Value: "42"
- │   - Value: 24
+ │   - Value: "24"
+ │   - Value: 42
	`)

	res, err := Compare(s1, s2)
	require.NoError(t, err)
	assert.Equal(t, expect, res.ToYaml(opts))

	// Removal.
	expect = strings.TrimSpace(`
- │ Value:
- │   - Value: "42"
- │   - Value: 24
	`)

	res, err = Compare(s1, nil)
	require.NoError(t, err)
	assert.Equal(t, expect, res.ToYaml(opts))

	// Addition.
	expect = strings.TrimSpace(`
+ │ Value:
+ │   - Value: "42"
+ │   - Value: 24
	`)

	res, err = Compare(nil, s1)
	require.NoError(t, err)
	assert.Equal(t, expect, res.ToYaml(opts))

	// No change.
	res, err = Compare(s1, s1)
	require.NoError(t, err)
	assert.Equal(t, "", res.ToYaml(opts))
}

func TestOutputYaml_ListOfStructs_NoOptions(t *testing.T) {
	type Struct struct{ Value any }

	s1 := Struct{Value: []Struct{{"42"}, {24}}}
	s2 := Struct{Value: []Struct{{"42"}, {24}}}

	expect := strings.TrimSpace(`
Value:
  - Value: "42"
  - Value: 24
	`)

	// Modification.
	res, _ := Compare(s1, s2)
	assert.Equal(t, expect, res.ToYaml())

	// Removal.
	res, _ = Compare(s1, nil)
	assert.Equal(t, expect, res.ToYaml())

	// Addition.
	res, _ = Compare(nil, s1)
	assert.Equal(t, expect, res.ToYaml())

	// No change.
	res, _ = Compare(s1, s1)
	assert.Equal(t, expect, res.ToYaml())
}

func TestOutputYaml_ListOfStructs_ShowSliceId(t *testing.T) {
	opts := FormatOptions{ShowDiffOnly: true, ShowChangeTypePrefix: true}

	type Struct struct {
		Id   string `cmp:",id"`
		List []Struct
	}

	s1 := Struct{List: []Struct{{Id: "42"}}}
	s2 := Struct{List: []Struct{{Id: "42", List: []Struct{{Id: "24"}}}}}

	expect := strings.TrimSpace(`
~ │ List:
  │   - Id: "42"
+ │     List:
+ │       - Id: "24"
+ │         List: <nil>
	`)

	res, _ := Compare(s1, s2)
	assert.Equal(t, expect, res.ToYaml(opts))
}

func TestOutputYaml_ComplexType(t *testing.T) {
	opts := FormatOptions{ShowChangeTypePrefix: true}

	type Struct struct {
		V1 any
		V2 any
	}

	s1 := Struct{V1: [][]string{{"xyz"}, {"abc"}, {"123"}}, V2: Struct{V1: map[string][]string{"1": {"foo"}, "2": {"bar"}}, V2: "bar"}}
	s2 := Struct{V1: [][]string{{"xyz"}, nil, {"abc", "cba"}}, V2: Struct{V1: map[string][]string{"1": {"bar", "foo"}, "3": {"bar1"}}, V2: "foobar"}}

	// Modification, removal and addition.
	expect := strings.TrimSpace(`
~ │ V1:
  │   - - "xyz"
- │   - - "abc"
~ │   - - "123" -> "abc"
+ │     - "cba"
~ │ V2:
~ │   V1:
~ │     1:
~ │       - "foo" -> "bar"
+ │       - "foo"
- │     2:
- │       - "bar"
+ │     3:
+ │       - "bar1"
~ │   V2: "bar" -> "foobar"
	`)

	res, err := Compare(s1, s2, Options{RespectSliceOrder: true})
	require.NoError(t, err)
	assert.Equal(t, expect, res.ToYaml(opts))

	// Removal.
	expect = strings.TrimSpace(`
- │ V1:
- │   - - "xyz"
- │   - - "abc"
- │   - - "123"
- │ V2:
- │   V1:
- │     1:
- │       - "foo"
- │     2:
- │       - "bar"
- │   V2: "bar"
	`)

	res, err = Compare(s1, nil, Options{RespectSliceOrder: true})
	require.NoError(t, err)
	assert.Equal(t, expect, res.ToYaml(opts))

	// Addition.
	expect = strings.TrimSpace(`
+ │ V1:
+ │   - - "xyz"
+ │   - - "abc"
+ │   - - "123"
+ │ V2:
+ │   V1:
+ │     1:
+ │       - "foo"
+ │     2:
+ │       - "bar"
+ │   V2: "bar"
	`)

	res, err = Compare(nil, s1, Options{RespectSliceOrder: true})
	require.NoError(t, err)
	assert.Equal(t, expect, res.ToYaml(opts))

	// No change.
	res, err = Compare(s1, s1, Options{RespectSliceOrder: true})
	require.NoError(t, err)
	assert.Equal(t, "", res.ToYaml(FormatOptions{ShowDiffOnly: true}))
}

func TestOutputYaml_IgnoreEmptyChanges(t *testing.T) {
	type Struct struct{ Value *string }

	v := "42"
	s1 := Struct{}
	s2 := Struct{Value: &v}

	res, err := Compare(s1, nil, Options{IgnoreEmptyChanges: true})
	require.NoError(t, err)
	assert.Equal(t, false, res.HasChanges())
	assert.Equal(t, "", res.ToYaml(FormatOptions{ShowChangeTypePrefix: true}))

	res, err = Compare(s2, nil, Options{IgnoreEmptyChanges: true})
	require.NoError(t, err)
	assert.Equal(t, `Value: "42"`, res.ToYaml())
}
