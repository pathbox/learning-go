package main

import (
	"reflect"
	"unsafe"
)

func incX(d interface{}) int64 {
	v := reflect.ValueOf(d).Elem()
	f := v.FieldByName("X")

	x := f.Int()
	x++
	f.SetInt(x)

	return x
}

var offset uintptr = 0xFFFF // 避开 offset = 0 的字段

func unsafeIncX(d interface{}) int64 {
	if offset == 0xFFFF {
		t := reflect.TypeOf(d).Elem()
		x, _ := t.FieldByName("X")
		offset = x.Offset
	}

	p := (*[2]uintptr)(unsafe.Pointer(&d))
	px := (*int64)(unsafe.Pointer(p[1] + offset))
	*px++
	return *px
}

var cache = map[*uintptr]map[string]uintptr{}

func unsafeCacheIncX(d interface{}) int64 {
	itab := *(**uintptr)(unsafe.Pointer(&d))

	m, ok := cache[itab]
	if !ok {
		m = make(map[string]uintptr)
		cache[itab] = m
	}

	offset, ok := m["X"]
	if !ok {
		t := reflect.TypeOf(d).Elem()
		x, _ := t.FieldByName("X")
		offset = x.Offset

		m["X"] = offset
	}

	p := (*[2]uintptr)(unsafe.Pointer(&d))
	px := (*int64)(unsafe.Pointer(p[1] + offset))
	*px++
	return *px
}

func main() {
	d := struct{ X int }{100}
	println(incX(&d))
}
