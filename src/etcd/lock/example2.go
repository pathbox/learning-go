package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
)

func main() {
	var name = flag.String("name", "", "give a name")
	flag.Parse()
	// Create a etcd client
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	// create a sessions to aqcuire a lock
	s, _ := concurrency.NewSession(cli)
	defer s.Close()
	l := concurrency.NewMutex(s, "/distributed-lock/")
	ctx := context.Background()
	// acquire lock (or wait to have it)
	if err := l.Lock(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("acquired lock for ", *name)
	fmt.Println("Do some work in", *name)
	time.Sleep(5 * time.Second)
	if err := l.Unlock(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("released lock for ", *name)
}
