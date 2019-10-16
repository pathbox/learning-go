package main

import (
	"fmt"
	"time"
)

func main() {
	go binlogListener()

	time.Sleep(2 * time.Minute)
	fmt.Print("Thx for watching, goodbuy")
}