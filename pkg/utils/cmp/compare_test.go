package cmp

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ptr returns a pointer to the value.
func ptr(v any) *any {
	return &v
}

type Test struct {
	Value1  any
	Value2  any
	Equal   bool
	Options Options
}

func runTests(t *testing.T, tests []Test) {
	for _, test := range tests {
		v1 := test.Value1
		v2 := test.Value2

		res, err := Compare(v1, v2, test.Options)
		require.NoErrorf(t, err, "%q: Comparison of '%v' and '%v' produced an unexpected error!", v1, v2)

		if test.Equal && res.HasChanges() {
			assert.Fail(t, fmt.Sprintf("%q: Values '%v' and '%v' should be equal!", t.Name(), v1, v2))
		} else if !test.Equal && !res.HasChanges() {
			assert.Fail(t, fmt.Sprintf("%q: Values '%v' and '%v' should not be equal!", t.Name(), v1, v2))
		}
	}
}

func TestTypeMismatchError(t *testing.T) {
	tests := []struct {
		Value1 any
		Value2 any
	}{
		// Type mismatches.
		{Value1: true, Value2: "true"},
		{Value1: 123, Value2: int32(123)},
		{Value1: 123, Value2: uint8(123)},
		{Value1: 1.3, Value2: float32(1.3)},
		{Value1: 123, Value2: complex128(123)},
		{Value1: int(123), Value2: int32(123)},
		{Value1: int32(123), Value2: int64(123)},
		{Value1: uint32(123), Value2: uint64(123)},
		{Value1: float32(1.3), Value2: float64(1.3)},
		{Value1: complex64(1.3), Value2: complex128(1.3)},
		{Value1: []int{}, Value2: &[]int{}},
		{Value1: []int{}, Value2: map[int]int{}},
		// Unsupported comparisons.
		{Value1: uintptr(123), Value2: uintptr(123)},
		{Value1: func() {}, Value2: func() {}},
		{Value1: make(chan int), Value2: make(chan int)},
	}

	for _, test := range tests {
		v1 := test.Value1
		v2 := test.Value2
		k1 := reflect.TypeOf(v1).Kind()
		k2 := reflect.TypeOf(v2).Kind()

		expect := NewTypeMismatchError(k1, k2).Error()

		_, err := Compare(v1, v2)
		assert.EqualError(t, err, expect, "Values '%v' and '%v' have the same type (%v), but it should be different!", v1, v2, k1)
	}
}

func TestBasic(t *testing.T) {
	tests := []Test{
		// Nil
		{Value1: nil, Value2: nil, Equal: true},
		// Bool
		{Value1: true, Value2: true, Equal: true},
		{Value1: false, Value2: false, Equal: true},
		{Value1: true, Value2: false, Equal: false},
		{Value1: true, Value2: nil, Equal: false},
		// String
		{Value1: "test", Value2: "test", Equal: true},
		{Value1: "test", Value2: "abcd", Equal: false},
		{Value1: "test", Value2: nil, Equal: false},
		{Value1: nil, Value2: "", Equal: false},
		// Int
		{Value1: 123, Value2: 123, Equal: true},
		{Value1: -12, Value2: -12, Equal: true},
		{Value1: -12, Value2: 123, Equal: false},
		{Value1: nil, Value2: 123, Equal: false},
		{Value1: -12, Value2: nil, Equal: false},
		{Value1: int(42), Value2: int(42), Equal: true},
		{Value1: int(42), Value2: int(77), Equal: false},
		{Value1: int8(42), Value2: int8(42), Equal: true},
		{Value1: int8(42), Value2: int8(77), Equal: false},
		{Value1: int16(42), Value2: int16(42), Equal: true},
		{Value1: int16(42), Value2: int16(77), Equal: false},
		{Value1: int32(42), Value2: int32(42), Equal: true},
		{Value1: int32(42), Value2: int32(77), Equal: false},
		{Value1: int64(42), Value2: int64(42), Equal: true},
		{Value1: int64(42), Value2: int64(77), Equal: false},
		// Uint
		{Value1: uint(42), Value2: uint(42), Equal: true},
		{Value1: uint(42), Value2: uint(77), Equal: false},
		{Value1: uint8(42), Value2: uint8(42), Equal: true},
		{Value1: uint8(42), Value2: uint8(77), Equal: false},
		{Value1: uint16(42), Value2: uint16(42), Equal: true},
		{Value1: uint16(42), Value2: uint16(77), Equal: false},
		{Value1: uint32(42), Value2: uint32(42), Equal: true},
		{Value1: uint32(42), Value2: uint32(77), Equal: false},
		{Value1: uint64(42), Value2: uint64(42), Equal: true},
		{Value1: uint64(42), Value2: uint64(77), Equal: false},
		// Float
		{Value1: float32(4.2), Value2: float32(4.2), Equal: true},
		{Value1: float32(4.2), Value2: float32(7.7), Equal: false},
		{Value1: float64(4.2), Value2: float64(4.2), Equal: true},
		{Value1: float64(4.2), Value2: float64(7.7), Equal: false},
		// Complex
		{Value1: complex64(4.2), Value2: complex64(4.2), Equal: true},
		{Value1: complex64(4.2), Value2: complex64(7.7), Equal: false},
		{Value1: complex128(4.2), Value2: complex128(4.2), Equal: true},
		{Value1: complex128(4.2), Value2: complex128(7.7), Equal: false},
	}

	runTests(t, tests)
}

