package cmp

import (
	"fmt"
	"reflect"
	"unsafe"
)

type flag uintptr

const (
	flagStickyRO flag = 1 << 5
	flagEmbedRO  flag = 1 << 6
	flagRO       flag = flagStickyRO | flagEmbedRO
)

// exportInterface returns an interface of the reflect value.
func exportInterface(v reflect.Value) interface{} {
	if !v.CanInterface() {
		vPtr := unsafe.Pointer(&v)
		offset := unsafe.Sizeof(uintptr(0)) << 1
		flagTmp := (*flag)(unsafe.Pointer(uintptr(vPtr) + uintptr(offset)))
		*flagTmp &^= flagRO
	}

	return v.Interface()
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
