package cmp

import "reflect"

func (c *Comparator) cmpPointer(a, b reflect.Value) (*DiffNode, error) {
	if a.Kind() == reflect.Invalid {
		if !b.IsNil() {
			return c.compare(reflect.ValueOf(nil), reflect.Indirect(b))
		}

		return NewLeaf(MODIFY, nil, b.Interface()), nil
	}

	if b.Kind() == reflect.Invalid {
		if !a.IsNil() {
			return c.compare(reflect.Indirect(a), reflect.ValueOf(nil))
		}

		return NewLeaf(DELETE, a.Interface(), nil), nil
	}

	if a.IsNil() && b.IsNil() {
		return NewLeaf(NONE, nil, nil), nil
	}

	return c.compare(reflect.Indirect(a), reflect.Indirect(b))
}
