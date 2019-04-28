package jingo

import (
	"reflect"
	"unsafe"
)

type SliceEncoder struct {
	instruction func(t unsafe.Pointer, w *Buffer)
	tt          reflect.Type
	offset      uintptr
}

// Marshal executes the instruction set built up by NewSliceEncoder
func (e *SliceEncoder) Marshal(s interface{}, w *Buffer) {
	p := unsafe.Pointer(reflect.ValueOf(s).Pointer()) // 使用反射得到 s的pointer
	e.instruction(p, w)
}

// NewSliceEncoder builds a new SliceEncoder
func NewSliceEncoder(t interface{}) *SliceEncoder {
	e := &SliceEncoder{}

	e.tt = reflect.TypeOf(t)
	e.offset = e.tt.Elem().Size()

	switch e.tt.Elem().Kind() {
	case reflect.Slice:
		e.sliceInstr()
	case reflect.Struct:
		e.structInstr()
	case reflect.String:
		e.stringInstr()
	case reflect.Ptr:
		switch e.tt.Elem().Elem().Kind() {
		case reflect.Slice:
			e.ptrSliceInstr()
		case reflect.Struct:
			e.ptrStructInstr()
		case reflect.String:
			e.ptrStringInstr()
		default:
			e.ptrOtherInstr()
		}
	default:
		e.otherInstr()
	}

	return e
}

//  avoid allocs in the instruction
var (
	null = []byte("null")
	zero = uintptr(0)
)

func (e *SliceEncoder) sliceInstr() {
	enc := NewSliceEncoder(reflect.New(e.tt.Elem()).Elem().Interface())
	e.instruction = func(v unsafe.Pointer, w *Buffer) {
		w.WriteByte('[')

		sl := *(*reflect.SliceHeader)(v)
		for i := uintptr(0); i < uintptr(sl.Len); i++ {
			if i > zero {
				w.WriteByte(',')
			}
			s := unsafe.Pointer(sl.Data + (i * e.offset))
			enc.Marshal(s, w)
		}

		w.WriteByte(']')
	}
}

func (e *SliceEncoder) structInstr() {
	enc := NewStructEncoder(reflect.New(e.tt.Elem()).Elem().Interface())
	e.instruction = func(v unsafe.Pointer, w *Buffer) {
		w.WriteByte('[')

		sl := *(*reflect.SliceHeader)(v)
		for i := uintptr(0); i < uintptr(sl.Len); i++ {
			if i > zero {
				w.WriteByte(',')
			}
			s := unsafe.Pointer(sl.Data + (i * e.offset))
			enc.Marshal(s, w)
		}
		w.WriteByte(']')
	}
}

func (e *SliceEncoder) stringInstr() {
	e.instruction = func(v unsafe.Pointer, w *Buffer) {
		w.WriteByte('[')

		sl := *(*reflect.SliceHeader)(v)
		for i := uintptr(0); i < uintptr(sl.Len); i++ {
			if i > zero {
				w.WriteByte(',')
			}
			w.WriteByte('"')
			ptrStringToBuf(unsafe.Pointer(sl.Data+(i*e.offset)), w)
			w.WriteByte('"')
		}

		w.WriteByte(']')
	}
}

func (e *SliceEncoder) otherInstr() {

	conv, ok := typeconv[e.tt.Elem().Kind()]
	if !ok {
		return
	}

	e.instruction = func(v unsafe.Pointer, w *Buffer) {
		w.WriteByte('[')

		sl := *(*reflect.SliceHeader)(v)
		for i := uintptr(0); i < uintptr(sl.Len); i++ {
			if i > zero {
				w.WriteByte(',')
			}
			conv(unsafe.Pointer(sl.Data+(i*e.offset)), w)
		}

		w.WriteByte(']')
	}
}

func (e *SliceEncoder) ptrSliceInstr() {
	enc := NewSliceEncoder(reflect.New(e.tt.Elem()).Elem().Elem().Interface())
	e.instruction = func(v unsafe.Pointer, w *Buffer) {
		w.WriteByte('[')

		sl := *(*reflect.SliceHeader)(v)
		for i := uintptr(0); i < uintptr(sl.Len); i++ {
			if i > zero {
				w.WriteByte(',')
			}

			s := unsafe.Pointer(*(*uintptr)(unsafe.Pointer(sl.Data + (i * e.offset))))
			if s == unsafe.Pointer(nil) {
				w.Write(null)
				continue
			}
			enc.Marshal(s, w)
		}

		w.WriteByte(']')
	}
}

func (e *SliceEncoder) ptrStructInstr() {
	enc := NewStructEncoder(reflect.New(e.tt.Elem().Elem()).Elem().Interface())
	e.instruction = func(v unsafe.Pointer, w *Buffer) {
		w.WriteByte('[')

		sl := *(*reflect.SliceHeader)(v)
		for i := uintptr(0); i < uintptr(sl.Len); i++ {
			if i > zero {
				w.WriteByte(',')
			}

			s := unsafe.Pointer(*(*uintptr)(unsafe.Pointer(sl.Data + (i * e.offset))))
			if s == unsafe.Pointer(nil) {
				w.Write(null)
				continue
			}
			enc.Marshal(s, w)
		}

		w.WriteByte(']')
	}
}

func (e *SliceEncoder) ptrStringInstr() {
	e.instruction = func(v unsafe.Pointer, w *Buffer) {
		w.WriteByte('[')

		sl := *(*reflect.SliceHeader)(v)
		for i := uintptr(0); i < uintptr(sl.Len); i++ {
			if i > zero {
				w.WriteByte(',')
			}

			s := unsafe.Pointer(*(*uintptr)(unsafe.Pointer(sl.Data + (i * e.offset))))
			if s == unsafe.Pointer(nil) {
				w.Write(null)
				continue
			}
			w.WriteByte('"')
			ptrStringToBuf(s, w)
			w.WriteByte('"')
		}

		w.WriteByte(']')
	}
}

func (e *SliceEncoder) ptrOtherInstr() {

	conv, ok := typeconv[e.tt.Elem().Elem().Kind()]
	if !ok {
		return
	}

	e.instruction = func(v unsafe.Pointer, w *Buffer) {
		w.WriteByte('[')

		sl := *(*reflect.SliceHeader)(v)
		for i := uintptr(0); i < uintptr(sl.Len); i++ {
			if i > zero {
				w.WriteByte(',')
			}

			s := unsafe.Pointer(*(*uintptr)(unsafe.Pointer(sl.Data + (i * e.offset))))
			if s == unsafe.Pointer(nil) {
				w.Write(null)
				continue
			}
			conv(s, w)
		}

		w.WriteByte(']')
	}
}
