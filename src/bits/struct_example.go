package main

import (
	"encoding/binary"
	"fmt"
)

// 将struct 序列化和反序列化为 bit []byte 的例子

type MyStruct struct {
	Field1 int32
	Field2 string
	Field3 []int16
}

func main() {
	var s1 = MyStruct{255, "0aA", []int16{255, 256, 257}}
	var s2 MyStruct

	tmp := s1.Marshal()

	s2.Unmarsharl(tmp)

	fmt.Println("tmp:", tmp)

	fmt.Println(s2)
}

func (s *MyStruct) binarySize() int {
	return 4 + // Field1 int32 占4个字节
		2 + len(s.Field2) + // Len + Field2  开头的2个字节用于记录Field2的长度
		2 + 2*len(s.Field3) // Len + Field3 in16占2个字节 开头的2个字节用于记录Field3的长度
}

// 序列化
func (s *MyStruct) Marshal() []byte {
	b := make([]byte, s.binarySize())
	n := 0

	binary.BigEndian.PutUint32(b[n:], uint32(s.Field1)) // 处理Field1
	n += 4

	binary.BigEndian.PutUint16(b[n:], uint16(len(s.Field2))) // 处理Field2的长度
	n += 2

	copy(b[n:], s.Field2) // 处理Field2
	n += len(s.Field2)

	binary.BigEndian.PutUint16(b[n:], uint16(len(s.Field3)))
	n += 2

	for i := 0; i < len(s.Field3); i++ {
		binary.BigEndian.PutUint16(b[n:], uint16(s.Field3[i]))
		n += 2
	}

	return b

}

func (s *MyStruct) Unmarsharl(b []byte) {
	n := 0

	s.Field1 = int32(binary.BigEndian.Uint32(b[n:]))
	n += 4

	x := int(binary.BigEndian.Uint16(b[n:]))
	n += 2

	s.Field2 = string(b[n : n+x])
	n += x

	s.Field3 = make([]int16, binary.BigEndian.Uint16(b[n:]))
	n += 2

	for i := 0; i < len(s.Field3); i++ {
		s.Field3[i] = int16(binary.BigEndian.Uint16(b[n:]))
		n += 2
	}
}

// Field1 + len(Field2) + Field2 + len(Field3) + Field3
// tmp: [0 0 0 255 0 3 48 97 65 0 3 0 1 0 2 0 3]
// 想法记录： 将大整数用 bit的方式存储表示，可以大大减少存储空间。 如果将大整数用字符串表示是最浪费存储空间的
