// channel  可以模拟 sync.Cond的Broadcast 的方法，但是 不能替代 mutex锁的功能

package main 

import "sync"

func main() {
	cv := sync.NewCond(new(sync.Mutex))
	done := false 

	go func() {
		// do something

		cv.L.Lock()
		done = true 
		cv.Signal()
		cv.L.Unlock()
	}()

	// wait something is done 
	cv.L.Lock()
	for !done {
		cv.Wait()
	}
	cv.L.Unlock()
}

package main

func main() {
	done := make(chan struct{})

	go func() {
		// do something

		close(done)
		// done <- struct{}
	}()

	// wait something is done 
	<-done 
}