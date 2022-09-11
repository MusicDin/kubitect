package cmp

import "reflect"

type Pair struct {
	A *reflect.Value
	B *reflect.Value
}

// PairMap is a map of comparative pairs.
type PairMap struct {
	m map[interface{}]*Pair
}

func NewPairMap() *PairMap {
	return &PairMap{
		m: make(map[interface{}]*Pair),
	}
}

func (pl *PairMap) addA(k interface{}, v *reflect.Value) {
	if (*pl).m[k] == nil {
		(*pl).m[k] = &Pair{}
	}
	(*pl).m[k].A = v
}

func (pl *PairMap) addB(k interface{}, v *reflect.Value) {
	if (*pl).m[k] == nil {
		(*pl).m[k] = &Pair{}
	}
	(*pl).m[k].B = v
}
