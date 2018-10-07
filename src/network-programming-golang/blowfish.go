package main

import (
	"bytes"

	blowfish "golang.org/x/crypto/blowfish"

	"fmt"
)

func main() {

	key := []byte("my key")

	cipher, err := blowfish.NewCipher(key)

	if err != nil {

		fmt.Println(err.Error())

	}

	src := []byte("hello\n\n\n")

	var enc [512]byte

	cipher.Encrypt(enc[0:], src)

	var decrypt [8]byte

	cipher.Decrypt(decrypt[0:], enc[0:])

	result := bytes.NewBuffer(nil)

	result.Write(decrypt[0:8])

	fmt.Println(string(result.Bytes()))

}
