package main

import (
	"fmt"
	"time"

	"./syncreader"
)

func main() {
	path := "/Users/pathbox/tickets.json"

	r, _ := syncreader.ReadFile(path, 100)

	for {
		select {
		case s := <-r:
			fmt.Println("line:", s)
		case <-time.After(10 * time.Second):
			return
		default:
			fmt.Println("Block")
		}
	}

}
