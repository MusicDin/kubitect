package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestF_List(t *testing.T) {
	assert.Equal(t, []interface{}{}, fList())
	assert.Equal(t, []interface{}{"heii"}, fList("heii"))
	assert.Equal(t, []interface{}{"heii", 123}, fList("heii", 123))
}

func TestF_Append(t *testing.T) {
	assert.Equal(t, []interface{}{"42"}, fAppend(nil, "42"))
	assert.Equal(t, []interface{}{"test", "42"}, fAppend([]interface{}{"test"}, "42"))
}

func TestF_Prepend(t *testing.T) {
	assert.Equal(t, []interface{}{"42"}, fPrepend(nil, "42"))
	assert.Equal(t, []interface{}{"42", "test"}, fPrepend([]interface{}{"test"}, "42"))
}

type TestObject struct {
	Name  string
	Value interface{}
}

func TestF_Map(t *testing.T) {
	objs := []TestObject{
		{
			Name:  "obj1",
			Value: "val1",
		},
		{
			Name:  "obj2",
			Value: "val2",
		},
	}

	res, err := fMap("Name", []TestObject{})
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{}, res)

	res, err = fMap("Value", []TestObject{})
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{}, res)

	res, err = fMap("Name", objs)
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{"obj1", "obj2"}, res)

	res, err = fMap("Value", objs)
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{"val1", "val2"}, res)
}

func TestF_Map_Error(t *testing.T) {
	_, err := fMap("invalid", "X")
	assert.EqualError(t, err, "extractField: list is not a slice")

	_, err = fMap("X", []string{"invalid"})
	assert.EqualError(t, err, "extractField: list element is not a struct")

	_, err = fMap("X", []TestObject{{Name: "test"}})
	assert.EqualError(t, err, "extractField: field X not found in struct")
}

func TestF_Join(t *testing.T) {
	assert.Equal(t, "", fJoin(", ", nil))
	assert.Equal(t, "test42", fJoin("", []interface{}{"test", "42"}))
	assert.Equal(t, "test, 123", fJoin(", ", []interface{}{"test", 123}))
}

func TestF_Contains(t *testing.T) {
	objs := []interface{}{
		TestObject{
			Name:  "obj1",
			Value: "val1",
		},
		TestObject{
			Name:  "obj2",
			Value: "val2",
		},
	}

	assert.True(t, fContains(objs[1], objs))
	assert.True(t, fContains("test", []interface{}{12, "test", 123}))
	assert.False(t, fContains(false, []interface{}{12, "test", 123}))
	assert.False(t, fContains("test", nil))
}

func TestF_Deref(t *testing.T) {
	var str string = "test"
	var nilPtr *string = nil

	assert.Equal(t, nil, fDeref(nil))
	assert.Equal(t, "test", fDeref(str))
	assert.Equal(t, "test", fDeref(&str))
	assert.Equal(t, nil, fDeref(nilPtr))
}
