package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Println("do something")
				time.Sleep(1 * time.Second)
				// other works ...
			}
		}
	}()

	time.Sleep(3 * time.Second)
	// signal all relevant goroutines
	// close(done)
	done <- true
}
