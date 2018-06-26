package main

import (
	"fmt"
	"strconv"
	"time"
)

func makeCakeAndSend(cs chan string, count int) {
	for i := 1; i <= count; i++ {
		cakeName := "Strawberry Cake " + strconv.Itoa(i)
		cs <- cakeName //send a strawberry cake
	}
	close(cs) // 如果不对close cs， range cs会一直阻塞从cs读取数据，receiveCakeAndPack goroutine会一直在运行，除非 main 都结束了
}

func receiveCakeAndPack(cs chan string) {
	for s := range cs {
		fmt.Println("Packing received cake: ", s)
	}
}

func main() {
	cs := make(chan string)
	go makeCakeAndSend(cs, 50)
	go receiveCakeAndPack(cs)

	//sleep for a while so that the program doesn’t exit immediately
	time.Sleep(3 * 1e9)
}

// Go提供了range关键字，将其使用在channel上时，会自动等待channel的动作一直到channel被关闭,
// 通过对channel使用range关键字，我们避免了给接收者写明要接收的数据个数这种不合理的需求——当channel被关闭时，接收者的for循环也被自动停止了
