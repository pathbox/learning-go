// package main

// import (
// 	"fmt"
// 	"math/rand"
// )

// const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// const (
// 	letterIdxBits = 6                    // 6 bits to represent a letter index
// 	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
// 	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
// )

// func RandStringBytes(n int) string {
// 	b := make([]byte, n)
// 	for i := range b {
// 		b[i] = letterBytes[rand.Intn(len(letterBytes))]
// 	}

// 	return string(b)
// }

// func RandStringBytesMaskImpr(n int) string {
// 	b := make([]byte, n)

// 	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
// 		if remain == 0 {
// 			cache, remain = rand.Int63(), letterIdxMax
// 		}
// 		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
// 			b[i] = letterBytes[idx]
// 			i--
// 		}
// 		cache >>= letterIdxBits
// 		remain--
// 	}
// 	return string(b)
// }

// func main() {
// 	fmt.Println(RandStringBytes(10))
// 	fmt.Println(RandStringBytesMaskImpr(10))
// }

package main

import (
	"fmt"
)

func fb(n int) int {
	if n == 1 {
		return 1
	}

	if n == 2 {
		return 1
	}

	if n < 0 {
		return 0
	}

	s := fb(n-1) + fb(n-2)
	return s
}

func main() {
	n := 7
	r := fb(n)
	fmt.Println(r)

}
