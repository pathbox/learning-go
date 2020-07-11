package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

//创建租约注册服务
type ServiceReg struct {
	client        *clientv3.Client
	lease         clientv3.Lease
	leaseResp     *clientv3.LeaseGrantResponse
	canclefunc    func()
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
}

func NewServiceReg(addr []string, timeNum int64) (*ServiceReg, error) {
	conf := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	}

	var (
		client *clientv3.Client
	)

	if clientTem, err := clientv3.New(conf); err == nil {
		client = clientTem
	} else {
		return nil, err
	}

	serReg := &ServiceReg{
		client: client,
	}
	if err := serReg.setLease(timeNum); err != nil {
		return nil, err
	}
	go serReg.ListenLeaseRespChan()
	return serReg, nil
}

//设置租约
func (this *ServiceReg) setLease(timeNum int64) error {
	lease := clientv3.NewLease(this.client)

	//设置租约时间
	leaseResp, err := lease.Grant(context.Background(), timeNum)
	if err != nil {
		return err
	}

	//设置续租
	ctx, cancelFunc := context.WithCancel(context.Background())
	leaseRespChan, err := lease.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return err
	}

	this.lease = lease
	this.leaseResp = leaseResp
	this.canclefunc = cancelFunc
	this.keepAliveChan = leaseRespChan
	return nil
}

//监听 续租情况
func (this *ServiceReg) ListenLeaseRespChan() {
	for {
		select {
		case leaseKeepResp := <-this.keepAliveChan:
			if leaseKeepResp == nil {
				fmt.Printf("已经关闭续租功能\n")
				return
			} else {
				fmt.Printf("续租成功\n")
			}
		}
	}
}

//通过租约 注册服务
func (this *ServiceReg) PutService(key, val string) error {
	kv := clientv3.NewKV(this.client)
	_, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(this.leaseResp.ID))
	return err
}

func main() {
	ser, _ := NewServiceReg([]string{"127.0.0.1:2379"}, 5) // 127.0.0.1:2379 etcd的服务地址
	ser.PutService("/node/001", "127.0.0.1:1212")
	select {}
}

/*
etcd的 租约模式:客户端申请 一个租约 并设置 过期时间，每隔一段时间 就要 请求 etcd 申请续租。客户端可以通过租约存key。如果不续租 ，过期了，etcd 会删除这个租约上的 所有key-value。类似于心跳模式。

一般相同的服务存的 key 的前缀是一样的 比如 “/node/001"=> "127.0.0.1:1212" 和 ”server/002"=>"127.0.0.1:1313" 这种模式，然后 客户端 就直接 匹配 “server/” 这个key
*/
