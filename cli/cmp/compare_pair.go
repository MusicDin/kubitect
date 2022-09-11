package cmp

import (
	"reflect"
)

func (c *Comparator) diffPairs(parent *DiffNode, key interface{}, pl *PairMap) error {
	for k := range pl.m {
		null := reflect.ValueOf(nil)

		if pl.m[k].A == nil {
			pl.m[k].A = &null
		}

		if pl.m[k].B == nil {
			pl.m[k].B = &null
		}

		err := c.compare(parent, k, *pl.m[k].A, *pl.m[k].B)
		if err != nil {
			return err
		}
	}

	return nil
}
