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

	// if m, err = zcrypt.DecryptWithPrivateKey("$PRIVATE_KEY", "+3mY3lZfdowmVsUHF7cPN9vC+8HI7pK5mMnT53cZYi6H6YqDxA3ZEfSNU66DukZ9ppE08RCzYl6sHT394I8XIN+8wBvMQk990pGj/hIoelr/JBaPxS9vfKkEABiGx+0x"); err == nil {
	// 	return errors.Wrap(err, "couldn't decrypt the secret message")
	// }
	// fmt.Print(m)
}
