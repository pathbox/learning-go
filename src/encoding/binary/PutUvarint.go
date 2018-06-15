package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	buf := make([]byte, binary.MaxVarintLen64)

	for _, x := range []uint64{1, 2, 16, 127, 128, 255, 256} {
		n := binary.PutUvarint(buf, x)
		// fmt.Println(n)
		fmt.Printf("%x\n", buf[:n]) // 转为16进制打印数据
		fmt.Println(buf)
	}
	fmt.Println("======================")
	for _, x := range []int64{-65, -64, -2, -1, 0, 1, 2, 63, 64, 16, 127, 128, 255, 256} {
		n := binary.PutVarint(buf, x)
		fmt.Printf("%x\n", buf[:n])
		fmt.Println(buf)
	}
}

// PutVarint encodes an int64 into buf and returns the number of bytes written. If the buffer is too small, PutVarint will panic.每次buf都会被重置

/*
01
[1 0 0 0 0 0 0 0 0 0]
02
[2 0 0 0 0 0 0 0 0 0]
10
[16 0 0 0 0 0 0 0 0 0]
7f
[127 0 0 0 0 0 0 0 0 0]
8001
[128 1 0 0 0 0 0 0 0 0]
ff01
[255 1 0 0 0 0 0 0 0 0]
8002
[128 2 0 0 0 0 0 0 0 0]
======================
8101
[129 1 0 0 0 0 0 0 0 0]
7f
[127 1 0 0 0 0 0 0 0 0]
03
[3 1 0 0 0 0 0 0 0 0]
01
[1 1 0 0 0 0 0 0 0 0]
00
[0 1 0 0 0 0 0 0 0 0]
02
[2 1 0 0 0 0 0 0 0 0]
04
[4 1 0 0 0 0 0 0 0 0]
7e
[126 1 0 0 0 0 0 0 0 0]
8001
[128 1 0 0 0 0 0 0 0 0]
20
[32 1 0 0 0 0 0 0 0 0]
fe01
[254 1 0 0 0 0 0 0 0 0]
8002
[128 2 0 0 0 0 0 0 0 0]
fe03
[254 3 0 0 0 0 0 0 0 0]
8004
[128 4 0 0 0 0 0 0 0 0]
*/
