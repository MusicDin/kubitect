package cmp

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func toInterface(v interface{}) interface{} {
	return v
}

func TestTypeMismatchError(t *testing.T) {

	a := 5
	b := true

	at := reflect.TypeOf(a).Kind()
	bt := reflect.TypeOf(b).Kind()

	expected := NewTypeMismatchError(at, bt)

	_, err := Compare(a, b)
	assert.Error(t, expected, err)
}

func TestBool(t *testing.T) {
	d, _ := Compare(true, true)
	assert.False(t, d.hasChanged())

	d, _ = Compare(true, false)
	assert.True(t, d.hasChanged())

	d, _ = Compare(true, nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, false)
	assert.True(t, d.hasChanged())
}

func TestInt(t *testing.T) {
	d, _ := Compare(42, 42)
	assert.False(t, d.hasChanged())

	d, _ = Compare(42, 24)
	assert.True(t, d.hasChanged())

	d, _ = Compare(42, nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, 42)
	assert.True(t, d.hasChanged())
}

func TestString(t *testing.T) {
	d, _ := Compare("test", "test")
	assert.False(t, d.hasChanged())

	d, _ = Compare("test", "tst")
	assert.True(t, d.hasChanged())

	d, _ = Compare("test", nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, "test")
	assert.True(t, d.hasChanged())
}

func TestPointer(t *testing.T) {
	a := "a"
	b := "b"

	d, _ := Compare(&a, &a)
	assert.False(t, d.hasChanged())

	d, _ = Compare(&a, &b)
	assert.True(t, d.hasChanged())

	d, _ = Compare(&a, nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, &b)
	assert.True(t, d.hasChanged())
}

func TestSlice(t *testing.T) {
	a := []int{1, 2}
	b := []int{2, 1}

	d, _ := Compare(a, a)
	assert.False(t, d.hasChanged())

	d, _ = Compare(a, b)
	assert.False(t, d.hasChanged())

	d, _ = Compare(a, nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, b)
	assert.True(t, d.hasChanged())
}

func TestSliceEmpty(t *testing.T) {
	a := []bool{true}
	b := []bool{}

	d, _ := Compare(a, a)
	assert.False(t, d.hasChanged())

	d, _ = Compare(a, b)
	assert.True(t, d.hasChanged())

	// Both slice of length 0 and nil represent slice zero value.
	d, _ = Compare(b, nil)
	assert.False(t, d.hasChanged())

	d, _ = Compare(&b, nil)
	assert.False(t, d.hasChanged())

	d, _ = Compare(nil, b)
	assert.False(t, d.hasChanged())
}

func TestSliceRespectOrder(t *testing.T) {
	a := []string{"a", "b"}
	b := []string{"b", "a"}

	cmp := NewComparator()
	cmp.RespectSliceOrder = true

	d, _ := cmp.Compare(a, b)
	assert.True(t, d.hasChanged())
}

func TestMap(t *testing.T) {
	m1 := map[string]interface{}{"test": 1}
	m2 := map[string]interface{}{"test": 2}
	m3 := map[string]interface{}{"test2": 1}
	m4 := map[string]int{"test": 1}

	d, _ := Compare(m1, m1)
	assert.False(t, d.hasChanged())

	d, _ = Compare(m1, m4)
	assert.False(t, d.hasChanged())

	d, _ = Compare(m1, m2)
	assert.True(t, d.hasChanged())

	d, _ = Compare(m3, m1)
	assert.True(t, d.hasChanged())

	d, _ = Compare(m1, nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, m2)
	assert.True(t, d.hasChanged())
}

func TestStructSimple(t *testing.T) {

	type SimpleStruct struct {
		value interface{}
	}

	a := SimpleStruct{42}
	b := SimpleStruct{24}

	d, _ := Compare(a, a)
	assert.False(t, d.hasChanged())

	d, _ = Compare(a, b)
	assert.True(t, d.hasChanged())

	d, _ = Compare(a, nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, a)
	assert.True(t, d.hasChanged())
}

func TestStructComplex(t *testing.T) {

	type SimpleStruct struct {
		value interface{}
	}

	type ComplexStruct struct {
		simple SimpleStruct
	}

	a := ComplexStruct{
		simple: SimpleStruct{
			value: 42,
		},
	}

	b := ComplexStruct{
		simple: SimpleStruct{
			value: 24,
		},
	}

	d, _ := Compare(a, a)
	assert.False(t, d.hasChanged())

	d, _ = Compare(a, b)
	assert.True(t, d.hasChanged())

	d, _ = Compare(&a, &b)
	assert.True(t, d.hasChanged())
}

func TestStructFieldName(t *testing.T) {

	type SimpleStruct struct {
		value interface{} `customTag:"v"`
	}

	a := SimpleStruct{
		value: 42,
	}

	b := SimpleStruct{
		value: 24,
	}

	d, _ := Compare(a, b)
	ch := d.Changes()
	assert.Equal(t, "value", ch[0].Path[0])

	cmp := NewComparator()
	cmp.TagName = "customTag"

	d, _ = cmp.Compare(a, b)
	ch = d.Changes()
	assert.Equal(t, "v", ch[0].Path[0])
}

func TestSliceId(t *testing.T) {

	type SimpleStruct struct {
		id interface{} `opt:",id"`
	}

	a := []SimpleStruct{
		{
			id: 10,
		},
		{
			id: 42,
		},
	}

	b := []SimpleStruct{
		{
			id: 42,
		},
		{
			id: 20,
		},
	}

	expect := []Change{
		{
			Path:   []string{"[0]", "id"},
			Before: 10,
			After:  nil,
			Action: DELETE,
		},
		{
			Path:   []string{"[1]", "id"},
			Before: nil,
			After:  20,
			Action: CREATE,
		},
	}

	cmp := NewComparator()

	d, _ := cmp.Compare(a, b)
	ch := d.Changes()

	assert.False(t, reflect.DeepEqual(ch, expect))

	cmp.TagName = "opt"

	expect[0].Path[0] = "10"
	expect[1].Path[0] = "20"

	d, _ = cmp.Compare(a, b)
	ch = d.Changes()

	assert.False(t, reflect.DeepEqual(ch, expect))
}
