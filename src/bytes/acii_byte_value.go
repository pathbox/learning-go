package main

import "fmt"

func main() {

	c := 'a'
	fmt.Println(byte(c))
	cr := byte(c) + '0' // 97 + 48 = 145
	fmt.Println(cr)

	crr := cr - '0'
	fmt.Println(crr)

	i := 0
	fmt.Println(byte(i))
	ir := byte(i) + '0'
	fmt.Println(ir)

	irr := ir - '0'
	fmt.Println(irr)
}
