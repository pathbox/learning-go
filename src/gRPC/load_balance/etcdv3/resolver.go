package etcdv3

import (
    "errors"
    "fmt"
    "strings"

    etcd3 "github.com/coreos/etcd/clientv3"
    "google.golang.org/grpc/naming"
)

type resolver struct {
	serviceName string
}

// 命名解析实现


func NewResolver(serviceName string) *resolver {
	return &resolver{serviceName: serviceName}
}

// Resolve to resolve the service from etcd, target is the dial address of etcd
// target example: "http://127.0.0.1:2379,http://127.0.0.1:12379,http://127.0.0.1:22379"

func (re *resolver) Resolve(target string) (naming.Watcher, error) {
	if re.serviceName == ""{
		return nil, errors.New("grpclb: no service name provided")
	}

	client, err := etcd3.New(etcd3.Config{
		Endpoints: strings.Split(target, ","),
	})

	if err != nil {
		return nil, fmt.Errorf("grpclb: creat etcd3 client failed: %s", err.Error())
	}

	return &watcher{re: re, client: *client}, nil
}

