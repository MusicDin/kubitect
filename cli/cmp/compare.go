package cmp

import (
	"reflect"
	"strings"
)

const (
	TAG_OPTION_ID = "id"
)

var fieldNameTags = []string{"json", "yaml"}

type CompareFunc func(reflect.Value, reflect.Value) (*DiffNode, error)

type Comparator struct {
	TagName           string
	RespectSliceOrder bool
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

// Compare compares the given values and returns a comparison tree.
func (c *Comparator) Compare(a, b interface{}) (*DiffNode, error) {
	return c.compare(reflect.ValueOf(a), reflect.ValueOf(b))
}

// compare recursively compares given values.
func (c *Comparator) compare(a, b reflect.Value) (*DiffNode, error) {
	if a.Kind() == reflect.Invalid && b.Kind() == reflect.Invalid {
		return NewLeaf(NONE, nil, nil), nil
	}

	cmpFunc := c.getCompareFunc(a, b)

	if cmpFunc == nil {
		return nil, NewTypeMismatchError(a.Kind(), b.Kind())
	}

	return cmpFunc(a, b)
}

// getCompareFunc returns a compare function based on the type of a
// comparative values.
func (c *Comparator) getCompareFunc(a, b reflect.Value) CompareFunc {
	switch {
	case areOfKind(a, b, reflect.Invalid, reflect.Bool):
		return c.cmpBasic
	case areOfKind(a, b, reflect.Invalid, reflect.String):
		return c.cmpBasic
	case areOfKind(a, b, reflect.Invalid, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64):
		return c.cmpBasic
	case areOfKind(a, b, reflect.Invalid, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64):
		return c.cmpBasic
	case areOfKind(a, b, reflect.Invalid, reflect.Float32, reflect.Float64):
		return c.cmpBasic
	case areOfKind(a, b, reflect.Invalid, reflect.Complex64, reflect.Complex128):
		return c.cmpBasic
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

func tagName(tagName string, field reflect.StructField) string {
	tags := append([]string{tagName}, fieldNameTags...)

	for _, tag := range tags {
		tName := strings.SplitN(field.Tag.Get(tag), ",", 2)[0]

		if len(tName) > 0 {
			return tName
		}
	}

	return ""
}

func tagOptionId(tagName string, v reflect.Value) interface{} {
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
	tag := field.Tag.Get(tagName)
	options := strings.Split(tag, ",")

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
