package defaults

import (
	def "github.com/creasty/defaults"
)

// Set sets default values for the structure on a
// given pointer and returns a potential error.
func Set(ptr interface{}) error {
	return def.Set(ptr)
}

// Assign assigns default values for the structure
// on a given pointer and returns back the pointer to
// allow method chaining.
//
// It panics if an error happens.
func Assign[T any](ptr T) T {
	def.MustSet(ptr)
	return ptr
}

// CanUpdate returns true when the given value is an
// initial value of its type.
func CanUpdate(v interface{}) bool {
	return def.CanUpdate(v)
}

// Default returns default value if the given value
// can be updated (is an initial value of its type).
// Otherwise the value itself is returned.
func Default[T any](v T, def T) T {
	if CanUpdate(v) {
		return def
	}
	return v
}
