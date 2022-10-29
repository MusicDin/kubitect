package validation

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var err = ValidationError{
	Namespace:       "test.error",
	Field:           "error",
	StructNamespace: "Test.Error",
	StructField:     "Error",
	Value:           "1",
	Err:             "{.Field}:{.Namespace}:{.StructField}:{.StructNamespace}:{.Value}",
}

var errs = ValidationErrors{
	{
		Namespace:       "test.error",
		Field:           "error",
		StructNamespace: "Test.Error",
		StructField:     "Error",
		Value:           "1",
		Err:             "{.Field}:{.Namespace}:{.StructField}:{.StructNamespace}:{.Value}",
	},
	{
		Namespace:       "Error",
		Field:           "Error",
		StructNamespace: "Error",
		StructField:     "Error",
		Value:           "2",
		Err:             "{.Field}:{.Namespace}:{.StructField}:{.StructNamespace}:{.Value}",
	},
	{
		Err: "{.Field}:{.Namespace}:{.StructField}:{.StructNamespace}:{.Value}",
	},
}

var errs_errors = []string{
	"error:test.error:Error:Test.Error:1",
	"Error:Error:Error:Error:2",
	"::::<nil>",
}

func TestAppendIncompatibleError(t *testing.T) {
	fn := func() {
		(&ValidationErrors{}).append(fmt.Errorf(""))
	}
	assert.PanicsWithError(t, ErrorAppendIncompatibleError.Error(), fn)
}

func TestSubAppendIncompatibleError(t *testing.T) {
	fn := func() {
		(&ValidationErrors{}).subAppend(fmt.Errorf(""), "", "")
	}
	assert.PanicsWithError(t, ErrorSubAppendIncompatibleError.Error(), fn)
}

func TestError_Empty(t *testing.T) {
	assert.Empty(t, ValidationError{}.Error())
}

func TestError_Population(t *testing.T) {
	assert.Equal(t, "error:test.error:Error:Test.Error:1", err.Error())
}

func TestError_PrependPath(t *testing.T) {
	e := err
	e.prependPath("Test", "test")

	assert.Equal(t, "error", e.Field)
	assert.Equal(t, "Error", e.StructField)
	assert.Equal(t, "test.test.error", e.Namespace)
	assert.Equal(t, "Test.Test.Error", e.StructNamespace)
}

func TestErrors_Empty(t *testing.T) {
	assert.Empty(t, ValidationErrors{}.Error())
}

func TestErrors_Population(t *testing.T) {
	assert.Equal(t, strings.Join(errs_errors, "\n"), errs.Error())
}

func TestErrors_Append(t *testing.T) {
	es := errs
	es.append(ValidationError{Err: "test"})

	assert.Equal(t, "test", es[len(es)-1].Error())

	es.append(ValidationErrors{
		ValidationError{Err: "test1"},
		ValidationError{Err: "test2"},
	})

	assert.Equal(t, "test1", es[len(es)-2].Error())
	assert.Equal(t, "test2", es[len(es)-1].Error())

	prevLen := len(es)
	es.append(nil)
	assert.Equal(t, prevLen, len(es))
}

func TestErrors_SubAppend(t *testing.T) {
	es := errs
	es.subAppend(ValidationError{Err: "test"}, "T", "t")

	e := es[len(es)-1]
	assert.Equal(t, "test", e.Error())
	assert.Equal(t, "t", e.Field)
	assert.Equal(t, "T", e.StructField)

	es.subAppend(ValidationErrors{
		ValidationError{Err: "test1", Field: "test1", Namespace: "ns"},
		ValidationError{Err: "test2"},
	}, "T", "t")

	e = es[len(es)-1]
	assert.Equal(t, "test2", e.Error())
	assert.Equal(t, "t", e.Field)
	assert.Equal(t, "T", e.StructField)

	e = es[len(es)-2]
	assert.Equal(t, "test1", e.Error())
	assert.Equal(t, "test1", e.Field)
	assert.Equal(t, "t.ns", e.Namespace)
}
