package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	inputs := [][]byte{
		[]byte{0x81, 0x01},
		[]byte{0x7f},
		[]byte{0x03},
		[]byte{0x01},
		[]byte{0x00},
		[]byte{0x02},
		[]byte{0x04},
		[]byte{0x7e},
		[]byte{0x80, 0x01},
	}
	for _, b := range inputs {
		x, n := binary.Varint(b)
		if n != len(b) {
			fmt.Println("Varint did not consume all of in")
		}
		fmt.Println(x)
	}
}
