package main

import (
	"fmt"
	"time"
)

func main() {
	var Ball int
	table := make(chan int)

	// 36个玩家
	for i := 0; i < 36; i++ {
		go player(table)
	}
	table <- Ball
	time.Sleep(2 * time.Second)
	<-table
}

func player(table chan int) {
	for {
		ball := <-table
		fmt.Println(ball)
		ball++
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}
