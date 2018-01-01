package main

import (
	"fmt"
)

func main() {
	ch := make(chan int)
	go generate(ch)
	for i := 0; i < 100; i++ {
		prime := <-ch
		fmt.Println(prime)
		out := make(chan int)

		go filter(ch, out, prime)
		ch = out
	}
}

func generate(ch chan<- int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

func filter(src <-chan int, dst chan<- int, prime int) {
	for i := range src {
		if i%prime != 0 {
			dst <- i
		}
	}
}
