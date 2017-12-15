package main

import (
	"context"
	"log"
	"time"

	"github.com/coreos/etcd/client"
)

func main() {
	cfg := client.Config{
		Endpoints: []string{"http://localhost:2379"},
		Transport: client.DefaultTransport,
	}

	c, err := client.New(cfg)

	if err != nil {
		log.Fatalln(err)
	}

	KAPI := client.NewKeysAPI(c)

	// create
	_, err = KAPI.Create(context.Background(), "/foo1", "bar1")

	if err != nil {
		log.Println(err)
	}

	bg := context.Background()

	// 设置key值
	resp, err := KAPI.Get(bg, "/foo1", nil)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("value:", resp.Node.Value)
		log.Println("ttl:", resp.Node.TTL)
	}

	// 设置ttl
	opt := client.SetOptions{TTL: time.Duration(300) * time.Second}
	_, err = KAPI.Set(bg, "/foo1", "bar4", &opt)
	if err != nil {
		log.Println(err)
	}
	resp, err = KAPI.Get(bg, "/foo1", nil)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("value:", resp.Node.Value)
		log.Println("ttl:", resp.Node.TTL)
	}

	// 删除key
	_, err = KAPI.Delete(bg, "/foo1", nil)
	if err != nil {
		log.Println("del key failed:", err)
	}

	resp, err = KAPI.Get(bg, "/foo1", nil)
	if err != nil {
		log.Println("after del, error:", err)
	} else {
		log.Println("after del, value:", resp.Node.Value)
	}
}
