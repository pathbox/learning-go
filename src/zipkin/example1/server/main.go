package main

import (
	"log"
	"net"

	"github.com/opentracing/opentracing-go"

	zipkin "github.com/openzipkin/zipkin-go-opentracing"

	cache "../cache/client"
	pb "../proto/add"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type AddServer struct{}

func (s *AddServer) DoAdd(ctx context.Context, in *pb.AddRequest) (*pb.AddReply, error) {
	log.Printf("input %d %d", in.GetNum1(), in.GetNum2())

	tracer := opentracing.GlobalTracer()
	val := cache.GetCache(ctx, tracer, in.GetNum1())
	log.Printf("cache value %d", val)

	return &pb.AddReply{Result: val + in.GetNum2()}, nil
}

func main() {
	collector, err := zipkin.NewHTTPCollector("http://localhost:9411/api/v1/spans")
	if err != nil {
		log.Fatal(err)
		return
	}

	tracer, err := zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, "localhost:0", "grpc_server"),
		zipkin.ClientServerSameSpan(true),
		zipkin.TraceID128Bit(true),
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	opentracing.InitGlobalTracer(tracer)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(
				tracer,
				otgrpc.LogPayloads(),
			),
		),
	)
	pb.RegisterAddServer(s, &AddServer{})
	log.Println("Server is Listening ", port)
	s.Serve(lis)

}
