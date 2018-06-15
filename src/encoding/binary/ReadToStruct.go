package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func main() {
	b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40, 0xff, 0x01, 0x02, 0x03, 0xbe, 0xef}

	l := len(b)
	fmt.Println("len b:", l)

	r := bytes.NewReader(b) // transform byte b to io.Reader r

	var data Data

	if err := binary.Read(r, binary.LittleEndian, &data); err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	fmt.Println(data)
	fmt.Println(data.PI)
	fmt.Println(data.Uate)
	fmt.Printf("% x\n", data.Mine)
	fmt.Println(data.Too)
}

type Data struct {
	PI   float64
	Uate uint8
	Mine [3]byte
	Too  uint16
}

// byte b 数据，反序列化为了 data struct数据. 14bytes 数据转为 data struct数据
// 就大小上来说，真是秒杀文本协议数据比如json
