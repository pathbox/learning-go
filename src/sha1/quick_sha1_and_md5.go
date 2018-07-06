package main

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
)

func main() {
	passwd := "123456789"
	sha1Str := sha1.Sum([]byte(passwd))
	md5Str := md5.Sum([]byte(passwd))
	fmt.Println("sha1 password: ", sha1Str)
	fmt.Println("md5 password: ", md5Str)
	fmt.Printf("sha1 password: %s\n", sha1Str)
	fmt.Printf("md5 password: %s\n", md5Str)
	fmt.Printf("sha1 password: %x\n", sha1Str)
	fmt.Printf("md5 password: %x\n", md5Str)
}

// quick sha1 or md5, you get the 0x 16进制的byte数据，不能直接格式化为字符串，而是要格式化为十六进制的字符串 使用%x, 能得到16进制的字符串
