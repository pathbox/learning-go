package main

import (
	"fmt"
	"io"
	"os"
	"unsafe"
)

type iface struct {
	tab  *itab
	data unsafe.Pointer
}

type itab struct {
	inter uintptr
	_type uintptr
	link  uintptr
	hash  uint32
	_     [4]byte
	fun   [1]uintptr
}

type eface struct {
	_type uintptr
	data  unsafe.Pointer
}

func main() {
	var r io.Reader
	fmt.Printf("initial r: %T, %v\n", r, r)

	tty, _ := os.OpenFile("/Users/qcrao/Desktop/test", os.O_RDWR, 0)
	fmt.Printf("tty: %T, %v\n", tty, tty)

	// 给 r 赋值
	r = tty
	fmt.Printf("r: %T, %v\n", r, r)

	rIface := (*iface)(unsafe.Pointer(&r))
	fmt.Printf("r: iface.tab._type = %#x, iface.data = %#x\n", rIface.tab._type, rIface.data)

	// 给 w 赋值
	var w io.Writer
	w = r.(io.Writer)
	fmt.Printf("w: %T, %v\n", w, w)

	wIface := (*iface)(unsafe.Pointer(&w))
	fmt.Printf("w: iface.tab._type = %#x, iface.data = %#x\n", wIface.tab._type, wIface.data)

	// 给 empty 赋值
	var empty interface{}
	empty = w
	fmt.Printf("empty: %T, %v\n", empty, empty)

	emptyEface := (*eface)(unsafe.Pointer(&empty))
	fmt.Printf("empty: eface._type = %#x, eface.data = %#x\n", emptyEface._type, emptyEface.data)
}
