package cmp

import (
	"reflect"
)

func (c *Comparator) cmpMap(a, b reflect.Value) (*DiffNode, error) {
	if !a.IsValid() {
		return c.cmpMapToNil(Create, b)
	}

	if !b.IsValid() {
		return c.cmpMapToNil(Delete, a)
	}

	// Create new pairs with change type None, since we know that none
	// of the values is nil.
	pairs := NewPairs(b.Type())
	pairs.changeType = None

	for _, k := range a.MapKeys() {
		if a.IsValid() {
			v := a.MapIndex(k)
			pairs.addA(k, v)
		}
	}

	for _, k := range b.MapKeys() {
		if b.IsValid() {
			v := b.MapIndex(k)
			pairs.addB(k, v)
		}
	}

	return c.cmpPairs(pairs)
}

// cmpMapToNil recursively adds all elements of the map to the diff tree
// by comparing them to a nil value.
func (c *Comparator) cmpMapToNil(t ChangeType, v reflect.Value) (*DiffNode, error) {
	pairs := NewPairs(v.Type())

	for _, k := range v.MapKeys() {
		v := v.MapIndex(k)
		switch t {
		case Create:
			pairs.addB(k, v)
		case Delete:
			pairs.addA(k, v)
		}
	}

	node, err := c.cmpPairs(pairs)
	if err != nil {
		return nil, err
	}

	if node != nil {
		node.setChangeTypeOfChildren(t)
	}

	return node, nil
}
