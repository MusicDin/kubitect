package cmp

import (
	"fmt"
	"reflect"
)

const skipPrivateFields = false

func (c *Comparator) cmpStruct(parent *DiffNode, key interface{}, a, b reflect.Value) error {
	node := parent.addNode(key)

	if a.Kind() == reflect.Invalid {
		return c.addPlainStruct(CREATE, node, key, b)
	}

	if b.Kind() == reflect.Invalid {
		return c.addPlainStruct(DELETE, node, key, a)
	}

	for i := 0; i < a.NumField(); i++ {
		if skipPrivateFields && !a.CanInterface() {
			continue
		}

		var af, bf reflect.Value
		field := a.Type().Field(i)

		if a.Kind() != reflect.Invalid {
			af = a.Field(i)
		}

		if b.Kind() != reflect.Invalid {
			bf = b.FieldByName(field.Name)
		}

		err := c.compare(node, field.Name, af, bf)
		if err != nil {
			return err
		}
	}

	return nil
}

// addPlainStruct recursively adds all elements of the struct to the diff tree
// by comparing them to a nil value.
func (c *Comparator) addPlainStruct(a ActionType, n *DiffNode, k interface{}, v reflect.Value) error {
	if a != CREATE && a != DELETE {
		return fmt.Errorf("Invalid action: %v", a)
	}

	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("Cannot add plain struct! Invalid kind: %s", v.Kind())
	}

	x := reflect.New(v.Type()).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)

		vf := v.Field(i)
		xf := x.FieldByName(field.Name)

		var err error
		switch a {
		case CREATE:
			err = c.compare(n, field.Name, xf, vf)
		case DELETE:
			err = c.compare(n, field.Name, vf, xf)
		}

		if err != nil {
			return err
		}
	}

	n.setActionToLeafs(a)

	return nil
}
