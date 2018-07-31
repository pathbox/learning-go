package hashmap

import (
	"sync/atomic"
	"unsafe"
)

type ListElement struct {
	keyHash         uintptr
	previousElement unsafe.Pointer
	nextElement     unsafe.Pointer
	key             interface{}
	value           unsafe.Pointer
	deleted         uintptr
}

func (e *ListElement) Value() (value interface{}) {
	return *(*interface{})(atomic.LoadPointer(&e.value))
}

func (e *ListElement) Next() *ListElement {
	return (*ListElement)(atomic.LoadPointer(&e.nextElement))
}

func (e *ListElement) Previous() *ListElement {
	return (*ListElement)(atomic.LoadPointer(&e.previousElement))
}

func (e *ListElement) setValue(value unsafe.Pointer) {
	atomic.StorePointer(&e.value, value)
}

func (e *ListElement) casValue(from interface{}, to unsafe.Pointer) bool {
	old := atomic.LoadPointer(&e.value)
	if *(*interface{})(old) != from {
		return false
	}
	return atomic.CompareAndSwapPointer(&e.value, old, to)
}
