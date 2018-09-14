package main

import (
	"fmt"

	"github.com/xxtea/xxtea-go/xxtea"
)

func main() {
	str := "Hello World! 你好，中国！"
	key := "123456"
	encrypt_data := xxtea.Encrypt([]byte(str), []byte(key))
	fmt.Println("encrypt_data: ", string(encrypt_data))
	decrypt_data := xxtea.Decrypt(encrypt_data, []byte(key))
	fmt.Println("decrypt_data: ", string(decrypt_data))
	if str == string(decrypt_data) {
		fmt.Println("success!")
	} else {
		fmt.Println("fail!")
	}
}
