package main

import "fmt"

func main() {
	ch := make(chan bool, 2)
	ch <- true
	ch <- true
	close(ch)

	for v := range ch {
		fmt.Println(v)
	}
}
