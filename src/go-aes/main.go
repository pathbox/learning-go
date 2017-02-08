package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func main() {
	testAes()
}

func testAes() {
	// AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
	key := []byte("sfe023f_9fd&fwfl")                       // key 的长度只能为 16, 24, 或32
	result, err := AesEncrypt([]byte("Good morning!"), key) // 加密
	if err != nil {
		panic(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(result)) // 解密
	origData, err := AesDecrypt(result, key)
	if err != nil {
		panic(err)
	}
	//fmt.Println(origData) // 是 []byte
	fmt.Println(string(origData))
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key) // 根据key，得到block
	if err != nil {
		return nil, err
	}
	fmt.Println(block)
	blockSize := block.BlockSize()
	fmt.Println(blockSize)
	origData = PKCS5Padding(origData, blockSize) // 进行padding操作,给origData后面加一个字节的byte
	fmt.Println(origData)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	fmt.Println(blockMode)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData) // 将origData最后加密然后赋值给crypted
	fmt.Println(crypted)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)

	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
