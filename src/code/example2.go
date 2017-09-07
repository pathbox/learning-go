package main

import (
	"fmt"
	// "time"
)

func main() {

	c := make(chan bool)

	b := make(chan bool)

	done := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			<-c
			fmt.Println(i)
			b <- true
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			<-b
			fmt.Println(i)
			c <- true
		}
		done <- true
	}()

	c <- true
	<-done
}

/*
oo
xx
oo
xx
oo
xx
oo
xx
oo
xx
oo
xx
oo
xx
oo
xx
oo
xx
oo
xx
fatal error: all goroutines are asleep - deadlock!

c := make(chan bool, 1) then no deadlock
*/
