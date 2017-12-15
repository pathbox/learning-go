package main

import (
	"context"
	"log"
	"time"

	"github.com/coreos/etcd/client"
)

type Worker struct {
	Addr     string
	Api      client.KeysAPI
	Key      string
	Interval time.Duration
	TTL      time.Duration
}

func NewWorker(addr string, api client.KeysAPI, key string, interval time.Duration, ttl time.Duration) *Worker {
	w := &Worker{Addr: addr, Api: api, Key: key, Interval: interval, TTL: ttl}
	go w.HeartBeat()
	return w
}

func (w *Worker) HeartBeat() { // 服务注册
	for {
		Key := w.Key + w.Addr
		_, err := w.Api.Set(context.Background(), Key, w.Addr, &client.SetOptions{}) // 对etcd 进行 Set操作
		if err != nil {
			log.Println("worker set failed:", err)
		} else {
			log.Println("worker set once")
		}
		time.Sleep(w.Interval)
	}
}

func (w *Worker) Work() {
	for {
		// do something
		time.Sleep(w.Interval)
	}
}

type Master struct {
	M        map[string]string
	Api      client.KeysAPI
	Key      string
	Interval time.Duration
}

func NewMaster(api client.KeysAPI, key string, interval time.Duration) *Master {
	m := &Master{M: map[string]string{}, Api: api, Key: key, Interval: interval}
	go m.WatchWorkers()
	return m
}

func (this *Master) WatchWorkers() { // 服务发现
	watcher := this.Api.Watcher(this.Key, &client.WatcherOptions{Recursive: true})
	for {
		time.Sleep(this.Interval)
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("master watch worker failed:", err)
			continue
		}
		value := res.Node.Value
		var prevalue string
		if res.PrevNode != nil {
			prevalue = res.PrevNode.Value
		}
		log.Println("master new event, action:", res.Action, "value:", value, "prevalue:", prevalue)
		if res.Action == "set" {
			log.Println("Set Action") // 走的是这一步
			this.M[value] = ""
		}
		if res.Action == "update" {
			delete(this.M, prevalue)
			this.M[value] = ""
		}
		if res.Action == "delete" || res.Action == "expire" {
			delete(this.M, prevalue)
		}
	}
}

func (this *Master) PrintWorkers() {
	//log.Println("workers:", this.WorkerAddrs)
	for {
		time.Sleep(time.Duration(1) * time.Second)
		log.Println(this.M)
	}
}

func startWorker() {
	cfg := client.Config{Endpoints: []string{"http://localhost:2379"}, Transport: client.DefaultTransport}
	wrokerClient, err := client.New(cfg)
	if err != nil {
		log.Fatalln("create worker client failed:", err)
	}
	wrokerApi := client.NewKeysAPI(wrokerClient)
	interval := time.Duration(1) * time.Second
	ttl := time.Duration(10) * time.Second
	worker := NewWorker("127.0.0.1:80", wrokerApi, "workers/", interval, ttl)
	worker.Work()
}

func startMaster() {
	cfg := client.Config{Endpoints: []string{"http://localhost:2379"}, Transport: client.DefaultTransport}
	masterClient, err := client.New(cfg)
	if err != nil {
		log.Fatalln("create master client failed:", err)
	}
	masterApi := client.NewKeysAPI(masterClient)
	interval := time.Duration(0) * time.Second
	m := NewMaster(masterApi, "workers/", interval)
	m.PrintWorkers()
}

func main() {
	go startWorker()
	go startMaster()
	select {}
}

// 代码的思路很简单, worker启动时向etcd注册自己的信息,并设置一个过期时间TTL,每隔一段时间更新这个TTL,如果该worker挂掉了,这个TTL就会expire.
// master则监听 workers/ 这个etcd directory, 根据检测到的不同action来增加, 更新, 或删除worker
