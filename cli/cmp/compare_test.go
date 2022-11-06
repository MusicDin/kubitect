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

	_, err := Compare(a, b)
	assert.EqualError(t, err, NewTypeMismatchError(at, bt).Error())
}

func TestBasic_Bool(t *testing.T) {
	d, _ := Compare(true, true)
	assert.False(t, d.hasChanged())

	d, _ = Compare(true, false)
	assert.True(t, d.hasChanged())

	d, _ = Compare(true, nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, false)
	assert.True(t, d.hasChanged())
}

func TestBasic_String(t *testing.T) {
	d, _ := Compare("test", "test")
	assert.False(t, d.hasChanged())

	d, _ = Compare("test", "tst")
	assert.True(t, d.hasChanged())

	d, _ = Compare("test", nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, "test")
	assert.True(t, d.hasChanged())
}

func TestBasic_Int(t *testing.T) {
	d, _ := Compare(1, 1)
	assert.False(t, d.hasChanged())

	d, _ = Compare(1, 2)
	assert.True(t, d.hasChanged())

	d, _ = Compare(1, nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, 1)
	assert.True(t, d.hasChanged())
}

func TestBasic_IntX(t *testing.T) {
	d, _ := Compare(int8(1), int8(1))
	assert.False(t, d.hasChanged())

	d, _ = Compare(int8(1), int8(2))
	assert.True(t, d.hasChanged())

	d, _ = Compare(int16(1), int16(1))
	assert.False(t, d.hasChanged())

	d, _ = Compare(int16(1), int16(2))
	assert.True(t, d.hasChanged())

	d, _ = Compare(int32(1), int32(1))
	assert.False(t, d.hasChanged())

	d, _ = Compare(int32(1), int32(2))
	assert.True(t, d.hasChanged())

	d, _ = Compare(int64(1), int64(1))
	assert.False(t, d.hasChanged())

	d, _ = Compare(int64(1), int64(2))
	assert.True(t, d.hasChanged())

	_, err := Compare(1, int32(1))
	assert.EqualError(t, err, NewTypeMismatchError(reflect.Int, reflect.Int32).Error())
}

func TestBasic_Uint(t *testing.T) {
	d, _ := Compare(uint(1), uint(1))
	assert.False(t, d.hasChanged())

	d, _ = Compare(uint(1), uint(2))
	assert.True(t, d.hasChanged())

	d, _ = Compare(uint(1), nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, uint(1))
	assert.True(t, d.hasChanged())

	_, err := Compare(uint(1), 1)
	assert.EqualError(t, err, NewTypeMismatchError(reflect.Uint, reflect.Int).Error())
}

func TestBasic_UintX(t *testing.T) {
	d, _ := Compare(uint8(1), uint8(1))
	assert.False(t, d.hasChanged())

	d, _ = Compare(uint8(1), uint8(2))
	assert.True(t, d.hasChanged())

	d, _ = Compare(uint16(1), uint16(1))
	assert.False(t, d.hasChanged())

	d, _ = Compare(uint16(1), uint16(2))
	assert.True(t, d.hasChanged())

	d, _ = Compare(uint32(1), uint32(1))
	assert.False(t, d.hasChanged())

	d, _ = Compare(uint32(1), uint32(2))
	assert.True(t, d.hasChanged())

	d, _ = Compare(uint64(1), uint64(1))
	assert.False(t, d.hasChanged())

	d, _ = Compare(uint64(1), uint64(2))
	assert.True(t, d.hasChanged())

	_, err := Compare(uint(1), uint32(1))
	assert.EqualError(t, err, NewTypeMismatchError(reflect.Uint, reflect.Uint32).Error())
}

func TestBasic_Float(t *testing.T) {
	d, _ := Compare(float32(24.42), float32(24.42))
	assert.False(t, d.hasChanged())

	d, _ = Compare(float32(24.42), float32(42.24))
	assert.True(t, d.hasChanged())

	d, _ = Compare(float32(24.42), nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, float32(24.42))
	assert.True(t, d.hasChanged())

	d, _ = Compare(float64(24.42), float64(24.42))
	assert.False(t, d.hasChanged())

	d, _ = Compare(float64(24.42), float64(42.24))
	assert.True(t, d.hasChanged())

	d, _ = Compare(float64(24.42), nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, float64(24.42))
	assert.True(t, d.hasChanged())

	_, err := Compare(1.0, float32(1))
	assert.EqualError(t, err, NewTypeMismatchError(reflect.Float64, reflect.Float32).Error())
}

func TestBasic_Complex(t *testing.T) {
	var a128 complex128 = complex(4, 2)
	var b128 complex128 = complex(2, 4)
	var a64 complex64 = complex(4, 2)
	var b64 complex64 = complex(2, 4)

	d, _ := Compare(a64, a64)
	assert.False(t, d.hasChanged())

	d, _ = Compare(a64, b64)
	assert.True(t, d.hasChanged())

	d, _ = Compare(a64, nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, a64)
	assert.True(t, d.hasChanged())

	d, _ = Compare(a128, a128)
	assert.False(t, d.hasChanged())

	d, _ = Compare(a128, b128)
	assert.True(t, d.hasChanged())

	d, _ = Compare(a128, nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, a128)
	assert.True(t, d.hasChanged())

	_, err := Compare(a64, a128)
	assert.EqualError(t, err, NewTypeMismatchError(reflect.Complex64, reflect.Complex128).Error())
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

	var i *int

	d, _ = Compare(i, i)
	assert.False(t, d.hasChanged())

	d, _ = Compare(i, nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, i)
	assert.True(t, d.hasChanged())
}

