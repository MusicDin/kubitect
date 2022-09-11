package cmp

import (
	"fmt"
	"reflect"
)

func (c *Comparator) cmpString(parent *DiffNode, key interface{}, a, b reflect.Value) error {
	if a.Kind() == reflect.Invalid {
		parent.AddLeaf(CREATE, key, nil, exportInterface(b))
		return nil
	}

	if b.Kind() == reflect.Invalid {
		parent.AddLeaf(DELETE, key, exportInterface(a), nil)
		return nil
	}

	if a.Kind() != b.Kind() {
		return fmt.Errorf("Type mismatch: %v <> %v\n", a.Kind(), b.Kind())
	}

	if a.String() != b.String() {
		parent.AddLeaf(MODIFY, key, a.String(), b.String())
		return nil
	}

	parent.AddLeaf(NONE, key, a.String(), b.String())

	return nil
}
