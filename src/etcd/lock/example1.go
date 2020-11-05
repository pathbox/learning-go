package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
)

func main() {
	client, err := clientv3.New(clientv3.Config{Endpoints: []string{"http://127.0.0.1:2379"}, DialTimeout: 10 * time.Second})
	if err != nil {
		fmt.Errorf("client create fail - %v", err)
		return
	}
	defer client.Close()
	session, err := concurrency.NewSession(client, concurrency.WithTTL(10))
	if err != nil {
		fmt.Errorf("create session fail - %v", err)
		return
	}

	mutex := concurrency.NewMutex(session, "/lock")
	if mutex == nil {
		fmt.Errorf("create mutex fail")
		return
	}
	err = mutex.Lock(context.TODO())
	if err != nil {
		fmt.Errorf("lock fail - %v", err)
		return
	}
	fmt.Println("got lock, begin run work")

	go func() {
		select {
		case <-session.Done():
			// do what ever you want to process lock lost
			fmt.Println("lock lost")
		}
	}()

	go func() {
		// do real work here
	}()
	// prevent progress quit
	select {}
}

// 我们在一个 goroutine 中监听一个 <-session.Done() 的 channel ，这样，一旦锁出现了问题，就会得到通知，这样就可以在这里进行一些锁丢失的善后工作，比如在这里停止所有的需要锁才能进行的工作，这样就不会出现锁已经失效，但是工作进程却全然不知的状况了
