package cmp

import "reflect"

func (c *Comparator) cmpBasic(a, b reflect.Value) (*DiffNode, error) {
	if a.Kind() == reflect.Invalid {
		return c.newLeaf(CREATE, nil, b.Interface()), nil
	}

	if b.Kind() == reflect.Invalid {
		return c.newLeaf(DELETE, a.Interface(), nil), nil
	}

	if a.Kind() != b.Kind() {
		return nil, NewTypeMismatchError(a.Kind(), b.Kind())
	}

	ai := a.Interface()
	bi := b.Interface()

	if ai != bi {
		return c.newLeaf(MODIFY, ai, bi), nil
	}

	return c.newLeaf(NONE, ai, bi), nil
}
