package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	b := []byte{0xe8, 0x03, 0xd0, 0x07}
	x1 := binary.LittleEndian.Uint16(b[0:])
	x2 := binary.LittleEndian.Uint16(b[2:])
	fmt.Printf("%#06x %#04x %#x\n", x1, x2, x1)
}
