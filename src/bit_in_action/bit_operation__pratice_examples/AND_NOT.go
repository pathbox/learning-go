// 与 和 非 ,先与再非，可以用于清除低位或高位

package main

import "fmt"

func main() {
	var a byte = 0xAB
	fmt.Printf("before %08b\n", a)
	a &^= 0x0F
	fmt.Printf("after %08b\n", a)
}
