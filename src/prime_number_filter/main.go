package main

import (
	"fmt"
)

func generate(ch chan<- int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

func filter(src <-chan int, dst chan<- int, prime int) {
	for i := range src { // loop over values received from src
		if i%prime != 0 {
			dst <- i // send i to channel dst
		}
	}
}

func main() {
	ch := make(chan int)
	go generate(ch)
	for {
		prime := <-ch
		fmt.Println("prime: ", prime)
		ch1 := make(chan int)
		go filter(ch, ch1, prime)
		ch = ch1
	}
}
