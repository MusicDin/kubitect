package cmp

import (
	"reflect"
)

func (c *Comparator) cmpStruct(a, b reflect.Value) (*DiffNode, error) {
	var changeType ChangeType
	var node *DiffNode

	if c.options.PopulateAllNodes {
		node = c.newNode(toInterface(a), toInterface(b))
	}

	// If either 'a' or 'b' is nil, replace its value with and empty struct
	// and set the change type accordingly.
	if !a.IsValid() {
		changeType = Create
		a = reflect.New(b.Type()).Elem()
	}

	if !b.IsValid() {
		changeType = Delete
		b = reflect.New(a.Type()).Elem()
	}

	if node == nil {
		node = NewEmptyNode(a.Type(), a.Type().Kind())
	}

	// Compare fields of struct 'a' with corresponding fields in struct 'b'.
	for i := 0; i < a.NumField(); i++ {
		field := a.Type().Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		fName := field.Name
		tName := getFieldNameFromTag(field, c.nameTags())

		// Skip fields marked with "-" in their tags.
		if tName == "-" {
			continue
		}

		// Use the struct name if field name is not set with the tag.
		if tName == "" {
			tName = fName
		}

		af := a.Field(i)
		bf := b.FieldByName(fName)

		child, err := c.compare(af, bf)
		if err != nil {
			return nil, err
		}

		node.addChild(child, tName, fName)
	}

	// Consider a struct empty if it's a leaf. Struct can be a leaf only if
	// it has no fields or if before and after values are nil for all its
	// fields.
	if c.options.IgnoreEmptyChanges && node.IsLeaf() {
		return nil, nil
	}

	// Apply determined change type to child nodes if necessary.
	if changeType != Any {
		node.setChangeTypeOfChildren(changeType)
	}

	return node, nil
}

// toInterface converts reflect value into an interface.
func toInterface(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}

	if !v.CanInterface() {
		panic("Cannot retrieve the value of an unexported field.")
	}

	return v.Interface()
}
