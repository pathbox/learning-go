package main

import (
	"fmt"
	"io"
	"log"
)

// 多个 writer 将数据写入 pipe，一个 reader 读取数据
func main() {
	c := make(chan int)
	r, w := io.Pipe()
	go read(r, c)

	buf := []byte("abcdefg")
	for i := 0; i < len(buf); i++ {
		p := buf[i : i+1]
		n, err := w.Write(p)
		if n != len(p) {
			log.Fatalf("wrote %d, got %d", len(p), n)
		}
		if err != nil {
			log.Fatalf("write: %v", err)
		}
		nn := <-c // 读取了多少byte
		if nn != n {
			log.Fatalf("wrote %d, read got %d", n, nn)
		} else {
			log.Println(nn)
		}
	}

	w.Close()
	nn := <-c
	if nn != 0 {
		log.Fatalf("final read got %d", nn)
	}
}

func read(r io.Reader, c chan int) {
	for {
		var buf = make([]byte, 64)
		n, err := r.Read(buf)
		if err == io.EOF {
			c <- 0
			break
		}
		if err != nil {
			log.Fatalf("read fail: %v", err)
		}
		fmt.Printf("[read]: %s\n", buf[:n])
		c <- n
	}
}
