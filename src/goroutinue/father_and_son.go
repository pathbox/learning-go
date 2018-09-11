package main

import (
	"fmt"
	"time"
)

func main() {
	father()
	for {
	}
}

func father() {
	go son()
	fmt.Println("This is father")
}

func son() {
	time.Sleep(3 * time.Second)
	fmt.Println("This is son")
}

// This is father
// This is son
// son goroutinue is OK