func TestPointer(t *testing.T) {
	v1 := ptr("a")
	v2 := ptr("b")

	res, _ := Compare(v1, v1)
	assert.False(t, res.HasChanges())

	res, _ = Compare(v1, v2)
	assert.True(t, res.HasChanges())

	res, _ = Compare(v1, nil)
	assert.True(t, res.HasChanges())

	res, _ = Compare(nil, v2)
	assert.True(t, res.HasChanges())
}

func TestPointer_Empty(t *testing.T) {
	var p *int

	res, _ := Compare(p, p)
	assert.False(t, res.HasChanges())

	res, _ = Compare(p, nil)
	assert.True(t, res.HasChanges())

	res, _ = Compare(nil, p)
	assert.True(t, res.HasChanges())

	res, _ = Compare(p, nil, Options{IgnoreEmptyChanges: true})
	assert.True(t, res.HasChanges())
}

func TestStruct(t *testing.T) {
	type Struct struct{ Value any }

	s1 := Struct{42}
	s2 := Struct{24}

	res, _ := Compare(s1, s1)
	assert.False(t, res.HasChanges())

	res, _ = Compare(s1, s2)
	assert.True(t, res.HasChanges())

	res, _ = Compare(s1, nil)
	assert.True(t, res.HasChanges())

	res, _ = Compare(nil, s1)
	assert.True(t, res.HasChanges())
}

func TestStruct_Empty(t *testing.T) {
	type Struct struct{ Value any }

	s := Struct{}

	res, _ := Compare(s, s)
	assert.False(t, res.HasChanges())

	res, _ = Compare(s, nil)
	assert.True(t, res.HasChanges())

	res, _ = Compare(nil, &s)
	assert.True(t, res.HasChanges())

	res, _ = Compare(s, nil, Options{IgnoreEmptyChanges: true})
	assert.False(t, res.HasChanges())

	res, _ = Compare(nil, &s, Options{IgnoreEmptyChanges: true})
	assert.False(t, res.HasChanges())
}

func TestStruct_Complex(t *testing.T) {
	type Struct struct{ Value any }

	s1 := Struct{Value: Struct{Value: 42}}
	s2 := Struct{Value: Struct{Value: 24}}

	res, _ := Compare(s1, s1)
	assert.False(t, res.HasChanges())

	res, _ = Compare(s1, nil)
	changes := res.Changes()
	require.Len(t, changes, 1)
	assert.Equal(t, Delete, changes[0].Type)

	res, _ = Compare(nil, s2)
	changes = res.Changes()
	require.Len(t, changes, 1)
	assert.Equal(t, Create, changes[0].Type)

	res, _ = Compare(&s1, &s2)
	changes = res.Changes()
	require.Len(t, changes, 1)
	assert.Equal(t, Modify, changes[0].Type)
}

