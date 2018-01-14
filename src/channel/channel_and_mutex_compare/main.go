package main

import "sync"

var mutex = sync.Mutex{}
var ch = make(chan bool, 1)

func UseMutex() {
	mutex.Lock()

	mutex.Unlock()
}
func UseChan() {
	ch <- true
	<-ch
}
