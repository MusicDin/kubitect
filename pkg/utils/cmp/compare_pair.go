package cmp

import (
	"fmt"
	"reflect"
)

// Pair represents a pair of values from two Maps or Slices. The 'key' denotes
// the location in the parent map or slice. For maps, 'key' is the map key. For
// slices, it represents the slice ID or index.
type Pair struct {
	key reflect.Value
	A   reflect.Value
	B   reflect.Value
}

func (p *Pair) StringKey() string {
	return toString(p.key)
}

// Pairs encapsulates the differences between two corresponding data structures
// (either Maps or Slices). It provides metadata about the type and kind of
// pairs it contains, and offers control over how changes are represented
// within the collection.
type Pairs struct {
	pairs    []*Pair
	pairType reflect.Type
	pairKind reflect.Kind

	// compareById indicates whether the comparison should be based on ID
	// rather than index for slices.
	compareById bool

	// changeType represents the initial change type of Pairs. Default
	// value is Any, which means the node will inherit the change type
	// from children.
	changeType ChangeType
}

func NewPairs(t reflect.Type) Pairs {
	return Pairs{
		pairs:    []*Pair{},
		pairType: t,
		pairKind: t.Kind(),
	}
}

func (ps *Pairs) addA(key any, v reflect.Value) {
	p := ps.getOrCreate(toReflectValue(key))
	p.A = v
}

func (ps *Pairs) addB(key any, v reflect.Value) {
	p := ps.getOrCreate(toReflectValue(key))
	p.B = v
}

// getOrCreate retrieves a Pair using the provided key. If not found,
// it creates a new Pair, adds it to the list, and returns it.
func (ps *Pairs) getOrCreate(key reflect.Value) *Pair {
	newPair := &Pair{key: key}

	for _, p := range ps.pairs {
		if p.StringKey() == newPair.StringKey() {
			return p
		}
	}

	ps.pairs = append(ps.pairs, newPair)

	return newPair
}

// setNodeValues reconstructs the original before and after values of the node
// based on the pairs of reflected values.
func (ps *Pairs) setNodeValues(node *DiffNode) {
	var before reflect.Value
	var after reflect.Value

	switch ps.pairKind {
	case reflect.Map:
		before = reflect.MakeMap(ps.pairType)
		after = reflect.MakeMap(ps.pairType)

		for _, p := range ps.pairs {
			before.SetMapIndex(p.key, p.A)
			after.SetMapIndex(p.key, p.B)
		}

	case reflect.Slice | reflect.Array:
		before = reflect.MakeSlice(ps.pairType, 0, 0)
		after = reflect.MakeSlice(ps.pairType, 0, 0)

		for _, p := range ps.pairs {
			if p.A.IsValid() {
				reflect.Append(before, p.A)
			}

			if p.B.IsValid() {
				reflect.Append(after, p.B)
			}
		}
	}

	node.valueBefore = before.Interface()
	node.valueAfter = after.Interface()
}

func (c *Comparator) cmpPairs(ps Pairs) (*DiffNode, error) {
	if c.options.IgnoreEmptyChanges && len(ps.pairs) == 0 {
		return nil, nil
	}

	node := NewEmptyNode(ps.pairType, ps.pairKind)
	node.changeType = ps.changeType

	// If PopulateAllNodes is set to true then populate before and after
	// values of the node.
	if c.options.PopulateAllNodes {
		ps.setNodeValues(node)
	}

	for _, p := range ps.pairs {
		child, err := c.compare(p.A, p.B)
		if err != nil {
			return nil, err
		}

		if ps.compareById && child != nil {
			setSliceId(p, child, c.options.Tag)
		}

		node.addChild(child, p.StringKey(), p.StringKey())
	}

	return node, nil
}

// setSliceId iterates over node's children (slice indexes) and sets isSliceId
// property to true if children contains an id tag option. This is used to
// display slice ids even when no change has occurred.
func setSliceId(p *Pair, indexNode *DiffNode, tag string) {
	rv := p.A
	if !rv.IsValid() {
		rv = p.B
	}

	fName := findIdTaggedFieldName(rv, tag)
	for _, c := range indexNode.children {
		if c.pathStructKey == fName {
			c.isSliceId = true
			return
		}
	}
}

func toString(v any) string {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		return toString(reflect.Indirect(rv))
	}

	return fmt.Sprint(v)
}

func toReflectValue(v any) reflect.Value {
	if rv, ok := v.(reflect.Value); ok {
		return rv
	}

	return reflect.ValueOf(v)
}
