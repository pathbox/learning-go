package main

import (
	"fmt"
	"net"
	"time"

	"github.com/silenceper/pool"
)

func main() {
	addr := "127.0.0.1:9000"
	factory := func() (interface{}, error) {
		return net.Dial("tcp", addr)
	}

	closeF := func(v interface{}) error {
		return v.(net.Conn).Close()
	}

	poolConfig := &pool.PoolConfig{
		InitialCap:  5,
		MaxCap:      30,
		Factory:     factory,
		Close:       closeF,
		IdleTimeout: 150 * time.Second,
	}

	p, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		fmt.Println("err=", err)
	}

	// 从连接池中获取一个连接
	for i := 0; i < 100000; i++ {
		// i := 0
		// for {
		v, err := p.Get()
		if err != nil {
			panic(err)
		}

		conn := v.(net.Conn)
		// time.Sleep(1 * time.Second)
		conn.Write([]byte("Hello World!\n"))

		// 将连接放回连接池中
		p.Put(v)

		// 释放连接池中的所有连接
		// p.Release()
		// i++
		fmt.Println(i)
		// current := p.Len()
		// fmt.Println("len=", current)
	}

}
