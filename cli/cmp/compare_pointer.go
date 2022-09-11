package cmp

import (
	"fmt"
	"reflect"
)

func (c *Comparator) cmpPointer(parent *DiffNode, key interface{}, a, b reflect.Value) error {
	if a.Kind() == b.Kind() {
		if a.IsNil() && b.IsNil() {
			parent.AddLeaf(NONE, key, nil, nil)
			return nil
		}

		return c.compare(parent, key, reflect.Indirect(a), reflect.Indirect(b))
	}

	if a.Kind() == reflect.Invalid {
		if !b.IsNil() {
			return c.compare(parent, key, reflect.ValueOf(nil), reflect.Indirect(b))
		}

		parent.AddLeaf(MODIFY, key, nil, exportInterface(b))
		return nil
	}

	if b.Kind() == reflect.Invalid {
		if !a.IsNil() {
			return c.compare(parent, key, reflect.Indirect(a), reflect.ValueOf(nil))
		}

		parent.AddLeaf(DELETE, key, exportInterface(a), nil)
		return nil
	}

	return fmt.Errorf("Type mismatch: %v <> %v", a.Kind(), b.Kind())
}
