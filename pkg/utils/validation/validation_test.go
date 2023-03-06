package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TEST DATA

type SimpleStruct struct {
	Field string `yaml:"wrong" opt:"name"`
}

func (t *SimpleStruct) Validate() error {
	return Struct(&t,
		Field(&t.Field, Min(42)),
	)
}

type ComplexStruct struct {
	List []Elem          `yaml:"list"`
	Map  map[string]Elem `yaml:"map"`
}

func (p ComplexStruct) Validate() error {
	return Struct(&p,
		Field(&p.List),
		Field(&p.Map),
	)
}

type Elem struct {
	Id string `opt:",id"`
	A  *SimpleStruct
}

func (ie Elem) Validate() error {
	return Struct(&ie,
		Field(&ie.A, Required()),
	)
}

// TESTS

func TestErrorFieldIsNotPointer(t *testing.T) {
	test := SimpleStruct{}
	fn := func() {
		Struct(&test, Field(test.Field))
	}

	assert.PanicsWithError(t, ErrorFieldIsNotPointer.Error(), fn)
}

func TestErrorStructFieldNotFound(t *testing.T) {
	test := SimpleStruct{}
	fn := func() {
		Struct(&test, Field(&test))
	}

	assert.PanicsWithError(t, ErrorStructFieldNotFound.Error(), fn)
}

func TestErrorMissingStructPointer(t *testing.T) {
	fn := func() {
		Struct(nil)
	}

	assert.PanicsWithError(t, ErrorMissingStructPointer.Error(), fn)
}

func TestErrorInvalidStructPointer(t *testing.T) {
	fn := func() {
		Struct(1)
	}

	assert.PanicsWithError(t, ErrorInvalidStructPointer.Error(), fn)
}

func TestSetTag(t *testing.T) {
	test := SimpleStruct{}

	ve := test.Validate().(ValidationErrors)
	assert.Equal(t, "wrong", ve[0].Field)

	SetTagName("opt")

	ve = test.Validate().(ValidationErrors)
	assert.Equal(t, "name", ve[0].Field)
}

func TestListId(t *testing.T) {

	o := &ComplexStruct{
		List: []Elem{
			{Id: "e1"},
			{Id: "e1", A: &SimpleStruct{}},
		},
	}

	err := o.Validate().(ValidationErrors)

	assert.Equal(t, "list.e1.A", err[0].Namespace)
	assert.Equal(t, "List.0.A", err[0].StructNamespace)
}

func TestMap(t *testing.T) {

	o := &ComplexStruct{
		Map: map[string]Elem{
			"e1": {},
		},
	}

	err := o.Validate().(ValidationErrors)
	assert.Equal(t, "map.e1.A", err[0].Namespace)
	assert.Equal(t, "Map.e1.A", err[0].StructNamespace)
}
