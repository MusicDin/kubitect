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

var err_error = "error:test.error:Error:Test.Error:1"

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

func TestErrorAppendIncompatibleError(t *testing.T) {
	fn := func() {
		(&ValidationErrors{}).append(fmt.Errorf(""))
	}
	assert.PanicsWithError(t, ErrorAppendIncompatibleError.Error(), fn)
}

func TestErrorSubAppendIncompatibleError(t *testing.T) {
	fn := func() {
		(&ValidationErrors{}).subAppend(fmt.Errorf(""), "", "")
	}
	assert.PanicsWithError(t, ErrorSubAppendIncompatibleError.Error(), fn)
}

func TestEmptyError(t *testing.T) {
	assert.Empty(t, ValidationError{}.Error())
	assert.Empty(t, ValidationErrors{}.Error())
}

func TestSampleError(t *testing.T) {
	result := fmt.Sprintf("%s\n", strings.Join(errs_errors, "\n"))
	assert.Equal(t, errs.Error(), result)
}

func TestErrorPopulation(t *testing.T) {
	assert.Equal(t, err.Error(), err_error)
}
func TestErrorsPopulation(t *testing.T) {
	for i := range errs {
		assert.Equal(t, errs[i].Error(), errs_errors[i])
	}
}

func TestPrependPath(t *testing.T) {
	e := err
	e.prependPath("Test", "test")

	assert.Equal(t, e.Field, "error")
	assert.Equal(t, e.StructField, "Error")
	assert.Equal(t, e.Namespace, "test.test.error")
	assert.Equal(t, e.StructNamespace, "Test.Test.Error")
}

func TestErrorsAppend(t *testing.T) {
	es := errs
	es.append(ValidationError{Err: "test"})

	assert.Equal(t, es[len(es)-1].Error(), "test")

	es.append(ValidationErrors{
		ValidationError{Err: "test1"},
		ValidationError{Err: "test2"},
	})

	assert.Equal(t, es[len(es)-2].Error(), "test1")
	assert.Equal(t, es[len(es)-1].Error(), "test2")

	prevLen := len(es)
	es.append(nil)
	assert.Equal(t, prevLen, len(es))
}

func TestErrorsSubAppend(t *testing.T) {
	es := errs
	es.subAppend(ValidationError{Err: "test"}, "T", "t")

	e := es[len(es)-1]
	assert.Equal(t, e.Error(), "test")
	assert.Equal(t, e.Field, "t")
	assert.Equal(t, e.StructField, "T")

	es.subAppend(ValidationErrors{
		ValidationError{Err: "test1", Field: "test1", Namespace: "ns"},
		ValidationError{Err: "test2"},
	}, "T", "t")

	e = es[len(es)-1]
	assert.Equal(t, e.Error(), "test2")
	assert.Equal(t, e.Field, "t")
	assert.Equal(t, e.StructField, "T")

	e = es[len(es)-2]
	assert.Equal(t, e.Error(), "test1")
	assert.Equal(t, e.Field, "test1")
	assert.Equal(t, e.Namespace, "t.ns")
}
