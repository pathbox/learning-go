package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 4; i++ {
		x := rand.Intn(3)
		fmt.Println(x)
	}
}

func createPasswd() string {
	t := time.Now().UnixNano()
	h := md5.New()
	str := fmt.Sprintf("%d", t)
	str = str + "718e7a3b4d2c0e32c3f78a02b62f56c5"
	h.Write([]byte(str))
	hex := hex.EncodeToString(h.Sum(nil))
	passwd := fmt.Sprintf("%s", hex[0:16])
	return passwd
}

func createNumber() string { // 创建随机数字符串，保证为14位
	rand.Seed(time.Now().UnixNano())
	n := rand.Int63n(99999999999999)
	result := fmt.Sprintf("%014d", n)
	return result
}
