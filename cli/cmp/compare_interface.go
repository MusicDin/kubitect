package cmp

import (
	"fmt"
	"reflect"
)

func (c *Comparator) cmpInterface(parent *DiffNode, key interface{}, a, b reflect.Value) error {
	if a.Kind() == reflect.Invalid {
		parent.AddLeaf(CREATE, key, nil, exportInterface(b))
	}

	if b.Kind() == reflect.Invalid {
		parent.AddLeaf(DELETE, key, exportInterface(a), nil)
		return nil
	}

	if a.Kind() != b.Kind() {
		return fmt.Errorf("Type mismatch: %v <> %v\n", a.Kind(), b.Kind())
	}

	if a.IsNil() && b.IsNil() {
		return nil
	}

	return c.compare(parent, key, a.Elem(), b.Elem())
}
