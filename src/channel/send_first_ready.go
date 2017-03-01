package main

import "fmt"

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}

func main() {
	s := []int{7, 2, 8, -9, 4, 0}

	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c // receive from c
	// go sum(s[:len(s)/2], c)  deadlock
	// go sum(s[len(s)/2:], c)

	fmt.Println(x, y, x+y)
}

// send and recieve in diffrent goroutinue
// send(c<-) is must be ready before recieve(<-c)
// if recieve(<-c) is ready befor send(c<-), deadlock happend
// recieve(<-c) is waiting for send(c<-), it is blocking until send operator

// Sends to a buffered channel block only when the buffer is full. Receives block when the buffer is empty.
