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

func TestF_First(t *testing.T) {
	assert.Equal(t, nil, fFirst(nil))
	assert.Equal(t, nil, fFirst([]interface{}{}))
	assert.Equal(t, "42", fFirst([]interface{}{"42"}))
	assert.Equal(t, "42", fFirst([]interface{}{"42", "50", 1}))
}

type (
	TestObject interface {
		GetName() string
		GetValue() interface{}
	}

	testObject struct {
		Name  string
		Value interface{}
	}
)

func (o testObject) GetName() string {
	return o.Name
}

func (o testObject) GetValue() interface{} {
	return o.Name
}

func TestF_Map(t *testing.T) {
	objs := []testObject{
		{
			Name:  "obj1",
			Value: "val1",
		},
		{
			Name:  "obj2",
			Value: "val2",
		},
	}

	res, err := fMap("Name", []testObject{})
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{}, res)

	res, err = fMap("Value", []testObject{})
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
	assert.EqualError(t, err, "map: list must be either of type slice or array (actual: string)")

	_, err = fMap("X", []string{"invalid"})
	assert.EqualError(t, err, "map: list elements must be of type struct (actual: string)")

	_, err = fMap("X", []testObject{{Name: "test"}})
	assert.EqualError(t, err, "map: field X not found in a struct")
}

func TestF_Select(t *testing.T) {
	objs := []struct {
		Name  string
		Value interface{}
	}{
		{
			Name:  "obj1",
			Value: "val1",
		},
		{
			Name:  "obj2",
			Value: "val2",
		},
	}

	res, err := fSelect("Invalid", nil, objs)
	assert.NoError(t, err)
	assert.Empty(t, res)

	res, err = fSelect("Name", nil, objs)
	assert.NoError(t, err)
	assert.ElementsMatch(t, objs, res)

	res, err = fSelect("Name", "obj1", objs)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, objs[0], res[0])
}

func TestF_Select_Interface(t *testing.T) {
	objs := []TestObject{
		testObject{
			Name:  "obj1",
			Value: "val1",
		},
		testObject{
			Name:  "obj2",
			Value: "val2",
		},
	}

	res, err := fSelect("Invalid", nil, objs)
	assert.NoError(t, err)
	assert.Empty(t, res)

	res, err = fSelect("Name", nil, objs)
	assert.NoError(t, err)
	assert.ElementsMatch(t, objs, res)

	res, err = fSelect("Name", "obj1", objs)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, objs[0], res[0])
}

func TestF_Select_Error(t *testing.T) {
	_, err := fSelect("", nil, "invalid")
	assert.EqualError(t, err, "select: list must be either of type slice or array (actual: string)")

	_, err = fSelect("", nil, []string{"invalid"})
	assert.EqualError(t, err, "select: list elements must be of type struct (actual: string)")
}

func TestF_Join(t *testing.T) {
	assert.Equal(t, "", fJoin(", ", nil))
	assert.Equal(t, "test42", fJoin("", []interface{}{"test", "42"}))
	assert.Equal(t, "test, 123", fJoin(", ", []interface{}{"test", 123}))
}

func TestF_Contains(t *testing.T) {
	objs := []interface{}{
		testObject{
			Name:  "obj1",
			Value: "val1",
		},
		testObject{
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
