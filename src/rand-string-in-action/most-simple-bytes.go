package main

import (
	"fmt"
	"math/rand"
	"time"
)

var lbytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func init() {
	rand.Seed(time.Now().UnixNano())
}
func main() {
	n := 10 // the length of the rand string
	s := RandStringBytes(n)
	fmt.Println("rand string: ", s)
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = lbytes[rand.Intn(len(lbytes))]
	}

	return string(b)
}
