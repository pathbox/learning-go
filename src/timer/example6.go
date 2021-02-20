package main

import (
	"fmt"
	"time"
)

func main() {
		t := time.NewTicker(time.Second*2)
		defer t.Stop()
		for {
			<- t.C
			fmt.Println("Ticker running...")
		}
}