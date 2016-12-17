package main

import (
	"fmt"
	"os"
	//  "time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func main() {
	p, err := pool.New("tcp", "localhost:6379", 3)
	errHndlr(err)
	conn, err := p.Get()
	errHndlr(err)

	if conn.Cmd("SET", "foo", "fooval").Err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
	p.Put(conn)

	conn, err = p.Get()
	errHndlr(err)
	defer p.Put(conn)

	if conn.Cmd("SET", "foo", "fooval1").Err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	if conn.Cmd("GET", "foo").Err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	// execute a single command
	r := p.Cmd("GET", "foo")
	if r.Err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	// custom connections
	//  time.Sleep(60 * 1000 * time.Millisecond)

	df := func(network, addr string) (*redis.Client, error) {
		client, err := redis.Dial(network, addr)
		if err != nil {
			return nil, err
		}
		if err = client.Cmd("AUTH", "SUPERSECRET").Err; err != nil {
			client.Close()
			return nil, err
		}
		return client, nil
	}
	p, err = pool.NewCustom("tcp", "localhost:6379", 5, df)
}
