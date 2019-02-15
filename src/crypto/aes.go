package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
)

//AESCBCEncrypt AesCBC加密PKCS5
//
//@param origData []byte 加密的字节数组
//
//@param key []byte 密钥字节数组
func AESCBCEncrypt(origData []byte, key []byte) ([]byte, error) {
	srcKey := MD5(key)
	block, err := aes.NewCipher(srcKey)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	cypted := make([]byte, len(origData))
	blockMode := cipher.NewCBCEncrypter(block, srcKey[:blockSize])
	blockMode.CryptBlocks(cypted, origData)
	return []byte(hex.EncodeToString(cypted)), nil
}

//AESECBEncrypt Aes ECB加密PKCS5
//
//@param origData []byte 加密的字节数组
//
//@param key []byte 密钥字节数组
func AESECBEncrypt(origData []byte, key []byte) ([]byte, error) {
	srckey := MD5(key)
	block, err := aes.NewCipher(srckey)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := newECBEncrypter(block)
	cypted := make([]byte, len(origData))
	blockMode.CryptBlocks(cypted, origData)
	return []byte(hex.EncodeToString(cypted)), nil
}

//AESCBCDecrypt Aes CBC解密PKCS5
//
//@param origData []byte 解密的字节数组
//
//@param key []byte 密钥字节数组
func AESCBCDecrypt(cypted []byte, key []byte) ([]byte, error) {
	cypted, err := hex.DecodeString(string(cypted))
	if err != nil {
		return nil, err
	}
	srckey := MD5(key)
	block, err := aes.NewCipher(srckey)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, srckey[:blockSize])
	origData := make([]byte, len(cypted))
	blockMode.CryptBlocks(origData, cypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

//AESECBDecrypt Aes ECB解密PKCS5
//
//@param origData []byte 解密的字节数组
//
//@param key []byte 密钥字节数组
func AESECBDecrypt(cypted []byte, key []byte) ([]byte, error) {
	cypted, err := hex.DecodeString(string(cypted))
	if err != nil {
		return nil, err
	}
	srckey := MD5(key)
	block, err := aes.NewCipher(srckey)
	if err != nil {
		return nil, err
	}
	blockMode := newECBDecrypter(block)
	origData := make([]byte, len(cypted))
	blockMode.CryptBlocks(origData, cypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize // 最后一个block中占多少size
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//PKCS5UnPadding PKCS5密钥填充方式返填充
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func MD5(key []byte) []byte {
	m := md5.New()
	m.Write(key)
	return m.Sum(nil)
}
