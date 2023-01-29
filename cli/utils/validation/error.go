package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	ErrorAppendIncompatibleError    = fmt.Errorf("validation.ValidationErrors.Append: Attempt to subAppend incompatible error.")
	ErrorSubAppendIncompatibleError = fmt.Errorf("validation.ValidationErrors.SubAppend: Attempt to subAppend incompatible error.")
)

// Custom validation error.
type ValidationError struct {
	Valid           bool
	Namespace       string
	Field           string
	StructNamespace string
	StructField     string
	Tag             string
	ActualTag       string
	Kind            reflect.Kind
	Type            reflect.Type
	Param           string
	Value           interface{}
	Err             string
	RealErr         string
}

// Error returns validation error as a string. It also populates the variables
// in the error message.
func (e ValidationError) Error() string {
	err := e.Err
	re := getDeepValue(reflect.ValueOf(e))

	if len(err) == 0 {
		return err
	}

	for i := 0; i < re.NumField(); i++ {
		fName := re.Type().Field(i).Name
		fValue := re.Field(i).Interface()

		old := fmt.Sprintf("{.%s}", fName)
		new := fmt.Sprintf("%v", fValue)

		err = strings.ReplaceAll(err, old, new)
	}

	err = strings.ReplaceAll(err, " '' ", " ")

	return err
}

// prependPath populates Namespace, Field, StructNamespace and StructField
// values of the validation error.
func (e *ValidationError) prependPath(realKey, key interface{}) {
	sep := "."

	if len(e.Namespace) == 0 {
		sep = ""
	}

	if len(e.Field) == 0 {
		e.Field = fmt.Sprintf("%v", key)
		e.StructField = fmt.Sprintf("%v", realKey)
	}

	e.Namespace = fmt.Sprintf("%v%s%s", key, sep, e.Namespace)
	e.StructNamespace = fmt.Sprintf("%v%s%s", realKey, sep, e.StructNamespace)
}

type ValidationErrors []ValidationError

// Error returns all validation errors as a string.
func (es ValidationErrors) Error() string {
	var out []string

	for _, e := range es {
		out = append(out, e.Error())
	}

	return strings.Join(out, "\n")
}

// append appends validation error(s).
//
// It panics if error is not of type ValidationError or ValidationErrors.
func (es *ValidationErrors) append(err error) {

	if err == nil {
		return
	}

	switch err.(type) {
	case ValidationErrors:
		*es = append(*es, err.(ValidationErrors)...)
	case ValidationError:
		*es = append(*es, err.(ValidationError))
	default:
		panic(ErrorAppendIncompatibleError)
	}
}

// subAppend appends validation error(s) and corrects their paths (Namespace
// and StructNamespace) accordingly.
//
// It panics if error is not of type ValidationError or ValidationErrors.
func (es *ValidationErrors) subAppend(err error, realKey, key interface{}) {

	if err == nil {
		return
	}

	if errs, ok := err.(ValidationErrors); ok {
		for i := range errs {
			errs[i].prependPath(realKey, key)
		}

		es.append(errs)
		return
	}

	if e, ok := err.(ValidationError); ok {
		e.prependPath(realKey, key)
		es.append(e)
		return
	}

	panic(ErrorSubAppendIncompatibleError)
}

// toValidationErrors converts validator.ValidationErrors to validation.ValidationErrors.
func toValidationErrors(err error) ValidationErrors {

	if err == nil {
		return nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return ValidationErrors{{Err: err.Error()}}
	}

	es := make(ValidationErrors, 0)

	for _, fe := range err.(validator.ValidationErrors) {
		es.append(ValidationError{
			Valid:           true,
			Namespace:       fe.Namespace(),
			Field:           fe.Field(),
			StructNamespace: fe.StructNamespace(),
			StructField:     fe.StructField(),
			Tag:             fe.Tag(),
			ActualTag:       fe.ActualTag(),
			Kind:            fe.Kind(),
			Type:            fe.Type(),
			Param:           fe.Param(),
			Value:           fe.Value(),
			Err:             fe.Error(),
		})
	}

	return es
}
