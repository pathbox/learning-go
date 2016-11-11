package main

import (
	"time"
)

func main() {
	global := 0
	var lock = make(chan bool)

	lock <- true
	go func() {
		<-lock
		global = 1
		println(global)
		lock <- true
	}()

	go func() {
		<-lock
		global = 2
		println(global)
		lock <- true
	}()

	for {
		time.Sleep(1 * time.Second)
	}

}
