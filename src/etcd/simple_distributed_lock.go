package main

import (
	"log"
	"time"

	"context"

	"github.com/coreos/etcd/client"
)

type Worker struct {
	Hostname      string
	Endpoints     []string
	Api           client.KeysAPI
	Key           string
	SleepDuration time.Duration
}

func NewWorker(hostname string, endpoints []string, key string, duration time.Duration) (*Worker, error) {
	cfg := client.Config{Endpoints: endpoints, Transport: client.DefaultTransport}
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	api := client.NewKeysAPI(c)
	w := &Worker{Hostname: hostname, Endpoints: endpoints, Key: key, Api: api, SleepDuration: duration}
	return w, nil
}

func (this *Worker) Work() {
	duration := time.Duration(10) * time.Second
	for {
		_, err := this.Api.Set(context.Background(), this.Key, "", &client.SetOptions{PrevExist: "false", TTL: duration})
		if err != nil {
			continue
		}
		log.Println("hostname:", this.Hostname, "lock succed")
		time.Sleep(this.SleepDuration)
		log.Println("hostname:", this.Hostname, "exec finish")
		this.Api.Delete(context.Background(), this.Key, nil)
	}
}

func main() {
	endpoints := []string{"http://localhost:2379"}
	w1, _ := NewWorker("hostname1", endpoints, "lock", time.Duration(1)*time.Second)
	w2, _ := NewWorker("hostname2", endpoints, "lock", time.Duration(2)*time.Second)
	w3, _ := NewWorker("hostname3", endpoints, "lock", time.Duration(3)*time.Second)
	go w1.Work()
	go w2.Work()
	go w3.Work()
	select {}
}
