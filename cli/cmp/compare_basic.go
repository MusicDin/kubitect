package cmp

import "reflect"

func (c *Comparator) cmpBasic(a, b reflect.Value) (*DiffNode, error) {
	if a.Kind() == reflect.Invalid {
		return NewLeaf(CREATE, nil, b.Interface()), nil
	}

	if b.Kind() == reflect.Invalid {
		return NewLeaf(DELETE, a.Interface(), nil), nil
	}

	if a.Kind() != b.Kind() {
		return nil, NewTypeMismatchError(a.Kind(), b.Kind())
	}

	ai := a.Interface()
	bi := b.Interface()

	if ai != bi {
		return NewLeaf(MODIFY, ai, bi), nil
	}

	return NewLeaf(NONE, ai, bi), nil
}
