package main

import (
	"encoding/hex"
	"fmt"
)

const MsgIDLength = 16

type MessageID [MsgIDLength]byte

func main() {
	i := int64(1000)
	r := Hex(i)
	fmt.Println(r)
	var bb []byte

	for _, re := range r {
		bb = append(bb, re)
	}
	fmt.Println(string(bb))
}

func Hex(id int64) MessageID {
	var h MessageID
	var b [8]byte // 8字节

	b[0] = byte(id >> 56) // 右移
	b[1] = byte(id >> 48)
	b[2] = byte(id >> 40)
	b[3] = byte(id >> 32)
	b[4] = byte(id >> 24)
	b[5] = byte(id >> 16)
	b[6] = byte(id >> 8)
	b[7] = byte(id)

	hex.Encode(h[:], b[:])
	return h
}

/*
[48 48 48 48 48 48 48 48 48 48 48 48 48 51 101 56]
00000000000003e8
*/
