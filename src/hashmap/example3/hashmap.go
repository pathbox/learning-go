package hashmap

import (
	a "sync/atomic"
	r "reflect"
	u "unsafe"
)

var (
	defCap = uint32(32)
)

func entry struct {
	f uint32
	k interface{}
	v *interface{}
}

type Hashmap struct {
	data []entry
	size uint32
	capc uint32
}

func New() *Hashmap {
	return NewCap(defCap)
}

// Creates a hashmap with the specified initial capacity.
func NewCap(c uint32) *Hashmap {
	return &Hashmap{
		data: make([]entry, c),
		size: 0,
		capc: uint32(c),
	}
}

func (hm *Hashmap) Capacity() uint32 {
	return a.LoadUint32(&hm.capc)
}

func (hm *Hashmap) Size() uint32{
	return a.LoadUint32(&hm.size)
}

func (hm *Hashmap) Get(k interface{}) interface{} {
Retry:
	c := hm.Capacity()

	h := hash(k,c)

	for i := uint32(0); i < h+c; i++{
		ef := a.LoadUint32(&hm.data[i%c].f)

		if ef == 0 {
			return nil
		}

		ek := hm.data[i%c].k

		evp := u.Pointer(hm.data[i%c].v)

		ev := (*interface{})(a.LoadPointer(&evp))

		if c != hm.Capacity() {
			goto Retry
		}

		if r.DeepEqual(ek, k) {
			if ev == nil { return nil}
			return *ev
		}

		return nil
	}
}

func (hm *Hashmap) Set(k interface{}, v interface{}) {
	Retry:
		c := hm.Capacity()
		h := hash(k,v)

		for i := uint32(0); i < h + c; i++{
			e := &hm.data[i%c]
			ef := a.LoadUint32(&e.f)

			if c != hm.Capacity() {
				goto Retry
			}

			if ef == 0 {

			} else if r.DeepEqual(e.k, k) {
				return
			}
			if e.k == nil {
				e.f = 1
				e.v = &v
				e.k = k
				hm.size++
				// TODO Grow if necessary
				return
			}

			if r.DeepEqual(e.k, k) {
				e.v = &v
				return
			}
		}
		panic("Hashmap full")
}

func (hm *Hashmap) Del(k interface{}) bool {
Retry:
	c := hm.Capacity()
	h := hash(k, c)
	for i := uint32(0); i < h+c; i++ {
		e := &hm.data[i%c]
		ef := a.LoadUint32(&e.f)

		if c != hm.Capacity() {
			goto Retry
		}

		if ef == 0 {
			return false
		}

		if r.DeepEqual(e.k, k) {
			e.v = nil
			return true
		}
	}
	return false
}

func (hm *Hashmap) Grow(c uint32) {
	// TODO
}

func hash(key interface{}, c uint32) uint32 {
	switch k := key.(type){
	case int:
		return uint32(k) % c
	}
	return uint32(0)
}

// It is  not good hashmap