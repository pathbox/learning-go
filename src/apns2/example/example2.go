package main

import (
	"log"
	"runtime"
	"sync"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

func init() {
	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)
}

func main() {
	clients := make(chan *apns2.Client, 100)
	cert, err := certificate.FromP12File("/home/user/ios_dev.p12", "123456")

	if err != nil {
		log.Fatal("cert error: ", err)
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		// create an apns2 Client per CPU core
		clients <- apns2.NewClient(cert).Development()
	}

	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		client := <-clients // grab a client from pool
		clients <- client   // add the client back to the pool
		wg.Add(1)
		go func() {
			// grab a notification from your channel filled with notifications
			notification := <-notifications
			res, err := client.Push(notification)
			if err != nil {
				log.Panic(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
