package main

import (
	"fmt"
	// "time"
)

func main() {

	c := make(chan bool, 1)

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
	r := <-c
	fmt.Println("========")
	fmt.Println(r)
}

/*
0
0
1
1
2
2
3
3
4
4
5
5
6
6
7
7
8
8
9
9

fatal error: all goroutines are asleep - deadlock!

c := make(chan bool, 1) then no deadlock
*/
