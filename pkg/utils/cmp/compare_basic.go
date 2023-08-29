package cmp

import "reflect"

func (c *Comparator) cmpBasic(a, b reflect.Value) (*DiffNode, error) {
	if !a.IsValid() {
		return c.newLeaf(Create, nil, b.Interface()), nil
	}

	if !b.IsValid() {
		return c.newLeaf(Delete, a.Interface(), nil), nil
	}

	if a.Kind() != b.Kind() {
		return nil, NewTypeMismatchError(a.Kind(), b.Kind())
	}

	ai := a.Interface()
	bi := b.Interface()

	if ai != bi {
		return c.newLeaf(Modify, ai, bi), nil
	}

	return c.newLeaf(None, ai, bi), nil
}
