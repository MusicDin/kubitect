package cmp

import (
	"fmt"
	"reflect"
)

func (c *Comparator) cmpInt(parent *DiffNode, key interface{}, a, b reflect.Value) error {
	if a.Kind() == reflect.Invalid {
		parent.addLeaf(CREATE, key, nil, exportInterface(b))
		return nil
	}

	if b.Kind() == reflect.Invalid {
		parent.addLeaf(DELETE, key, exportInterface(a), nil)
		return nil
	}

	if a.Kind() != b.Kind() {
		return fmt.Errorf("Type mismatch: %v <> %v\n", a.Kind(), b.Kind())
	}

	if a.Int() != b.Int() {
		parent.addLeaf(MODIFY, key, a.Int(), b.Int())
		return nil
	}

	parent.addLeaf(NONE, key, a.Int(), b.Int())

	return nil
}
