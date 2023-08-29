package cmp

import (
	"fmt"
	"reflect"
)

// TypeMismatchError represents a type mismatch error when comparing two
// values of different types.
type TypeMismatchError struct {
	aKind reflect.Kind
	bKind reflect.Kind
}

func NewTypeMismatchError(aKind, bKind reflect.Kind) *TypeMismatchError {
	return &TypeMismatchError{
		aKind: aKind,
		bKind: bKind,
	}
}

func (e *TypeMismatchError) Error() string {
	return fmt.Sprintf("Compared values have either unsupported or mismatched types (%v <> %v).", e.aKind, e.bKind)
}
