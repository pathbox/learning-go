package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	// a channel to tell `tick()` and `tock()` to stop
	stopChan := make(chan struct{})

	// a channel for `tick()` to tell us they've stopped
	tickStoppedChan := make(chan struct{})
	go tick(stopChan, tickStoppedChan)

	// a channel for `tock()` to tell us they've stopped
	tockStoppedChan := make(chan struct{})
	go tock(stopChan, tockStoppedChan)

	// listen for C-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("main: received C-c - shutting down")

	// tell the goroutine to stop
	fmt.Println("main: telling goroutines to stop")
	close(stopChan)
	// and wait for them to reply back
	<-tickStoppedChan
	<-tockStoppedChan
	fmt.Println("main: all goroutines have told us they've finished")
}

func tick(stop, stopped chan struct{}) {
	// tell the caller we've stopped
	defer close(stopped)

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			fmt.Printf("tick: tick %s\n", now.UTC().Format("20060102-150405.000000000"))
		case <-stop:
			fmt.Println("tick: caller has told us to stop")
			return
		}
	}
}

func tock(stop, stopped chan struct{}) {
	// tell the caller we've stopped
	defer close(stopped)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			fmt.Printf("tock: tock %s\n", now.UTC().Format("20060102-150405.000000000"))
		case <-stop:
			fmt.Println("tock: caller has told us to stop")
			return
		}
	}
}

// 明白了 原理其实就是 两个chan， 一个 用来由外传到goroutine中， 告诉goroutine可以关闭了，一个由内传到外，告诉外层main goroutine
// 我这个goroutine已经关闭了，你可以继续执行下去，感觉这个不一定需要让main等子goroutine
