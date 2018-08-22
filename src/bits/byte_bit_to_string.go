package main

import (
	"fmt"
	"strings"
)

func main() {
	bitmaps := make([]byte, 64)
	bitmaps[0] = byte(255)
	bitmaps[1] = byte(254)
	bitmaps[2] = byte('1') // => byte(49) 49是 1字符串的ACII码
	bitmaps[3] = byte('a')
	bitmaps[4] = byte('A')

	fmt.Println("bitmaps: ", bitmaps)
	s := GetBitString(bitmaps)
	fmt.Println("Result: ", s)
}

func GetBitString(bitmaps []byte) string {
	var offset, bitTotal uint64
	bitTotal = uint64(len(bitmaps) * 8)          //  一个字节占8位
	bitArr := make([]string, bitTotal, bitTotal) // 每个元素是0或1字符串

	for offset = 0; offset < bitTotal; offset++ {
		bitValue := GetBitValue(bitmaps, offset)
		bitArr[offset] = fmt.Sprintf("%d", bitValue)
	}

	return strings.Join(bitArr, "") // 将所有元素进行拼接
}

func GetBitValue(bitmaps []byte, offset uint64) uint8 {
	bitsize, index, pos := uint64(len(bitmaps)*8), offset/8, 7-offset%8

	if bitsize < offset {
		return 0
	}

	return (bitmaps[index] >> pos) & 0x01
}
