package client

import (
	"log"

	pb "../../proto/cache"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	// zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

const (
	address = "localhost:50052"
)

func GetCache(ctx context.Context, tracer opentracing.Tracer, id int32) int32 {
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
	)

	defer conn.Close()
	c := pb.NewCacheClient(conn)

	r, err := c.Get(ctx, &pb.CacheRequest{Id: id})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return -1
	}

	return r.GetResult()
}
