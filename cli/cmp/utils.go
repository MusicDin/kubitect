package cmp

import (
	"fmt"
	"reflect"
	"unsafe"
)

// toString converts an interface to a string.
func toString(v interface{}) string {
	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Pointer:
		return toString(reflect.Indirect(rv))
	default:
		return fmt.Sprint(v)
	}
}

// exportInterface returns an interface of the reflect value.
func exportInterface(v reflect.Value) interface{} {
	if !v.CanAddr() {
		return nil
	}

	if v.CanInterface() {
		return v.Interface()
	}

	ptr := unsafe.Pointer(v.UnsafeAddr())
	copy := reflect.NewAt(v.Type(), ptr)

	return copy.Elem().Interface()
}

// getDeepValue recursively returns the actual value that a
// reflect value contains.
func getDeepValue(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Interface:
		return getDeepValue(v.Elem())
	case reflect.Pointer:
		return getDeepValue(reflect.Indirect(v))
	default:
		return v
	}
}
