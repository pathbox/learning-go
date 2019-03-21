package main

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

const INT_SIZE int = int(unsafe.Sizeof(0))

// 判断系统中的字符序类型
func systemEdian() {
	var i int = 0x1 // 用16进制的表示方法表示整数
	fmt.Println("i value:", i)
	bs := (*[INT_SIZE]byte)(unsafe.Pointer(&i))
	fmt.Println(bs)
	if bs[0] == 0 {
		fmt.Println("system edian is little endian")
	} else {
		fmt.Println("system edian is big endian")
	}
}

func testBigEndian() {
	var testInt int32 = 256 // int32  uint32 表示4个字节一组，一个数值
	fmt.Printf("%d use big endian: \n", testInt)
	var testBytes []byte = make([]byte, 4)
	binary.BigEndian.PutUint32(testBytes, uint32(testInt))
	fmt.Println("int32 to bytes:", testBytes)

	convInt := binary.BigEndian.Uint32(testBytes)
	fmt.Printf("bytes to int32: %d\n\n", convInt)
}

func testLittleEndian() {

	var testInt int32 = 256
	fmt.Printf("%d use little endian: \n", testInt)
	var testBytes []byte = make([]byte, 4)
	binary.LittleEndian.PutUint32(testBytes, uint32(testInt))
	fmt.Println("int32 to bytes:", testBytes)

	convInt := binary.LittleEndian.Uint32(testBytes)
	fmt.Printf("bytes to int32: %d\n\n", convInt)
}

func main() {
	systemEdian()
	fmt.Println("")
	testBigEndian()
	testLittleEndian()
}
