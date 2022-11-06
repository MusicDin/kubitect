package cmp

import "reflect"

func (c *Comparator) cmpStruct(a, b reflect.Value) (*DiffNode, error) {
	var action ActionType

	if a.Kind() == reflect.Invalid {
		a = reflect.New(b.Type()).Elem()
		action = CREATE
	}

	if b.Kind() == reflect.Invalid {
		b = reflect.New(a.Type()).Elem()
		action = DELETE
	}

	t := a.Type()
	node := NewEmptyNode(t, t.Kind())

	for i := 0; i < a.NumField(); i++ {
		field := a.Type().Field(i)

		if !field.IsExported() {
			continue
		}

		fName := field.Name
		tName := tagName(c.TagName, field)

		if tName == "" {
			tName = fName
		}

		if tName == "-" {
			continue
		}

		af := a.Field(i)
		bf := b.FieldByName(fName)

		child, err := c.compare(af, bf)

		if err != nil {
			return nil, err
		}

		node.addChild(child, tName, fName)
	}

	if action != UNKNOWN {
		node.setActionToLeafs(action)
	}

	return node, nil
}
