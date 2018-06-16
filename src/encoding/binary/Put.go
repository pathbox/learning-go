package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	b := make([]byte, 16)
	// binary.LittleEndian.PutUint16(b[0:], 0x03e8)
	binary.BigEndian.PutUint64(b[0:], 0x03e8)
	// binary.LittleEndian.PutUint16(b[2:], 0x07d0)
	binary.BigEndian.PutUint64(b[8:], 0x07d0)
	fmt.Printf("% x\n", b)
}
