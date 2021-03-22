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
	go grandson()
}

func grandson() {
	time.Sleep(5 * time.Second)
	fmt.Println("This is grandson")
}

// This is father
// This is son
// son goroutinue is OK
