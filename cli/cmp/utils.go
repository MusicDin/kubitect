package cmp

import (
	"fmt"
	"reflect"
	"strings"
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

// toSliceKey wraps index into square brackets.
func toSliceKey(key interface{}) string {
	return "[" + toString(key) + "]"
}

// fromSliceKey unwraps index out of square brackets.
func fromSliceKey(k interface{}) string {
	if !isSliceKey(k) {
		return ""
	}

	key := toString(k)
	key = strings.TrimPrefix(key, "[")
	key = strings.TrimSuffix(key, "]")

	return key
}

// isSliceKey checks whether given key represents a slice key,
// which means that it starts with "[" and ends with "]".
func isSliceKey(k interface{}) bool {
	s := toString(k)
	return strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]")
}

// exportInterface returns an interface of the reflect value.
func exportInterface(v reflect.Value) interface{} {
	if !v.CanInterface() {
		ptr := unsafe.Pointer(v.UnsafeAddr())
		typ := v.Type()
		return reflect.NewAt(typ, ptr).Elem().Interface()
	}
	return v.Interface()
}

// getDeepValue recursively returns the actual value that a
// reflect value contains.
func getDeepValue(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Interface:
		return getDeepValue(v.Elem())
	case reflect.Ptr:
		return getDeepValue(reflect.Indirect(v))
	default:
		return v
	}
}
