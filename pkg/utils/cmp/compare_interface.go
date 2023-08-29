package cmp

import "reflect"

func (c *Comparator) cmpInterface(a, b reflect.Value) (*DiffNode, error) {
	if !a.IsValid() {
		return c.newLeaf(Create, nil, b.Interface()), nil
	}

	if !b.IsValid() {
		return c.newLeaf(Delete, a.Interface(), nil), nil
	}

	return c.compare(a.Elem(), b.Elem())
}
