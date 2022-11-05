package cmp

import (
	"fmt"
	"reflect"
)

func (c *Comparator) cmpStruct(a, b reflect.Value) (*DiffNode, error) {
	if a.Kind() == reflect.Invalid {
		return c.addPlainStruct(CREATE, b)
	}

	if b.Kind() == reflect.Invalid {
		return c.addPlainStruct(DELETE, a)
	}

	node := NewEmptyNode()

	for i := 0; i < a.NumField(); i++ {
		var af, bf reflect.Value

		field := a.Type().Field(i)

		if !field.IsExported() {
			continue
		}

		fName := field.Name
		tName := tagName(c.TagName, field)

		if tName == "-" {
			continue
		}

		if tName == "" {
			tName = fName
		}

		if a.Kind() != reflect.Invalid {
			af = a.Field(i)
		}

		if b.Kind() != reflect.Invalid {
			bf = b.FieldByName(field.Name)
		}

		child, err := c.compare(af, bf)

		if err != nil {
			return nil, err
		}

		node.addChild(child, tName, fName)
	}

	return node, nil
}

// addPlainStruct recursively adds all elements of the struct to the diff tree
// by comparing them to a nil value.
func (c *Comparator) addPlainStruct(a ActionType, v reflect.Value) (*DiffNode, error) {
	if a != CREATE && a != DELETE {
		return nil, fmt.Errorf("addPlainStruct: invalid action: %v", a)
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("addPlainStruct: invalid kind: %s", v.Kind())
	}

	n := NewEmptyNode()
	x := reflect.New(v.Type()).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)

		if !field.IsExported() {
			continue
		}

		fName := field.Name
		tName := tagName(c.TagName, field)

		if tName == "-" {
			continue
		}

		if tName == "" {
			tName = fName
		}

		vf := v.Field(i)
		xf := x.FieldByName(field.Name)

		var err error
		var child *DiffNode

		switch a {
		case CREATE:
			child, err = c.compare(xf, vf)
		case DELETE:
			child, err = c.compare(vf, xf)
		}

		if err != nil {
			return nil, err
		}

		n.addChild(child, tName, fName)
	}

	n.setActionToLeafs(a)

	return n, nil
}
