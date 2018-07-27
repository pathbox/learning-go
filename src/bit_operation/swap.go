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
		*a ^= *b
		*b ^= *a
		*a ^= *b
	}

}
