package cmp

import "reflect"

func (c *Comparator) cmpPointer(a, b reflect.Value) (*DiffNode, error) {
	if !a.IsValid() {
		if !b.IsNil() {
			return c.compare(reflect.ValueOf(nil), reflect.Indirect(b))
		}

		return c.newLeaf(Create, nil, b.Interface()), nil
	}

	if !b.IsValid() {
		if !a.IsNil() {
			return c.compare(reflect.Indirect(a), reflect.ValueOf(nil))
		}

		return c.newLeaf(Delete, a.Interface(), nil), nil
	}

	if a.IsNil() && b.IsNil() {
		return c.newLeaf(None, nil, nil), nil
	}

	return c.compare(reflect.Indirect(a), reflect.Indirect(b))
}
