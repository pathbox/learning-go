package main

import (
	"fmt"

	"github.com/zheng-ji/goSnowFlake"
)

func main() {
	iw, err := goSnowFlake.NewIdWorker(1)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < 100; i++ {
		if id, err := iw.NextId(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(id)
		}
	}

}
