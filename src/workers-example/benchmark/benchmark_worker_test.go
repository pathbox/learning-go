package main

import (
	"sync"
	"testing"
)

var eventChan = make(chan int, 10000)
var eventPoolChan = make(chan int, 10000)
var wg = &sync.WaitGroup{}

func BenchmarkMultiWorker(b *testing.B) {

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go printAction()
	}

	wg.Wait()
}

func BenchmarkEventWorker(b *testing.B) {

	go func() {
		for {
			select {
			case <-eventChan:
				printAction()
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		eventChan <- i
	}
	wg.Wait()
}

func BenchmarkEventPoolWorker(b *testing.B) {
	for n := 0; n < 10; n++ {
		go func() {
			for {
				select {
				case <-eventPoolChan:
					printAction()
				}
			}
		}()
	}

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		eventPoolChan <- i
	}
	wg.Wait()
}

func printAction() {
	wg.Done()
	for i := 0; i < 1000; i++ {
		a := 0
		a++
	}
}

// go test -v -bench . -benchmem
