package main

import (
	"context"
	"log"
	"sync"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

type ClientDis struct {
	client     *clientv3.Client
	serverList map[string]string // 客户端本地缓存 可用服务列表
	lock       sync.Mutex
}

func NewClientDis(addr []string) (*ClientDis, error) {
	conf := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	}
	if client, err := clientv3.New(conf); err == nil {
		return &ClientDis{
			client:     client,
			serverList: make(map[string]string),
		}, nil
	} else {
		return nil, err
	}
}

func (cliDis *ClientDis) GetService(prefix string) ([]string, error) {
	// 获取操作
	resp, err := cliDis.client.Get(context.Background(), prefix, clientv3.WithPrefix()) // 以前缀方式获取匹配
	if err != nil {
		return nil, err
	}
	addrs := cliDis.extractAddrs(resp)
	// 开启watch监听prefix
	go cliDis.watch(prefix)
	return addrs, nil
}

func (cliDis *ClientDis) watch(prefix string) {
	rch := cliDis.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for resp := range rch { //监听chan
		for _, ev := range resp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				cliDis.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				cliDis.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

// 处理resp，将得到的服务列表初始化存入 serverList缓存
func (cliDis *ClientDis) extractAddrs(resp *clientv3.GetResponse) []string {
	addrs := make([]string, 0)
	if resp == nil || resp.Kvs == nil {
		return addrs
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			cliDis.SetServiceList(string(resp.Kvs[i].Key), string(resp.Kvs[i].Value))
			addrs = append(addrs, string(v))
		}
	}
	return addrs
}

func (cliDis *ClientDis) SetServiceList(key, val string) {
	cliDis.lock.Lock()
	defer cliDis.lock.Unlock()
	cliDis.serverList[key] = val
	log.Println("set data key :", key, "val:", val)
}

func (cliDis *ClientDis) DelServiceList(key string) {
	cliDis.lock.Lock()
	defer cliDis.lock.Unlock()
	delete(cliDis.serverList, key)
	log.Println("del data key:", key)
}

func (cliDis *ClientDis) SerList2Array() []string {
	cliDis.lock.Lock()
	defer cliDis.lock.Unlock()
	addrs := make([]string, 0)

	for _, v := range cliDis.serverList {
		addrs = append(addrs, v)
	}
	return addrs
}

func main() {
	cli, _ := NewClientDis([]string{"127.0.0.1:2379"}) // 连接etcd得到客户端
	cli.GetService("/node")
}

/*
1.创建一个client 连到etcd。
2.匹配到所有相同前缀的 key。把值存到 serverList 这个map里面。
3.watch这个 key前缀，当有增加或者删除的时候 就 修改这个map。
4.所以这个map就是 实时的 服务列表
*/