func TestStruct_UnexportedField(t *testing.T) {
	type Struct struct{ value any }

	s1 := Struct{value: 42}
	s2 := Struct{value: 24}

	res, _ := Compare(s1, s1)
	assert.False(t, res.HasChanges())

	res, _ = Compare(s1, nil)
	assert.True(t, res.HasChanges())

	res, _ = Compare(nil, &s2)
	assert.True(t, res.HasChanges())

	// There should be no changes, since unexported fields are not
	// compared.
	res, _ = Compare(s1, s2)
	assert.False(t, res.HasChanges())

	res, _ = Compare(&s1, &s2)
	assert.False(t, res.HasChanges())
}

func TestStruct_FieldName(t *testing.T) {
	type Struct struct {
		Value any `customTag:"v"`
	}

	s1 := Struct{42}
	s2 := Struct{24}

	res, _ := Compare(s1, s2)
	assert.Equal(t, "Value", res.Changes()[0].Path)

	res, _ = Compare(s1, s2, Options{Tag: "customTag"})
	assert.Equal(t, "v", res.Changes()[0].Path)
}

func TestStruct_SkipField(t *testing.T) {
	type Struct struct {
		Value any `customTag:"-"`
	}

	s1 := Struct{42}
	s2 := Struct{24}

	res, _ := Compare(s1, s2, Options{Tag: "customTag"})
	assert.False(t, res.HasChanges())
}

func TestStruct_PopulateNodes(t *testing.T) {
	type Struct struct{ Value any }

	s := Struct{Value: Struct{Value: 42}}

	res, err := Compare(s, nil, Options{PopulateAllNodes: true})
	require.NoError(t, err)

	changes := res.Changes()
	require.Len(t, changes, 1)
	assert.Equal(t, changes[0].ValueBefore, 42)
	assert.Equal(t, changes[0].ValueAfter, nil)
}

func TestMap(t *testing.T) {
	m1 := map[string]any{"test": 1}
	m2 := map[string]any{"test": 2}
	m3 := map[string]any{"test2": 1}
	m4 := map[string]int{"test": 1}

	res, _ := Compare(m1, m1)
	assert.False(t, res.HasChanges())

	res, _ = Compare(m1, m2)
	assert.True(t, res.HasChanges())

	res, _ = Compare(m3, m1)
	assert.True(t, res.HasChanges())

	res, _ = Compare(m1, nil)
	assert.True(t, res.HasChanges())

	res, _ = Compare(nil, m2)
	assert.True(t, res.HasChanges())

	res, err := Compare(m1, m4)
	expect := "Compared values have either unsupported or mismatched types (interface <> int)."
	assert.EqualError(t, err, expect)
}

func TestMap_Empty(t *testing.T) {
	m := map[any]bool{}

	res, _ := Compare(m, m)
	assert.False(t, res.HasChanges())

	res, _ = Compare(m, nil)
	assert.True(t, res.HasChanges())

	res, _ = Compare(nil, &m)
	assert.True(t, res.HasChanges())

	res, _ = Compare(m, nil, Options{IgnoreEmptyChanges: true})
	assert.False(t, res.HasChanges())

	res, _ = Compare(nil, &m, Options{IgnoreEmptyChanges: true})
	assert.False(t, res.HasChanges())
}

func TestArray(t *testing.T) {
	s1 := [2]int{1, 2}
	s2 := [2]int{2, 1}

	res, err := Compare(s1, s1)
	require.NoError(t, err)
	assert.False(t, res.HasChanges())

	res, _ = Compare(s1, s2)
	assert.False(t, res.HasChanges())

	res, _ = Compare(s1, nil)
	assert.True(t, res.HasChanges())

	res, _ = Compare(nil, s2)
	assert.True(t, res.HasChanges())
}

func TestArray_Empty(t *testing.T) {
	s := [0]any{}

	res, _ := Compare(s, s)
	assert.False(t, res.HasChanges())

	res, _ = Compare(s, nil)
	assert.True(t, res.HasChanges())

	res, _ = Compare(nil, &s)
	assert.True(t, res.HasChanges())

	res, _ = Compare(s, nil, Options{IgnoreEmptyChanges: true})
	assert.False(t, res.HasChanges())

	res, _ = Compare(nil, &s, Options{IgnoreEmptyChanges: true})
	assert.False(t, res.HasChanges())
}

