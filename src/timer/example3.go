package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second * 7)
			c <- false
		}

		time.Sleep(time.Second * 7)
		c <- true
	}()

	// consumer routine
	go func() {
		timer := time.NewTimer(time.Second * 5)
		for {
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(time.Second * 5)
			select {
			case b := <-c:
				if b == false {
					fmt.Println(time.Now(), ": recv false. continue")
					continue
				}

				fmt.Println(time.Now(), ":recv true. return")
				return
			case <-timer.C:
				fmt.Println(time.Now(), ":timer expired")
				continue
			}
		}
	}()

	var s string
	fmt.Scanln(&s)
}
