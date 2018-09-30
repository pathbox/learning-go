package main

import (
	"fmt"
	"math"
)

func main() {
	f := float32(1)
	fmt.Println(f)

	r := math.Float32bits(f) //符合IEEE 754浮点数标准

	fmt.Println(r)
}
