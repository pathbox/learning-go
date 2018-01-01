package main

import (
	"fmt"
	"time"
)

func main() {
	var Ball int
	table := make(chan int)
	// 2 个玩家
	go player(table) // 两个player 直接互相传输 table的值，直到 time.Sleep(2 * time.Second)
	go player(table)
	// go player(table)

	table <- Ball // first 从外部传入值
	time.Sleep(2 * time.Second)
	<-table // last
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
