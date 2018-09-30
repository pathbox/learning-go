package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	lBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	lIdxBits = 6
	lIdxMask = (1 << lIdxBits) - 1
)

func main() {
	n := 10
	s := RandStringBytesMask(n)
	fmt.Println("rand string: ", s)
}

func RandStringBytesMask(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; {
		// 代替了取余的方式
		if idx := int(rand.Int63() & lIdxMask); idx < len(lBytes) {
			b[i] = lBytes[idx]
			i++
		}
	}
	return string(b)
}
