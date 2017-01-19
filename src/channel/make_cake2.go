package main

import (
	"fmt"
	"strconv"
	"time"
)

func makeCakeAndSend(cs chan string) {
	for i := 1; i <= 3; i++ {
		cakeName := "Strawberry Cake " + strconv.Itoa(i)
		fmt.Println("Making a cake and sending ...", cakeName)
		cs <- cakeName //send a strawberry cake
	}
}

func receiveCakeAndPack(cs chan string) {
	for i := 1; i <= 3; i++ {
		s := <-cs //get whatever cake is on the channel
		fmt.Println("Packing received cake: ", s)
	}
}

func main() {
	cs := make(chan string)
	go makeCakeAndSend(cs)
	go receiveCakeAndPack(cs)

	//sleep for a while so that the program doesn’t exit immediately
	time.Sleep(4 * 1e9)
}
