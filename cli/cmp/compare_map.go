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

	var pairs Pairs

	for _, k := range a.MapKeys() {
		ai := a.MapIndex(k)
		if ai.Kind() != reflect.Invalid {
			pairs.addA(k, k, &ai)
		}
	}

	for _, k := range b.MapKeys() {
		bi := b.MapIndex(k)
		if bi.Kind() != reflect.Invalid {
			pairs.addB(k, k, &bi)
		}
	}

	return c.diffPairs(pairs)
}

// addPlainMap recursively adds all elements of the map to the diff tree
// by comparing them to a nil value.
func (c *Comparator) addPlainMap(a ActionType, v reflect.Value) (*DiffNode, error) {
	if a != CREATE && a != DELETE {
		return nil, fmt.Errorf("addPlainMap: invalid action: %v", a)
	}

	x := reflect.New(v.Type()).Elem()
	node := NewEmptyNode()

	for _, k := range v.MapKeys() {
		vi := v.MapIndex(k)
		xi := x.MapIndex(k)

		var err error
		var child *DiffNode

		switch a {
		case CREATE:
			child, err = c.compare(xi, vi)
		case DELETE:
			child, err = c.compare(vi, xi)
		}

		if err != nil {
			return nil, err
		}

		node.addChild(child, k, k)
	}

	node.setActionToLeafs(a)

	return node, nil
}
