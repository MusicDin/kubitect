package cmp

import (
	"reflect"
)

func (c *Comparator) cmpSlice(a, b reflect.Value) (*DiffNode, error) {
	if !a.IsValid() {
		return c.cmpSliceToNil(Create, b)
	}

	if !b.IsValid() {
		return c.cmpSliceToNil(Delete, a)
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

		// If slice order is respected, check if slice 'b' contains the
		// element from slice 'a' on the given index. If slice order is
		// not respected, check if any element in slice 'b' equals
		// element from slice 'a' on the given index.
		if c.options.RespectSliceOrder && !containsAtIndex(b, ai, i) || !c.options.RespectSliceOrder && !contains(b, ai, &matched) {
			pairs.addA(i, ai)
		} else {
			// We know that elements are equal, so we just create
			// a pair with two identical values.
			pairs.addA(i, ai)
			pairs.addB(i, ai)
		}
	}

	matched = []bool{}
	lastIndex := a.Len()
	for i := 0; i < b.Len(); i++ {
		bi := b.Index(i)

		// If element from slice 'b' is matched with an element from
		// slice 'a', just add the element to the existing pair.
		if c.options.RespectSliceOrder && !containsAtIndex(a, bi, i) {
			pairs.addB(i, bi)
		}

		// Since each element from slice 'a' is already inserted in
		// pairs, elements from slice 'b' that are not matched with
		// any other element from slice 'a' are simply appended to
		// pairs. Therefore the last index for a new element is calculated
		// as length of a slice 'a' plus the number of already added
		// elements from slice 'b'.
		if !c.options.RespectSliceOrder && !contains(a, bi, &matched) {
			pairs.addB(lastIndex, bi)
			lastIndex++
		}
	}

	return c.cmpPairs(pairs)
}

// cmpSliceByIndex compares slice elements based on the id element that is
// set with a tag.
func (c *Comparator) cmpSliceById(a, b reflect.Value) (*DiffNode, error) {
	pairs := NewPairs(a.Type())
	pairs.compareById = true

	for i := 0; i < a.Len(); i++ {
		ai := a.Index(i)
		if id := findIdTaggedField(ai, c.options.Tag); id != nil {
			pairs.addA(id, ai)
		}
	}

	for i := 0; i < b.Len(); i++ {
		bi := b.Index(i)
		if id := findIdTaggedField(bi, c.options.Tag); id != nil {
			pairs.addB(id, bi)
		}
	}

	return c.cmpPairs(pairs)
}

// cmpSliceToNil recursively adds all elements of the slice to the diff tree
// by comparing them to a nil value.
func (c *Comparator) cmpSliceToNil(t ChangeType, v reflect.Value) (*DiffNode, error) {
	pairs := NewPairs(v.Type())

	for i := 0; i < v.Len(); i++ {
		vi := v.Index(i)

		id := findIdTaggedField(vi, c.options.Tag)
		if id == nil {
			id = i
		} else {
			pairs.compareById = true
		}

		switch t {
		case Create:
			pairs.addB(id, vi)
		case Delete:
			pairs.addA(id, vi)
		}
	}

	node, err := c.cmpPairs(pairs)
	if err != nil {
		return nil, err
	}

	if node != nil {
		node.setChangeTypeOfChildren(t)
	}

	return node, nil
}

// areComparativeById returns true if any slice element contains an id tag
// option.
func (c *Comparator) areComparativeById(a, b reflect.Value) bool {
	if a.Len() > 0 {
		ai := a.Index(0)
		av := getDeepValue(ai)
		if av.Kind() == reflect.Struct && findIdTaggedField(av, c.options.Tag) != nil {
			return true
		}
	}

	if b.Len() > 0 {
		bi := b.Index(0)
		bv := getDeepValue(bi)
		if bv.Kind() == reflect.Struct && findIdTaggedField(bv, c.options.Tag) != nil {
			return true
		}
	}

	return false
}

// contains checks whether slice s contains an element x. An element from
// slice s is skipped if it was already matched before.
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
		return reflect.DeepEqual(si.Interface(), x.Interface())
	}

	return false
}
