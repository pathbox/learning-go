package grpcsr

import (
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/consul/api"
)

func NewConsulRegister() *ConsulRegister {
	return &ConsulRegister{
		Address: "127.0.0.1:8500",
		Service: "example",
		Tag:     []string{},
		Port:    3000,
		DeregisterCriticalServiceAfter: time.Duration(1) * time.Minute,
		Interval:                       time.Duration(10) * time.Second,
	}
}

type ConsulRegister struct {
	Address                        string
	Service                        string
	Tag                            []string
	Port                           int
	DeregisterCriticalServiceAfter time.Duration
	Interval                       time.Duration
}

func (r *ConsulRegister) Register() error {
	config := api.DefaultConfig()
	config.Address = r.Address
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	agent := client.Agent()

	IP := localIP()
	reg := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v-%v", r.Service, IP, r.Port), // 服务节点的名称 将服务节点的ip和port注册到 consul，consul会对此进行health check
		Name:    fmt.Sprintf("grpc.health.v1.%v", r.Service),    // 服务名称
		Tags:    r.Tag,                                          // tag，可以为空
		Port:    r.Port,                                         // 服务端口
		Address: IP,                                             // 服务 IP
		Check: &api.AgentServiceCheck{ // 健康检查
			Interval: r.Interval.String(),                            // 健康检查间隔
			GRPC:     fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Service), // grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			DeregisterCriticalServiceAfter: r.DeregisterCriticalServiceAfter.String(), // 注销时间，相当于过期时间
		},
	}
	if err := agent.ServiceRegister(reg); err != nil {
		return err
	}
	return nil
}

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
