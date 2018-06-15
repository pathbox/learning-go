package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func main() {
	var pi float64
	b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &pi)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	fmt.Println(pi)

	buf = bytes.NewReader(b)
	err = binary.Read(buf, binary.BigEndian, &pi)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	fmt.Println(pi)

}
