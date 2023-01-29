package template

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

func BuiltInFuncs() template.FuncMap {
	return template.FuncMap{
		"list":     fList,
		"append":   fAppend,
		"prepend":  fPrepend,
		"map":      fMap,
		"join":     fJoin,
		"contains": fContains,
		"deref":    fDeref,
	}
}

// fList creates a list with given elements.
// If no elements are provided, an empty list is returned.
func fList(elements ...interface{}) []interface{} {
	if len(elements) == 0 {
		return make([]interface{}, 0)
	}
	return elements
}

// fAppend returns the list with a given element appended.
func fAppend(list []interface{}, element interface{}) []interface{} {
	return append(list, element)
}

// fPrepend returns the list with a given element appended.
func fPrepend(list []interface{}, element interface{}) []interface{} {
	return append([]interface{}{element}, list...)
}

// fMap maps objects from a list to a new list that contains only object
// fields that match a given name.
func fMap(fieldName string, list interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(list)

	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("extractField: list is not a slice")
	}

	fields := make([]interface{}, 0)

	for i := 0; i < v.Len(); i++ {
		el := v.Index(i)

		if el.Kind() != reflect.Struct {
			return nil, fmt.Errorf("extractField: list element is not a struct")
		}

		f := el.FieldByName(fieldName)

		if !f.IsValid() {
			return nil, fmt.Errorf("extractField: field %s not found in struct", fieldName)
		}

		fields = append(fields, f.Interface())
	}

	return fields, nil
}

// fMap joins a list of objects with a given separator.
func fJoin(sep string, list []interface{}) string {
	sList := make([]string, len(list))

	for i, e := range list {
		sList[i] = fmt.Sprint(e)
	}

	return strings.Join(sList, sep)
}

// fContains checks whether a list contains a given value.
func fContains(v interface{}, list []interface{}) bool {
	for _, e := range list {
		if e == v {
			return true
		}
	}

	return false
}

// fDeref dereferences a given value if possible. Otherwise
// it returns the initial value.
func fDeref(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	v := reflect.ValueOf(value)

	if v.Kind() != reflect.Pointer {
		return value
	}

	pv := v.Elem()

	if pv.Kind() == reflect.Invalid || !pv.CanInterface() {
		return nil
	}

	return pv.Interface()
}
