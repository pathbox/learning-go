package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	input := []byte("hello golang base64 hello golang base64+-*/%")

	// 演示base64编码
	encodeString := base64.StdEncoding.EncodeToString(input)
	fmt.Println(encodeString)

	// 对上面的编码结果进行base64解码
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(decodeBytes))

	fmt.Println()

	// 如果要用在url中，需要使用URLEncoding
	uEnc := base64.URLEncoding.EncodeToString([]byte(input))
	fmt.Println(uEnc)

	uDec, err := base64.URLEncoding.DecodeString(uEnc)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(uDec))

	fmt.Println()
	dst := make([]byte, 1<<10)
	base64.StdEncoding.Encode(dst, input)
	fmt.Println(string(dst))

	fmt.Println()
	eDst := make([]byte, 1<<10)
	_, err = base64.StdEncoding.Decode(eDst, dst)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("===", string(eDst))
}
