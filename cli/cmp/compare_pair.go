package cmp

import (
	"reflect"
)

type Pair struct {
	key string
	A   *reflect.Value
	B   *reflect.Value
}

type Pairs struct {
	pairs   []*Pair
	typ     reflect.Type
	kind    reflect.Kind
	cmpById bool
}

func NewPairs(t reflect.Type) Pairs {
	return Pairs{
		pairs: make([]*Pair, 0),
		typ:   t,
		kind:  t.Kind(),
	}
}

func (ps *Pairs) addA(key interface{}, v *reflect.Value) {
	p := ps.get(key)
	p.A = v
}

func (ps *Pairs) addB(key interface{}, v *reflect.Value) {
	p := ps.get(key)
	p.B = v
}

func (ps *Pairs) get(key interface{}) *Pair {
	for _, p := range ps.pairs {
		if p.key == toString(key) {
			return p
		}
	}

	p := &Pair{
		key: toString(key),
	}

	ps.pairs = append(ps.pairs, p)

	return p
}

func (c *Comparator) cmpPairs(ps Pairs) (*DiffNode, error) {
	node := NewEmptyNode(ps.typ, ps.kind)

	for _, p := range ps.pairs {
		null := reflect.ValueOf(nil)

		if p.A == nil {
			p.A = &null
		}

		if p.B == nil {
			p.B = &null
		}

		child, err := c.compare(*p.A, *p.B)

		if err != nil {
			return nil, err
		}

		setSliceId(ps, child, c.TagName)
		node.addChild(child, p.key, p.key)
	}

	return node, nil
}

func setSliceId(ps Pairs, child *DiffNode, tagName string) {
	if !ps.cmpById || len(ps.pairs) == 0 || child == nil {
		return
	}

	p := ps.pairs[0]

	var rv reflect.Value

	if p.A != nil && p.A.Kind() != reflect.Invalid {
		rv = *p.A
	} else {
		rv = *p.B
	}

	fName := tagOptionIdFieldName(tagName, rv)

	for _, c := range child.children {
		if c.structKey == fName {
			c.isSliceId = true
			return
		}
	}
}
