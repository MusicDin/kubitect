package cmp

import (
	"fmt"
	"reflect"
)

func (c *Comparator) cmpMap(a, b reflect.Value) (*DiffNode, error) {
	if a.Kind() == reflect.Invalid {
		return c.addPlainMap(CREATE, b)
	}

	if b.Kind() == reflect.Invalid {
		return c.addPlainMap(DELETE, a)
	}

	pairs := NewPairs(a.Type())

	for _, k := range a.MapKeys() {
		ai := a.MapIndex(k)
		if ai.Kind() != reflect.Invalid {
			pairs.addA(k, &ai)
		}
	}

	for _, k := range b.MapKeys() {
		bi := b.MapIndex(k)
		if bi.Kind() != reflect.Invalid {
			pairs.addB(k, &bi)
		}
	}

	return c.cmpPairs(pairs)
}

// addPlainMap recursively adds all elements of the map to the diff tree
// by comparing them to a nil value.
func (c *Comparator) addPlainMap(a ActionType, v reflect.Value) (*DiffNode, error) {
	if a != CREATE && a != DELETE {
		return nil, fmt.Errorf("addPlainMap: invalid action: %v", a)
	}

	x := reflect.New(v.Type()).Elem()
	pairs := NewPairs(v.Type())

	for _, k := range v.MapKeys() {
		vi := v.MapIndex(k)
		xi := x.MapIndex(k)

		switch a {
		case CREATE:
			pairs.addA(k, &xi)
			pairs.addB(k, &vi)
		case DELETE:
			pairs.addA(k, &vi)
			pairs.addB(k, &xi)
		}
	}

	node, err := c.cmpPairs(pairs)
	node.setActionToLeafs(a)

	return node, err
}