func TestStruct_Simple(t *testing.T) {
	type SimpleStruct struct {
		Value interface{}
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

func TestStruct_Complex(t *testing.T) {
	type SimpleStruct struct {
		Value interface{}
	}

	type ComplexStruct struct {
		Simple SimpleStruct
	}

	a := ComplexStruct{
		Simple: SimpleStruct{
			Value: 42,
		},
	}

	b := ComplexStruct{
		Simple: SimpleStruct{
			Value: 24,
		},
	}

	d, _ := Compare(a, a)
	assert.False(t, d.hasChanged())

	d, _ = Compare(a, b)
	assert.True(t, d.hasChanged())

	d, _ = Compare(&a, &b)
	assert.True(t, d.hasChanged())
}

func TestStruct_UnexportedField(t *testing.T) {
	type SimpleStruct struct {
		value interface{}
	}

	a := SimpleStruct{
		value: 42,
	}

	b := SimpleStruct{
		value: 24,
	}

	d, _ := Compare(a, a)
	assert.False(t, d.hasChanged())

	d, _ = Compare(a, b)
	assert.False(t, d.hasChanged())

	d, _ = Compare(&a, &b)
	assert.False(t, d.hasChanged())

	d, _ = Compare(&a, nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, &b)
	assert.True(t, d.hasChanged())
}

func TestStruct_FieldName(t *testing.T) {
	type SimpleStruct struct {
		Value interface{} `customTag:"v"`
	}

	a := SimpleStruct{42}
	b := SimpleStruct{24}

	d, _ := Compare(a, b)
	ch := d.Changes()
	assert.Equal(t, "Value", ch[0].Path)

	cmp := NewComparator()
	cmp.TagName = "customTag"

	d, _ = cmp.Compare(a, b)
	ch = d.Changes()
	assert.Equal(t, "v", ch[0].Path)
}

func TestStruct_SkipField(t *testing.T) {
	type SimpleStruct struct {
		Value interface{} `customTag:"-"`
	}

	a := SimpleStruct{42}
	b := SimpleStruct{24}

	cmp := NewComparator()
	cmp.TagName = "customTag"

	d, _ := cmp.Compare(a, b)
	assert.False(t, d.hasChanged())
}

func TestMap(t *testing.T) {
	m1 := map[string]interface{}{"test": 1}
	m2 := map[string]interface{}{"test": 2}
	m3 := map[string]interface{}{"test2": 1}
	m4 := map[string]int{"test": 1}

	d, _ := Compare(m1, m1)
	assert.False(t, d.hasChanged())

	d, _ = Compare(m1, m2)
	assert.True(t, d.hasChanged())

	d, _ = Compare(m3, m1)
	assert.True(t, d.hasChanged())

	d, _ = Compare(m1, nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, m2)
	assert.True(t, d.hasChanged())

	_, err := Compare(m1, m4)
	assert.EqualError(t, err, "Compared values have mismatched types (interface <> int).")
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

func TestSlice_Empty(t *testing.T) {
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

func TestSlice_RespectOrder(t *testing.T) {
	a := []string{"a", "b"}
	b := []string{"b", "a"}

	cmp := NewComparator()
	cmp.RespectSliceOrder = true

	d, _ := cmp.Compare(a, b)
	assert.True(t, d.hasChanged())
}

func TestSlice_SliceId(t *testing.T) {
	type SimpleStruct struct {
		Id interface{} `opt:"id,id"`
	}

	a := []SimpleStruct{
		{Id: 10},
		{Id: 42},
	}

	b := []SimpleStruct{
		{Id: 42},
		{Id: 20},
	}

	expect := []Change{
		{
			Path:        "0.Id",
			StructPath:  "0.Id",
			GenericPath: "*.Id",
			Before:      10,
			After:       nil,
			Action:      DELETE,
		},
		{
			Path:        "2.Id",
			StructPath:  "2.Id",
			GenericPath: "*.Id",
			Before:      nil,
			After:       20,
			Action:      CREATE,
		},
	}

	cmp := NewComparator()

	d, _ := cmp.Compare(a, b)
	ch := d.Changes()

	assert.ElementsMatch(t, expect, ch)

	cmp.TagName = "opt"

	expect[0].Path = "10.id"
	expect[1].Path = "20.id"
	expect[0].StructPath = "10.Id"
	expect[1].StructPath = "20.Id"

	d, _ = cmp.Compare(a, b)
	ch = d.Changes()

	assert.ElementsMatch(t, expect, ch)
}

func TestSlice_SliceIdToNilElem(t *testing.T) {
	type SimpleStruct struct {
		Id interface{} `cmp:"id,id"`
	}

	a := []SimpleStruct{
		{Id: 10},
		{Id: 42},
	}

	expect := []Change{
		{
			Path:        "10.id",
			StructPath:  "10.Id",
			GenericPath: "*.Id",
			Before:      10,
			After:       nil,
			Action:      DELETE,
		},
		{
			Path:        "42.id",
			StructPath:  "42.Id",
			GenericPath: "*.Id",
			Before:      42,
			After:       nil,
			Action:      DELETE,
		},
	}

	d, _ := Compare(a, nil)
	ch := d.Changes()

	assert.ElementsMatch(t, expect, ch)
}

func TestSlice_UnexportedField(t *testing.T) {
	type SimpleStruct struct {
		value interface{} `cmp:",id"`
	}

	a := []SimpleStruct{
		{42},
		{24},
	}

	d, _ := Compare(a, a)
	assert.False(t, d.hasChanged())

	d, _ = Compare(&a, nil)
	assert.True(t, d.hasChanged())

	d, _ = Compare(nil, a)
	assert.True(t, d.hasChanged())
}
