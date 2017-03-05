package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	go func() {
		fmt.Println("before send")

	}()

	<-c

	time.Sleep(1 * time.Second)
	fmt.Println("done")
}

// before send
// fatal error: all goroutines are asleep - deadlock!

// goroutine 1 [chan receive]:
// main.main()
//   /Users/pathbox/code/learning-go/src/channel/example/example_lock4.go:15 +0x7f
// exit status 2

// receive <-c in main goroutinue
// no send
// deadlock!