func TestSlice(t *testing.T) {
	s1 := []int{1, 2}
	s2 := []int{2, 1}

	res, _ := Compare(s1, s1)
	assert.False(t, res.HasChanges())

	res, _ = Compare(s1, s2)
	assert.False(t, res.HasChanges())

	res, _ = Compare(s1, nil)
	assert.True(t, res.HasChanges())

	res, _ = Compare(nil, s2)
	assert.True(t, res.HasChanges())
}

func TestSlice_Empty(t *testing.T) {
	s := []any{}

	res, _ := Compare(s, s)
	assert.False(t, res.HasChanges())

	res, _ = Compare(s, nil)
	assert.True(t, res.HasChanges())

	res, _ = Compare(nil, &s)
	assert.True(t, res.HasChanges())

	res, _ = Compare(s, nil, Options{IgnoreEmptyChanges: true})
	assert.False(t, res.HasChanges())

	res, _ = Compare(nil, &s, Options{IgnoreEmptyChanges: true})
	assert.False(t, res.HasChanges())
}

func TestSlice_SliceOrder(t *testing.T) {
	s1 := []string{"a", "b"}
	s2 := []string{"b", "a"}

	res, _ := Compare(s1, s2)
	assert.False(t, res.HasChanges())

	res, _ = Compare(s1, s2, Options{RespectSliceOrder: true})
	assert.True(t, res.HasChanges())
}

func TestSlice_SliceId(t *testing.T) {
	type Struct struct {
		Id any `opt:"id,id"`
	}

	s1 := []Struct{{Id: 10}, {Id: 42}}
	s2 := []Struct{{Id: 42}, {Id: 20}}

	expect := []Change{
		{
			Path:        "0.Id",
			StructPath:  "0.Id",
			ValueBefore: 10,
			ValueAfter:  nil,
			Type:        Delete,
		},
		{
			Path:        "2.Id",
			StructPath:  "2.Id",
			ValueBefore: nil,
			ValueAfter:  20,
			Type:        Create,
		},
	}

	res, _ := Compare(s1, s2)
	changesEqual(t, expect, res.Changes())

	expect[0].Path = "10.id"
	expect[1].Path = "20.id"
	expect[0].StructPath = "10.Id"
	expect[1].StructPath = "20.Id"

	res, _ = Compare(s1, s2, Options{Tag: "opt"})
	changesEqual(t, expect, res.Changes())
}

func TestSlice_SliceIdToNilElem(t *testing.T) {
	type Struct struct {
		Id any `cmp:"id,id"`
	}

	s := []Struct{{Id: 10}, {Id: 42}}

	expect := []Change{
		{
			Path:        "10.id",
			StructPath:  "10.Id",
			ValueBefore: 10,
			ValueAfter:  nil,
			Type:        Delete,
		},
		{
			Path:        "42.id",
			StructPath:  "42.Id",
			ValueBefore: 42,
			ValueAfter:  nil,
			Type:        Delete,
		},
	}

	res, _ := Compare(s, nil)
	changesEqual(t, expect, res.Changes())
}

func TestSlice_UnexportedField(t *testing.T) {
	type Struct struct{ value any }

	s := []Struct{{42}, {24}}

	res, _ := Compare(s, s)
	assert.False(t, res.HasChanges())

	res, _ = Compare(&s, nil)
	assert.True(t, res.HasChanges())

	res, _ = Compare(nil, s)
	assert.True(t, res.HasChanges())
}

// Ensures that Any/Unknown change type is not propagated across children
// nodes.
func Test_ChangeTypePropagation(t *testing.T) {
	type Struct struct {
		A string
		B []string // Produces ChangeType 'Any'
	}

	s1 := Struct{A: "test"}
	s2 := Struct{A: "test"}

	res, _ := Compare(s1, s2)
	assert.False(t, res.HasChanges())
}
