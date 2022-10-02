package main

import (
	"fmt"
	"log"

	"github.com/kevwan/mapreduce/v2"
)

func main() { // 传入了3个func
	val, err := mapreduce.MapReduce(func(source chan<- int) {
		for i := 0; i < 10; i++ {
			source <- i // 得到并发的次数
		}
	}, func(item int, writer mapreduce.Writer[int], cancel func(error)) {
		writer.Write(item * item) // 传入item 统计
	}, func(pipe <-chan int, writer mapreduce.Writer[int], cancel func(error)) {
		var sum int
		for i := range pipe { // 传入sum统计
			sum += i
		}
		writer.Write(sum)
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("result:", val)
}
