package cmp

import (
	"reflect"
)

func (c *Comparator) cmpInterface(parent *DiffNode, key interface{}, a, b reflect.Value) error {
	if a.Kind() == reflect.Invalid {
		parent.addLeaf(CREATE, key, nil, exportInterface(b))
	}

	if b.Kind() == reflect.Invalid {
		parent.addLeaf(DELETE, key, exportInterface(a), nil)
		return nil
	}

	if a.Kind() != b.Kind() {
		return NewTypeMismatchError(a.Kind(), b.Kind())
	}

	if a.IsNil() && b.IsNil() {
		return nil
	}

	return c.compare(parent, key, a.Elem(), b.Elem())
}
