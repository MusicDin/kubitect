package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
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

// regex tests whether a pattern matches a string.
// It panics if the provided pattern is incorrect.
func regex(regex string, value string) bool {
	matched, err := regexp.MatchString(regex, value)

	if err != nil {
		panic(err)
	}

	return matched
}

// getDeepValue returns the actual (final) reflect value.
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

// Returns true if the given value is empty.
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	case reflect.Map, reflect.Slice:
		return rv.Len() == 0

	case reflect.Ptr:
		if rv.IsNil() {
			return true
		}
		return isEmpty(rv.Elem().Interface())

	default:
		zero := reflect.Zero(rv.Type()).Interface()
		return reflect.DeepEqual(value, zero)
	}
}
