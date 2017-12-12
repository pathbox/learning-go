package pool

import (
	"reflect"
	"testing"
)

func TestGenericPoolGet(t *testing.T) {
	for _, test := range []struct {
		name     string
		min, max int
		init     func(int) interface{}
		get      int

		expSize int
		expObj  interface{}
	}{
		{
			min: 0,
			max: 1,
			get: 10,
			init: func(n int) interface{} {
				return uint(n)
			},

			expSize: 10,
			expObj:  uint(10),
		},
		{
			min: 0,
			max: 16,
			get: 10,
			init: func(n int) interface{} {
				return uint(n)
			},

			expSize: 16,
			expObj:  uint(16),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			p := New(test.min, test.max, test.init)
			x, n := p.Get(test.get)
			if n != test.expSize {
				t.Errorf("Get(%d) = _, %d; want %d", test.get, n, test.expSize)
			}
			if !reflect.DeepEqual(x, test.expObj) {
				t.Errorf("Get(%d) = %v, _; want %v", test.get, x, test.expObj)
			}
		})
	}
}
