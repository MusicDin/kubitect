package cmp

import (
	"reflect"
)

func (c *Comparator) cmpInterface(a, b reflect.Value) (*DiffNode, error) {
	if a.Kind() == reflect.Invalid {
		return NewLeaf(CREATE, nil, exportInterface(b)), nil
	}

	if b.Kind() == reflect.Invalid {
		return NewLeaf(DELETE, exportInterface(a), nil), nil
	}

	return c.compare(a.Elem(), b.Elem())
}
