package cmp

import (
	"reflect"
)

func (c *Comparator) cmpSlice(a, b reflect.Value) (*DiffNode, error) {
	if a.Kind() == reflect.Invalid {
		a = reflect.New(b.Type()).Elem()
	}

	if b.Kind() == reflect.Invalid {
		b = reflect.New(a.Type()).Elem()
	}

	if c.areComparativeById(a, b) {
		return c.cmpSliceById(a, b)
	}

	return c.cmpSliceByIndex(a, b)
}

// cmpSliceByIndex compares slice elements located on the same index.
func (c *Comparator) cmpSliceByIndex(a, b reflect.Value) (*DiffNode, error) {
	pairs := NewPairs(a.Type())

	matched := []bool{}
	for i := 0; i < a.Len(); i++ {
		ai := a.Index(i)

		if (c.RespectSliceOrder && !containsAtIndex(b, ai, i)) || (!c.RespectSliceOrder && !contains(b, ai, &matched)) {
			pairs.addA(i, &ai)
		} else {
			pairs.addA(i, &ai)
			pairs.addB(i, &ai)
		}
	}

	matched = []bool{}
	missingCount := 0
	for i := 0; i < b.Len(); i++ {
		bi := b.Index(i)

		if c.RespectSliceOrder && !containsAtIndex(a, bi, i) {
			pairs.addB(i, &bi)
		}

		if !c.RespectSliceOrder && !contains(a, bi, &matched) {
			j := a.Len() + missingCount
			pairs.addB(j, &bi)
			missingCount++
		}
	}

	return c.cmpPairs(pairs)
}

// cmpSliceByIndex compares slice elements based on the id element that is
// set with a tag.
func (c *Comparator) cmpSliceById(a, b reflect.Value) (*DiffNode, error) {
	pairs := NewPairs(a.Type())
	pairs.cmpById = true

	for i := 0; i < a.Len(); i++ {
		ai := a.Index(i)

		if id := tagOptionId(c.Tag, ai); id != nil {
			pairs.addA(id, &ai)
		}
	}

	for i := 0; i < b.Len(); i++ {
		bi := b.Index(i)

		if id := tagOptionId(c.Tag, bi); id != nil {
			pairs.addB(id, &bi)
		}
	}

	return c.cmpPairs(pairs)
}

// areComparativeById returns true if one of the values contains a
// tag option representing an ID element.
func (c *Comparator) areComparativeById(a, b reflect.Value) bool {
	if a.Len() > 0 {
		ai := a.Index(0)
		av := getDeepValue(ai)

		if av.Kind() == reflect.Struct {
			if tagOptionId(c.Tag, av) != nil {
				return true
			}
		}
	}

	if b.Len() > 0 {
		bi := b.Index(0)
		bv := getDeepValue(bi)

		if bv.Kind() == reflect.Struct {
			if tagOptionId(c.Tag, bv) != nil {
				return true
			}
		}
	}

	return false
}

// contains checks whether a slice s contains an element x
func contains(s, x reflect.Value, matched *[]bool) bool {
	if len(*matched) != s.Len() {
		*matched = make([]bool, s.Len())
	}

	for i := 0; i < s.Len(); i++ {
		if (*matched)[i] {
			continue
		}

		if containsAtIndex(s, x, i) {
			(*matched)[i] = true
			return true
		}
	}

	return false
}

// containsAtIndex checks whether a slice s contains an element x at index i.
func containsAtIndex(s, x reflect.Value, i int) bool {
	if i < s.Len() {
		si := s.Index(i)
		return reflect.DeepEqual(exportInterface(si), exportInterface(x))
	}

	return false
}
