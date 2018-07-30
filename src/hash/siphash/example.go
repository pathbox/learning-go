package main

import "fmt"
import "github.com/dchest/siphash"

const (
	sipHashKey1 = 0xdda7806a4847ec61
	sipHashKey2 = 0xb5940c2623a5aabd
)

func main() {

	x := []byte("100")
	r := siphash.Hash(sipHashKey1, sipHashKey2, x)
	fmt.Println(r)
}
