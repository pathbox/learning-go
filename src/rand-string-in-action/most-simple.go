package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var lRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func main() {
	n := 10 // the length of the rand string
	s := RandStringRunes(n)
	fmt.Println("rand string: ", s)
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = lRunes[rand.Intn(len(lRunes))]
	}

	return string(b)
}
