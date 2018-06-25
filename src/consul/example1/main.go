package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	addservice "github.com/hatlonely/hellogolang/sample/addservice/api"
	"github.com/hatlonely/hellogolang/sample/addservice/internal/grpcsr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// AddServiceImpl 实现 Add 服务
type AddServiceImpl struct{}

// Add 接口实现
func (s *AddServiceImpl) Add(ctx context.Context, request *addservice.AddRequest) (*addservice.AddResponse, error) {
	// 50% 概率 sleep，模拟超时场景
	if rand.Int()%2 == 0 {
		time.Sleep(time.Duration(200) * time.Millisecond)
	}
	// fmt.Println(request)
	response := &addservice.AddResponse{
		V: request.A + request.B,
	}
	logrus.WithField("request", request).WithField("response", response).Info()
	return response, nil
}

type HealthImpl struct{}

// Check 实现健康检查接口，这里直接返回健康状态，这里也可以有更复杂的健康检查策略，比如根据服务器负载来返回
func (h *HealthImpl) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func main() {
	port := pflag.IntP("register.port", "p", 3000, "service port")
	pflag.Parse()

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	entry := logrus.NewEntry(logger)

	grpc_logrus.ReplaceGrpcLogger(entry)

	server := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(
				grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
			),
			grpc_logrus.UnaryServerInterceptor(
				entry,
				grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
					return "grpc.time_ns", duration.Nanoseconds()
				}),
			),
			grpc_logrus.PayloadUnaryServerInterceptor(entry, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool { return true }),
		),
	)

	addservice.RegisterAddServiceServer(server, &AddServiceImpl{})
	grpc_health_v1.RegisterHealthServer(server, &HealthImpl{})

	// 从 consul 读取配置文件
	config := viper.New()
	config.AddRemoteProvider("consul", "127.0.0.1:8500", "config/addservice.json")
	config.SetConfigType("json")
	if err := config.ReadRemoteConfig(); err != nil {
		panic(err)
	}
	config.BindPFlags(pflag.CommandLine)

	// 使用 consul 注册服务
	register := grpcsr.NewConsulRegister()
	config.Sub("register").Unmarshal(register)
	register.Port = *port
	if err := register.Register(); err != nil {
		panic(err)
	}

	address, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", *port))
	if err != nil {
		panic(err)
	}

	if err := server.Serve(address); err != nil {
		panic(err)
	}

}
