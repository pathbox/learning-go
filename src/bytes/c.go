package main

import "fmt"

type Buffer struct {
	buf []byte // contents are the bytes buf[off : len(buf)]
	off int    // read at &buf[off], write at &buf[len(buf)]

}

func NewBufferString(s string) *Buffer {
	return &Buffer{buf: []byte(s)}
}

func main() {
	bb := NewBufferString("Hello World")
	bc := make([]byte, 6)

	n := copy(bc, bb.buf[0:])
	fmt.Println(n)
	fmt.Println(string(bc))
	fmt.Println(string(bb.buf))
}
