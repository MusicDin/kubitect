package cmp

import (
	"fmt"
	"reflect"
)

var emptyString string

func (c *Comparator) cmpString(parent *DiffNode, key interface{}, a, b reflect.Value) error {
	if a.Kind() == reflect.Invalid {
		if b.String() != emptyString {
			parent.addLeaf(CREATE, key, nil, exportInterface(b))
		}
		return nil
	}

	if b.Kind() == reflect.Invalid {
		if a.String() != emptyString {
			parent.addLeaf(DELETE, key, exportInterface(a), nil)
		}
		return nil
	}

	if a.Kind() != b.Kind() {
		return fmt.Errorf("Type mismatch: %v <> %v\n", a.Kind(), b.Kind())
	}

	if a.String() != b.String() {
		parent.addLeaf(MODIFY, key, a.String(), b.String())
		return nil
	}

	if a.String() != emptyString {
		parent.addLeaf(NONE, key, a.String(), b.String())
	}

	return nil
}
