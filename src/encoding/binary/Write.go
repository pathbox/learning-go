package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

func main() {
	var buf bytes.Buffer
	var pi float64 = math.Pi
	err := binary.Write(&buf, binary.LittleEndian, pi)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Printf("% x\n", buf.Bytes())
	fmt.Println(buf.Bytes())
}

/*
18 2d 44 54 fb 21 09 40
[24 45 68 84 251 33 9 64]
用8byte 表示了 math.Pi

Write: 将data 转为bytes
*/
