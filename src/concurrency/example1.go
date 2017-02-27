package main

import (
	"fmt"
	// "net/http"
	// "time"
)

func main() {
	chans := make(chan int64, 100)
	for i := 0; i < cap(chans); i++ { // 创建了100个goroutine 去处理逻辑
		go func() {

			// resp, _ := http.Get("http://www.baidu.com")
			// fmt.Println(resp)
			chans <- 100
			// defer resp.Body.Close()
		}()
	}

	// for c := range chans { // 这样会发生死锁
	j := 0
	for {
		// time.Sleep(1 * time.Second)
		fmt.Println("...........................")
		j++
		fmt.Println(j)
		select {
		case temp := <-chans:
			fmt.Println("============================")
			fmt.Println(temp)
		default:
			fmt.Println("hello")
			break
		}
	}
}
