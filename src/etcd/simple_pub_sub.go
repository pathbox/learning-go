// 简单的发布/订阅的实现。集群中每台机器都用同一套配置文件，在启动阶段，都会到该节点上获取配置信息，同时客户端还需要在在节点注册一个数据变更的watcher监听，一旦配置文件发生变更，就会受到通知信息。
package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/coreos/etcd/client"
)

type App struct {
	Api        client.KeysAPI
	Key        string
	Config     string
	ConfigLock sync.RWMutex
}

func NewApp(api client.KeysAPI, key string) *App {
	app := &App{Api: api, Key: key}
	res, err := app.Api.Get(context.Background(), app.Key, nil)
	if err == nil {
		app.Config = res.Node.Value
	}

	go app.WatchConfig()
	return app
}

func (this *App) WatchConfig() {
	watcher := this.Api.Watcher(this.Key, &client.WatcherOptions{Recursive: true})
	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("watch config failed:", err)
			continue
		}
		if res.Action == "set" || res.Action == "update" {
			this.ConfigLock.Lock()
			this.Config = res.Node.Value
			this.ConfigLock.Unlock()
			log.Println("config update once")
		}
	}
}

func (this *App) Run() {
	for {
		time.Sleep(time.Duration(1) * time.Second)
		this.ConfigLock.RLock()
		log.Println("config:", this.Config)
		this.ConfigLock.RUnlock()
	}
}

func startApp() {
	cfg := client.Config{
		Endpoints: []string{"http://localhost:2379"},
		Transport: client.DefaultTransport,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatalln("create client failed:", err)
	}
	api := client.NewKeysAPI(c)
	app := NewApp(api, "config")
	app.Run()
}

type ConfManager struct {
	Api client.KeysAPI
}

func NewConfManager(api client.KeysAPI) *ConfManager {
	cm := &ConfManager{Api: api}
	return cm
}

func (this *ConfManager) UpdateConfig(key string, config string) error {
	_, err := this.Api.Set(context.Background(), key, config, nil)
	return err
}

func updateConfig() {
	cfg := client.Config{
		Endpoints: []string{"http://localhost:2379"},
		Transport: client.DefaultTransport,
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Fatalln("create client failed:", err)
	}
	api := client.NewKeysAPI(c)

	confManager := NewConfManager(api)

	var i int
	for i = 0; i < 100; i++ {
		time.Sleep(time.Duration(1) * time.Second)
		confManager.UpdateConfig("config", fmt.Sprintf("%d", i))
	}
}

func main() {
	go startApp()
	go updateConfig()
	select {}
}
