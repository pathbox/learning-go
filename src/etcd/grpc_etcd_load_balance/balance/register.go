package balance

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"

	etcd3 "github.com/coreos/etcd/clientv3"
)

var Prefix = "etcd3_name"
var client etcd3.Client
var serviceKey string

var stopSignal = make(chan bool, 1)

// Register
// Register 就是Put key=>value 到etcd. key为有意义的唯一的key,value为实例服务器的地址.
// 并不是简单的注册一次,而是在后台守护进程,心跳的模式对etcd进行注册
func Register(name, host, target string, port, ttl int, interval time.Duration) error {
	serviceValue := fmt.Sprintf("%s:%d", host, port)
	serviceKey = fmt.Sprintf("/%s/%s/%s", Prefix, name, serviceValue)

	// get endpoints for register dial address
	var err error
	cfg := etcd3.Config{
		Endpoints: strings.Split(target, ","),
	}
	client, err := etcd3.New(cfg)
	if err != nil {
		return fmt.Errorf("grpclb: create etcd3 client failed: %v", err)
	}

	go func() {
		// invoke self-register with ticker
		ticker := time.NewTicker(interval)
		for {
			resp, _ := client.Grant(context.Background(), int64(ttl))
			// should get first, if not exist, set it
			_, err := client.Get(context.Background(), serviceKey)
			if err != nil {
				if err == rpctypes.ErrKeyNotFound {
					if _, err := client.Put(context.TODO(), serviceKey, serviceValue, etcd3.WithLease(resp.ID)); err != nil {
						log.Printf("grpclb: set service '%s' with ttl to etcd3 failed: %s", name, err.Error())
					} else {
						log.Printf("grpclb: service '%s' connect to etcd3 failed: %s", name, err.Error())
					}
				} else {
					// refresh set to true for not notifying the watcher
					if _, err := client.Put(context.Background(), serviceKey, serviceValue, etcd3.WithLease(resp.ID)); err != nil {
						log.Printf("grpclb: refresh service '%s' with ttl to etcd3 failed: %s", name, err.Error())
					}
				}

				// 使用select 控制for循环,中止或按照ticker间隔进行. 这样就不用sleep这样的方法了
				select {
				case <-stopSignal:
					return
				case <-ticker.C:
				}
			}
		}
	}()

	return nil
}

// UnRegister delete registered service from etcd
func UnRegister() error {
	stopSignal <- true
	stopSignal = make(chan bool, 1) // just a hack to avoid multi UnRegister deadlock
	var err error
	if _, err := client.Delete(context.Background(), serviceKey); err != nil {
		log.Printf("grpclb: deregister '%s' failed: %s", serviceKey, err.Error())
	} else {
		log.Printf("grpclb: deregister '%s' ok.", serviceKey)
	}
	return err
}
