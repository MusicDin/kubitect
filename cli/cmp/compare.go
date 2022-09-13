package cmp

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	TAG_OPTION_ID = "id"
)

type CompareFunc func(*DiffNode, interface{}, reflect.Value, reflect.Value) error

type Comparator struct {
	TagName           string
	RespectSliceOrder bool
	SkipPrivateFields bool
}

func NewComparator() *Comparator {
	return &Comparator{
		TagName: "cmp",
	}
}

// Compare initializes default comparator, compares given values
// and returns a comparison tree.
func Compare(a, b interface{}) (*DiffNode, error) {
	c := NewComparator()
	return c.Compare(a, b)
}

// Compare compares the given elements of the same type and returns
// a comparison tree.
func (c *Comparator) Compare(a, b interface{}) (*DiffNode, error) {
	diff := NewNode()
	err := c.compare(diff, c.TagName, reflect.ValueOf(a), reflect.ValueOf(b))
	diff = diff.getChild(c.TagName)
	return diff, err
}

// compare recursively compares given values.
func (c *Comparator) compare(parent *DiffNode, key interface{}, a, b reflect.Value) error {
	cmpFunc := c.getCompareFunc(a, b)

	if cmpFunc == nil {
		return fmt.Errorf("Invalid compare type. Type mismatch: %s <> %s\n", a.Kind().String(), b.Kind().String())
	}

	return cmpFunc(parent, key, a, b)
}

// getCompareFunc returns a compare function based on the type of a
// comparative values.
func (c *Comparator) getCompareFunc(a, b reflect.Value) CompareFunc {
	switch {
	case areOfKind(a, b, reflect.Invalid, reflect.Int):
		return c.cmpInt
	case areOfKind(a, b, reflect.Invalid, reflect.String):
		return c.cmpString
	case areOfKind(a, b, reflect.Invalid, reflect.Struct):
		return c.cmpStruct
	case areOfKind(a, b, reflect.Invalid, reflect.Slice):
		return c.cmpSlice
	case areOfKind(a, b, reflect.Invalid, reflect.Map):
		return c.cmpMap
	case areOfKind(a, b, reflect.Invalid, reflect.Pointer):
		return c.cmpPointer
	case areOfKind(a, b, reflect.Invalid, reflect.Interface):
		return c.cmpInterface
	default:
		return nil
	}
}

// areComparativeById returns true if one of the values contains a
// tag option representing an ID element.
func (c *Comparator) areComparativeById(a, b reflect.Value) bool {
	if a.Len() > 0 {
		ai := a.Index(0)
		av := getDeepValue(ai)

		if av.Kind() == reflect.Struct {
			if hasTagOptionId(c.TagName, av) != nil {
				return true
			}
		}
	}

	if b.Len() > 0 {
		bi := b.Index(0)
		bv := getDeepValue(bi)

		if bv.Kind() == reflect.Struct {
			if hasTagOptionId(c.TagName, bv) != nil {
				return true
			}
		}
	}

	return false
}

func areOfKind(a, b reflect.Value, kinds ...reflect.Kind) bool {
	var isA, isB bool

	for _, k := range kinds {
		if a.Kind() == k {
			isA = true
		}

		if b.Kind() == k {
			isB = true
		}
	}

	return isA && isB
}

func areOfType(a, b reflect.Value, types ...reflect.Type) bool {
	var isA, isB bool

	for _, t := range types {
		if a.Kind() != reflect.Invalid {
			if a.Type() == t {
				isA = true
			}
		}
		if b.Kind() != reflect.Invalid {
			if b.Type() == t {
				isB = true
			}
		}
	}

	return isA && isB
}

func hasTagOptionId(tagName string, v reflect.Value) interface{} {
	if v.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < v.NumField(); i++ {
		if hasTagOption(tagName, v.Type().Field(i), TAG_OPTION_ID) {
			return exportInterface(v.Field(i))
		}
	}

	return nil
}

func hasTagOption(tagName string, field reflect.StructField, option string) bool {
	options := strings.Split(field.Tag.Get(tagName), ",")

	if len(options) < 2 {
		return false
	}

	for _, o := range options[1:] {
		o = strings.TrimSpace(o)
		o = strings.ToLower(o)

		if o == option {
			return true
		}
	}

	return false
}
