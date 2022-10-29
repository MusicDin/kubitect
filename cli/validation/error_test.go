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
	assert.EqualError(t, err, "error:test.error:Error:Test.Error:1")
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

	es.append(ValidationErrors{
		ValidationError{Err: "test1"},
		ValidationError{Err: "test2"},
	})

	assert.EqualError(t, es[len(es)-3], "test")
	assert.EqualError(t, es[len(es)-2], "test1")
	assert.EqualError(t, es[len(es)-1], "test2")
}

func TestErrors_AppendNil(t *testing.T) {
	prevLen := len(errs)

	errs.append(nil)
	assert.Len(t, errs, prevLen)
}

func TestErrors_SubAppend(t *testing.T) {
	es := errs

	es.subAppend(ValidationError{Err: "test"}, "T", "t")

	es.subAppend(ValidationErrors{
		ValidationError{Err: "test1", Field: "test1", Namespace: "ns"},
		ValidationError{Err: "test2"},
	}, "T", "t")

	e := es[len(es)-3]
	assert.EqualError(t, e, "test")
	assert.Equal(t, "t", e.Field)
	assert.Equal(t, "T", e.StructField)

	e = es[len(es)-2]
	assert.EqualError(t, e, "test1")
	assert.Equal(t, "test1", e.Field)
	assert.Equal(t, "t.ns", e.Namespace)

	e = es[len(es)-1]
	assert.EqualError(t, e, "test2")
	assert.Equal(t, "t", e.Field)
	assert.Equal(t, "T", e.StructField)

}
