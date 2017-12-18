package main

import (
	"log"
	"net"
	"time"

	"github.com/opentracing/opentracing-go"

	zipkin "github.com/openzipkin/zipkin-go-opentracing"

	pb "../proto/cache"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

type CacheServer struct{}

func (s *CacheServer) Get(ctx context.Context, in *pb.CacheRequest) (*pb.CacheReply, error) {
	log.Printf("input %d", in.GetId())
	time.Sleep(time.Duration(50) * time.Millisecond)

	return &pb.CacheReply{Result: in.GetId() * 2}, nil

}

func main() {

	collector, err := zipkin.NewHTTPCollector("http://localhost:9411/api/v1/spans")
	if err != nil {
		log.Fatal(err)
		return
	}

	tracer, err := zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, "localhost:0", "grpc_cache"),
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
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads())))
	pb.RegisterCacheServer(s, &CacheServer{})
	log.Println("Cache Server Listening ", port)
	s.Serve(lis)
}
