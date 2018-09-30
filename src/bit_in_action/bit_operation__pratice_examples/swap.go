package main

import "fmt"

func main() {
	a := 1
	b := 2
	fmt.Printf("a:%d, b:%d\n ", a, b)
	swap(&a, &b)
	fmt.Printf("a:%d, b:%d\n ", a, b)

}

func swap(a, b *int) {
	if *a != *b {
		*a ^= *b // *a = (*a^*b)
		*b ^= *a // *b = *b^(*a^*b) => *b = *b^*b^*a 由于一个数和自己异或的结果为0并且任何数与0异或都会不变的，所以此时b被赋上了a的值。
		*a ^= *b // *a = *a^*b => *a^*b^*a => *a^*a^*b => *b
	}

}
