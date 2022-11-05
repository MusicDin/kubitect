package cmp

import (
	"reflect"
)

func (c *Comparator) cmpSlice(a, b reflect.Value) (*DiffNode, error) {
	if a.Kind() == reflect.Invalid {
		a = reflect.New(b.Type()).Elem()
		return c.cmpSliceByIndex(a, b)
	}

	if b.Kind() == reflect.Invalid {
		b = reflect.New(a.Type()).Elem()
		return c.cmpSliceByIndex(a, b)
	}

	if c.areComparativeById(a, b) {
		return c.cmpSliceById(a, b)
	}

	return c.cmpSliceByIndex(a, b)
}

// cmpSliceByIndex compares slice elements located on the same index.
func (c *Comparator) cmpSliceByIndex(a, b reflect.Value) (*DiffNode, error) {
	var pairs Pairs

	matched := []bool{}
	for i := 0; i < a.Len(); i++ {
		ai := a.Index(i)

		if (c.RespectSliceOrder && !containsAtIndex(b, ai, i)) || (!c.RespectSliceOrder && !contains(b, ai, &matched)) {
			pairs.addA(toSliceKey(i), i, &ai)
		} else {
			pairs.addA(toSliceKey(i), i, &ai)
			pairs.addB(toSliceKey(i), i, &ai)
		}
	}

	matched = []bool{}
	missingCount := 0
	for i := 0; i < b.Len(); i++ {
		bi := b.Index(i)

		if c.RespectSliceOrder && !containsAtIndex(a, bi, i) {
			pairs.addB(toSliceKey(i), i, &bi)
		}

		if !c.RespectSliceOrder && !contains(a, bi, &matched) {
			j := a.Len() + missingCount
			pairs.addB(toSliceKey(j), i, &bi)
			missingCount++
		}
	}

	return c.diffPairs(pairs)
}

// cmpSliceByIndex compares slice elements based on the id element that is
// set with a tag.
func (c *Comparator) cmpSliceById(a, b reflect.Value) (*DiffNode, error) {
	var pairs Pairs

	for i := 0; i < a.Len(); i++ {
		ai := a.Index(i)
		av := getDeepValue(ai)

		id := tagOptionId(c.TagName, av)
		if id != nil {
			pairs.addA(toSliceKey(id), i, &ai)
		}
	}

	for i := 0; i < b.Len(); i++ {
		bi := b.Index(i)
		bv := getDeepValue(bi)

		id := tagOptionId(c.TagName, bv)
		if id != nil {
			pairs.addB(toSliceKey(id), i, &bi)
		}
	}

	return c.diffPairs(pairs)
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
