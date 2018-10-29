package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	addr := "127.0.0.1:5055"
	count := 500
	client, _ := net.Dial("tcp", addr)
	for i := 0; i < count; i++ {
		go func() {
			_, err := client.Write([]byte("This is a tcp message\n"))
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	defer client.Close()
	time.Sleep(10 * time.Second)
}
