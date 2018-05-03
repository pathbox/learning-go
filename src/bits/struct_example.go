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

/*
255用二进制表达就是1111 1111，再加1就是1 0000 0000，多了一个1出来，显然我们需要再用额外的一个字节来存放这个1，但是这个1要存放在第一个字节还是第二个字节呢？这时候因为人们选择的不同，就出现了大端序和小端序的差异。

当我们把这个1放在第一个字节的时候，就称之为大端序格式。当我们把1放在第二个字节的时候，就称之为小端序格式。

这两种格式显然没办法说谁更好，所以两个格式一直都各自的支持者，如果是按标准实现一个通讯协议，那就得严格按照标准上说的字节序来实现。如果是自定义的二进制协议，选择哪个格式按自己喜好就可以了。

encoding/binary包中的全局变量BigEndian用于操作大端序数据，LittleEndian用于操作小端序数据

从上面结构体序列化和反序列化的代码中，大家不难看出，实现一个二进制协议是挺繁琐和容易出BUG的，只要稍微有一个数值计算错就解析出错了。

所以在工程实践中，不推荐大家手写二进制协议的解析代码，项目中通常会用自动化的工具来辅助生成代码
*/
