package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	const n = 100
	finish := make(chan struct{})
	var done sync.WaitGroup

	for i := 0; i < n; i++ {
		done.Add(1)
		go func() {
			select {
			case <-time.After(1 * time.Hour):
			case <-finish:
			}
			done.Done()
		}()
	}
	t0 := time.Now()
	close(finish)
	done.Wait()
	fmt.Printf("Waited %v for %d goroutines to stop\n", time.Since(t0), n)
}
