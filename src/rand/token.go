package main

import (
	"crypto/rand"
	// "encoding/base64"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	// return base64.URLEncoding.EncodeToString(b), err  // base64的编码
	return hex.EncodeToString(b), err // 16进制的编码
}

func GenerateRandomMd5StringTimeUnix(s int) (string, error) { // 为什么用md5，因为加密后能够防止别人轻易的解密
	t := time.Now().UnixNano()
	md5Newer := md5.New()
	timeStr := fmt.Sprintf("%d", t)
	salt := "718e7a3b4d2c0e32c3f78a02b62f56c5"
	str := timeStr + salt
	md5Newer.Write([]byte(str))
	hexResult := hex.EncodeToString(md5Newer.Sum([]byte("aaa"))) // “aaa”len 为3, 得到的最后结果长度： 32 + 2*3 = 38
	return hexResult, nil                                        // 得到32位长度的token字符串
}

func main() {
	token, err := GenerateRandomString(32)
	if err != nil {
		panic(err)
	}
	fmt.Println("Generate Token: ", token)
	md5Token, err := GenerateRandomMd5StringTimeUnix(32)
	fmt.Println(md5Token)
}
