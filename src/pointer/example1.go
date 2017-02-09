package main

import (
	"fmt"
)

func zeroval(ival int) {
	ival = 0
}

func zeroptr(iptr *int) {
	*iptr = 0
}

func main() {
	i := 1
	fmt.Println("initial:", i) // 1
	zeroval(i)
	fmt.Println("zeroval:", i) // 1

	zeroptr(&i)                // 参数为int指针，&i 就是传了这个变量参数的地址(指针指向的地址,也就是这个指针的值)
	fmt.Println("zeroptr:", i) // 0
	// Pointers can be printed too.
	fmt.Println("pointer:", &i) // pointer

}
