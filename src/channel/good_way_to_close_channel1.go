package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	// ...
	const MaxRandomNumber = 100000
	const NumReceivers = 100

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	// ...
	dataCh := make(chan int, 100)

	// the sender
	go func() {
		for {
			if value := rand.Intn(MaxRandomNumber); value == 0 {
				// the only sender can close the channel safely.
				close(dataCh)
				return
			} else {
				dataCh <- value
			}
		}
	}()

	// receivers
	for i := 0; i < NumReceivers; i++ {
		go func() {
			defer wgReceivers.Done()

			// receive values until dataCh is closed and
			// the value buffer queue of dataCh is empty.
			for value := range dataCh {
				log.Println(value)
			}
		}()
	}

	wgReceivers.Wait()
}

//  M个receivers，一个sender，sender通过关闭data channel说“不再发送”
// 这是最简单的场景了，就只是当sender不想再发送的时候让sender关闭data 来关闭channel
