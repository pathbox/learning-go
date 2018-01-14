package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	balance "../balance"
	pb "./pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	serv = flag.String("service", "hello_service", "service name")
	port = flag.Int("port", 50001, "listening port")
	reg  = flag.String("reg", "http://127.0.0.1:2379", "register etcd address")
)

type server struct{}

func main() {
	flag.Parse()

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		panic(err)
	}

	err = balance.Register(*serv, "127.0.0.1", *reg, *port, 15, time.Second*10)
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		log.Printf("receive signal %v", s)
		balance.UnRegister()
		os.Exit(1)
	}()
	log.Printf("starting hello service at %d", *port)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	s.Serve(listen)
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Printf("%v: Receive is %s\n", time.Now(), in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
