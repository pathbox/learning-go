package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "github.com/coreos/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	ss1, err := concurrency.NewSession(cli, concurrency.WithContext(ctx))
	if err != nil {
		log.Fatal(err)
	}
	defer ss1.Close()
	mu1 := concurrency.NewMutex(ss1, "/my-lock/")
	fmt.Println("Lock 1")
	fmt.Println(mu1.Lock(context.Background()))
	fmt.Println("Lock 2")

	doCancel := func() {
		fmt.Println("cancel 1")
		cancel()
		fmt.Println("cancel 2")
		mu1.Unlock(context.Background())
	}

	go func() {
		time.Sleep(5 * time.Second)
		doCancel()
	}()

	ss2, err := concurrency.NewSession(cli, concurrency.WithContext(ctx))
	if err != nil {
		log.Fatal(err)
	}
	defer ss2.Close()
	mu2 := concurrency.NewMutex(ss2, "/my-lock/")
	fmt.Println("Lock 3")
	fmt.Println(mu2.Lock(context.Background()))
	fmt.Println("Lock 4")
}
