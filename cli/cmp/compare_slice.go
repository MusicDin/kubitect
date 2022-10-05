package cmp

import (
	"reflect"
)

func (c *Comparator) cmpSlice(parent *DiffNode, key interface{}, a, b reflect.Value) error {
	node := parent.addNode(key)

	if a.Kind() == reflect.Invalid && b.Kind() == reflect.Invalid {
		return nil
	}

	if a.Kind() == reflect.Invalid {
		a = reflect.New(b.Type()).Elem()
		return c.cmpSliceByIndex(node, key, a, b)
	}

	if b.Kind() == reflect.Invalid {
		b = reflect.New(a.Type()).Elem()
		return c.cmpSliceByIndex(node, key, a, b)
	}

	if a.Kind() != b.Kind() {
		return NewTypeMismatchError(a.Kind(), b.Kind())
	}

	if c.areComparativeById(a, b) {
		return c.cmpSliceById(node, key, a, b)
	}

	return c.cmpSliceByIndex(node, key, a, b)
}

// cmpSliceByIndex compares slice elements located on the same index.
func (c *Comparator) cmpSliceByIndex(n *DiffNode, key interface{}, a, b reflect.Value) error {
	pairs := NewPairMap()

	matched := []bool{}
	for i := 0; i < a.Len(); i++ {
		ai := a.Index(i)

		if (c.RespectSliceOrder && !containsAtIndex(b, ai, i)) || (!c.RespectSliceOrder && !contains(b, ai, &matched)) {
			pairs.addA(toSliceKey(i), &ai)
		} else {
			pairs.addA(toSliceKey(i), &ai)
			pairs.addB(toSliceKey(i), &ai)
		}
	}

	matched = []bool{}
	missingCount := 0
	for i := 0; i < b.Len(); i++ {
		bi := b.Index(i)

		if c.RespectSliceOrder && !containsAtIndex(a, bi, i) {
			pairs.addB(toSliceKey(i), &bi)
		}

		if !c.RespectSliceOrder && !contains(a, bi, &matched) {
			j := a.Len() + missingCount
			pairs.addB(toSliceKey(j), &bi)
			missingCount++
		}
	}

	if len(pairs.m) > 0 {
		return c.diffPairs(n, key, pairs)
	}

	return nil
}

// cmpSliceByIndex compares slice elements based on the id element that is
// set with a tag.
func (c *Comparator) cmpSliceById(n *DiffNode, key interface{}, a, b reflect.Value) error {
	pairs := NewPairMap()

	for i := 0; i < a.Len(); i++ {
		ai := a.Index(i)
		av := getDeepValue(ai)

		id := tagOptionId(c.TagName, av)
		if id != nil {
			pairs.addA(toSliceKey(id), &ai)
		}
	}

	for i := 0; i < b.Len(); i++ {
		bi := b.Index(i)
		bv := getDeepValue(bi)

		id := tagOptionId(c.TagName, bv)
		if id != nil {
			pairs.addB(toSliceKey(id), &bi)
		}
	}

	return c.diffPairs(n, key, pairs)
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
