package main

import (
	"fmt"

	"./zcrypt"
)

func main() {
	m, err := zcrypt.EncryptToPubKey("MaRI8ibKsgg+QqvRPDPRrh8NbOR2nsB2Mk81ctU4KEE=", "this is a secret message")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(m) // 得到一个加密并且base64之后的字符串
}
