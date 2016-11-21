package main

import (
	"log"
	"time"

	"github.com/mediocregopher/okq-go/okq"
)

func stopTimeout(t <-chan time.Time) chan bool {
	b := make(chan bool, 1)
	select {
	case <-t:
		b <- true
		return b
	}
}

func main() {
	cl := okq.New("127.0.0.1:9999")
	defer cl.Close()

	stopCh := make(chan bool, 1)
	go func() {
		time.Sleep(2 * time.Second)
		log.Println("before stop channel")
		close(stopCh)
		log.Println("after stop channel")
	}()

	fn := func(e *okq.Event) bool {
		log.Printf("event received on %s: %s\n", e.Queue, e.Contents)
		time.Sleep(5 * time.Second)
		log.Println("event stop itself")
		return true
	}

	for {
		err := cl.Consumer(fn, stopCh, "super-queue")
		if err != nil {
			log.Printf("Error received from consumer: %s\n", err)
		}
	}
}
