package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(3 * time.Second)

	for {
		select {
		case now := <-ticker.C:
			fmt.Printf("tick %s\n", now.UTC().Format("20060102-150405.000000000"))
		}
	}
}

// Post url
// https://chilts.org/2017/06/12/cancelling-multiple-goroutines?utm_source=golangweekly&utm_medium=email
