package cmp

import (
	"reflect"
	"strings"
)

const tagOptionId = "id"

type CompareFunc func(reflect.Value, reflect.Value) (*DiffNode, error)

// Options define configurations that influence how the comparison is executed.
type Options struct {
	// Tag specifies the primary tag from which the field name and options
	// are derived.
	Tag string

	// ExtraNameTags lists additional tags checked for field names. The
	// primary Tag always takes precedence.
	ExtraNameTags []string

	// RespectSliceOrder, when set, enforces slice comparison by index.
	// Otherwise, slice's fields are first checked if they contain id tag
	// option.
	RespectSliceOrder bool

	// IgnoreEmptyChanges, when set, omits tracking fields that remain nil
	// both before and after a change.
	IgnoreEmptyChanges bool

	// PopulateAllNodes ensures that 'before' and 'after' values are
	// retained even for non-leaf changes.
	PopulateAllNodes bool
}

type Comparator struct {
	options Options
}

// Compare initializes default comparator, compares given values
// and returns a comparison tree.
func Compare(a any, b any, options ...Options) (*Result, error) {
	opts := Options{}
	if len(options) > 0 {
		opts = options[0]
	}

	if opts.Tag == "" {
		opts.Tag = "cmp"
	}

	c := &Comparator{
		options: opts,
	}

	d, err := c.compare(reflect.ValueOf(a), reflect.ValueOf(b))
	if err != nil {
		return nil, err
	}

	return &Result{d}, nil
}

// compare recursively compares given values.
func (c *Comparator) compare(a, b reflect.Value) (*DiffNode, error) {
	if !a.IsValid() && !b.IsValid() {
		return nil, nil
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
		fallthrough
	case areOfKind(a, b, reflect.Invalid, reflect.String):
		fallthrough
	case areOfKind(a, b, reflect.Invalid, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64):
		fallthrough
	case areOfKind(a, b, reflect.Invalid, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64):
		fallthrough
	case areOfKind(a, b, reflect.Invalid, reflect.Float32, reflect.Float64):
		fallthrough
	case areOfKind(a, b, reflect.Invalid, reflect.Complex64, reflect.Complex128):
		return c.cmpBasic
	case areOfKind(a, b, reflect.Invalid, reflect.Array, reflect.Slice):
		return c.cmpSlice
	case areOfKind(a, b, reflect.Invalid, reflect.Map):
		return c.cmpMap
	case areOfKind(a, b, reflect.Invalid, reflect.Struct):
		return c.cmpStruct
	case areOfKind(a, b, reflect.Invalid, reflect.Pointer):
		return c.cmpPointer
	case areOfKind(a, b, reflect.Invalid, reflect.Interface):
		return c.cmpInterface
	default:
		return nil
	}
}

// nameTags returns the list of tags that should be used for searching field's name.
// Note that the primary tag always takes precedence.
func (c *Comparator) nameTags() []string {
	return append([]string{c.options.Tag}, c.options.ExtraNameTags...)
}

// areOfKind checks whether both provided values equal one of the provided types.
// This is used to determine an appropriate compare function.
func areOfKind(a, b reflect.Value, kinds ...reflect.Kind) bool {
	isA := false
	isB := false

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

// hasTagOption returns true if field's tag contains the given option.
func hasTagOption(field reflect.StructField, tagName string, option string) bool {
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

// getFieldNameFromTag returns the field's name extracted from any given tag (tag's
// content before first comma). Tags are checked in the provided order.
func getFieldNameFromTag(field reflect.StructField, tags []string) string {
	for _, tag := range tags {
		tagName := strings.SplitN(field.Tag.Get(tag), ",", 2)[0]

		if len(tagName) > 0 {
			return tagName
		}
	}

	return ""
}

// findIdTaggedField returns the struct field that contains id as an option in
// its tag. If none is found, nil is returned.
func findIdTaggedField(v reflect.Value, tag string) any {
	v = getDeepValue(v)

	if v.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < v.NumField(); i++ {
		if hasTagOption(v.Type().Field(i), tag, tagOptionId) {
			return v.Field(i)
		}
	}

	return nil
}

// findIdTaggedFieldName returns the name of a struct field that contains id
// as an option in its tag. If none is found, nil is returned.
func findIdTaggedFieldName(v reflect.Value, tag string) any {
	v = getDeepValue(v)

	if v.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < v.NumField(); i++ {
		if hasTagOption(v.Type().Field(i), tag, tagOptionId) {
			return v.Type().Field(i).Name
		}
	}

	return nil
}

// getDeepValue recursively dereferences a reflect.Value until it reaches
// a non-pointer, non-interface value.
func getDeepValue(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Pointer:
		return getDeepValue(reflect.Indirect(v))
	case reflect.Interface:
		return v.Elem()
	default:
		return v
	}
}
