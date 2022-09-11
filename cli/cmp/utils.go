package cmp

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

// toString converts an interface to a string.
func toString(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	default:
		return fmt.Sprint(v)
	}
}

const isExportFlag uintptr = (1 << 5) | (1 << 6)

// exportInterface returns an interface of a reflect value.
func exportInterface(v reflect.Value) interface{} {
	if !v.CanInterface() {
		flagTmp := (*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + 2*unsafe.Sizeof(uintptr(0))))
		*flagTmp = (*flagTmp) & (^isExportFlag)
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
