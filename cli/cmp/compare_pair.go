package cmp

import (
	"reflect"
)

type Pair struct {
	key       interface{}
	structKey interface{}
	A         *reflect.Value
	B         *reflect.Value
}

type Pairs []*Pair

func (p *Pairs) get(key, structKey interface{}) *Pair {
	for _, pair := range *p {
		if pair.key == toString(key) {
			return pair
		}
	}

	pair := &Pair{
		key:       toString(key),
		structKey: toString(structKey),
	}

	*p = append(*p, pair)

	return pair
}

func (p *Pairs) addA(key, structKey interface{}, v *reflect.Value) {
	pair := p.get(key, structKey)
	pair.A = v
}

func (p *Pairs) addB(key, structKey interface{}, v *reflect.Value) {
	pair := p.get(key, structKey)
	pair.B = v
}

func (c *Comparator) diffPairs(pl Pairs) (*DiffNode, error) {
	node := NewEmptyNode()

	for _, p := range pl {
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

		node.addChild(child, p.key, p.structKey)
	}

	return node, nil
}
