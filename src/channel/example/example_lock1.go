package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("before received")

	}()
	fmt.Println("before send")
	c <- 1
	fmt.Println("after received")
}

// before send
// before received
// fatal error: all goroutines are asleep - deadlock!

// goroutine 1 [chan send]:
// main.main()
//   /Users/pathbox/code/learning-go/src/channel/example/example_lock1.go:16 +0x10c
// exit status 2

// sned in main goroutinue
// no receive, send the data to channel, deadlock!
