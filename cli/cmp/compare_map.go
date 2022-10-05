package cmp

import (
	"fmt"
	"reflect"
)

func (c *Comparator) cmpMap(parent *DiffNode, key interface{}, a, b reflect.Value) error {
	node := parent.addNode(key)

	if a.Kind() == reflect.Invalid {
		return c.addPlainMap(CREATE, node, key, b)
	}

	if b.Kind() == reflect.Invalid {
		return c.addPlainMap(DELETE, node, key, a)
	}

	pairs := NewPairMap()

	for _, k := range a.MapKeys() {
		ai := a.MapIndex(k)
		if ai.Kind() != reflect.Invalid {
			pairs.addA(exportInterface(k), &ai)
		}
	}

	for _, k := range b.MapKeys() {
		bi := b.MapIndex(k)
		if bi.Kind() != reflect.Invalid {
			pairs.addB(exportInterface(k), &bi)
		}
	}

	return c.diffPairs(node, key, pairs)
}

// addPlainMap recursively adds all elements of the map to the diff tree
// by comparing them to a nil value.
func (c *Comparator) addPlainMap(a ActionType, p *DiffNode, key interface{}, v reflect.Value) error {
	if a != CREATE && a != DELETE {
		return fmt.Errorf("Invalid action: %v", a)
	}

	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}

	if v.Kind() != reflect.Map {
		return fmt.Errorf("Cannot add plain map! Invalid kind: %s", v.Kind())
	}

	x := reflect.New(v.Type()).Elem()

	for _, k := range v.MapKeys() {
		vi := v.MapIndex(k)
		xi := x.MapIndex(k)

		var err error
		switch a {
		case CREATE:
			err = c.compare(p, k, xi, vi)
		case DELETE:
			err = c.compare(p, k, vi, xi)
		}

		if err != nil {
			return err
		}
	}

	p.setActionToLeafs(a)

	return nil
}
