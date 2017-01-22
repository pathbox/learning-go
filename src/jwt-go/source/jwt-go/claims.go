package main

import (
	"fmt"
	"math/rand"
	// "strings"
	"time"
)

func main() {
	code_flg := "9"
	rand.Seed(time.Now().UnixNano())
	n := rand.Int63n(9999999999999)
	ten := fmt.Sprintf("%010d", n)
	result := code_flg + ten

	if len(result) > 1 && len(result) < 13 {
		fmt.Println(result)
	}

}
