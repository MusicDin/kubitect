package cmp

import "reflect"

func (c *Comparator) cmpBasic(a, b reflect.Value) (*DiffNode, error) {
	if a.Kind() == reflect.Invalid {
		return NewLeaf(CREATE, nil, exportInterface(b)), nil
	}

	if b.Kind() == reflect.Invalid {
		return NewLeaf(DELETE, exportInterface(a), nil), nil
	}

	if a.Kind() != b.Kind() {
		return nil, NewTypeMismatchError(a.Kind(), b.Kind())
	}

	ai := exportInterface(a)
	bi := exportInterface(b)

	if ai != bi {
		return NewLeaf(MODIFY, ai, bi), nil
	}

	return NewLeaf(NONE, ai, bi), nil
}
