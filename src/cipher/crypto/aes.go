package main

import (
	"crypto/aes"
	"crypto/cipher"
	// "crypto/rand"
	"errors"
	"fmt"
	// "io"
)

func main() {
	key := []byte("this is a rand key")
	plaintext := []byte("Hello World!你好世界")

	ciphertext, err := encrypt(plaintext, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("============", ciphertext)

	get_plaintext, _ := decrypt(ciphertext, key)
	fmt.Println("============", get_plaintext)
}

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	fmt.Println(plaintext)
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	// if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	// 	return nil, err
	// }
	fmt.Println(iv, ciphertext)
	stream := cipher.NewCTR(aesCipher, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	fmt.Println(ciphertext)
	return ciphertext, nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	if len(ciphertext) <= aes.BlockSize {
		return nil, errors.New("Invalid cipher text")
	}
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	fmt.Println(ciphertext)
	plaintext := make([]byte, len(ciphertext)-aes.BlockSize)
	fmt.Println(plaintext)
	stream := cipher.NewCTR(aesCipher, ciphertext[:aes.BlockSize])
	stream.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])
	fmt.Println(plaintext)
	return plaintext, nil
}
