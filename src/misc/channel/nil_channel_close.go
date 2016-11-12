package main

import (
	"fmt"
	"time"
)

func WaitMany(a, b chan bool) {
	for a != nil || b != nil {
		select {
		case <-a:
			a = nil
		case <-b:
			b = nil
		}
	}
}

func main() {
	a, b := make(chan bool), make(chan bool)
	t := time.Now()
	go func() {
		close(a)
		close(b)
		// a <- true
		// b <- true
	}()
	WaitMany(a, b)
	fmt.Printf("waited %v for WaitMany\n", time.Since(t))
}
