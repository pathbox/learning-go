package main

import (
	"log"

	pb "../proto/add"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

const (
	address = "localhost:50051"
)

func main() {
	num1 := 10
	num2 := 20

	collector, err := zipkin.NewHTTPCollector("http://localhost:9411/api/v1/spans")
	if err != nil {
		log.Fatal(err)
		return
	}

	tracer, err := zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, "localhost:0", "grpc_client"),
		zipkin.ClientServerSameSpan(true),
		zipkin.TraceID128Bit(true),
	)

	if err != nil {
		log.Fatal(err)
		return
	}
	opentracing.InitGlobalTracer(tracer)
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAddClient(conn)

	span := opentracing.StartSpan("Start")

	ctx := opentracing.ContextWithSpan(context.Background(), span)
	r, err := c.DoAdd(ctx, &pb.AddRequest{Num1: int32(num1), Num2: int32(num2)})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("add(%d,%d), Result: %d", num1, num2, r.GetResult())

	span.Finish()
	collector.Close()
}

// curl -sSL https://zipkin.io/quickstart.sh | bash -s
// java -jar zipkin.jar --logging.level.zipkin=DEBUG --logging.level.zipkin2=DEBUG

/*
启动两个grpc服务
go run cache/main.go
go run server/main.go

go run client/main.go

*/

/*
文档资料
https://github.com/openzipkin/zipkin/tree/master/zipkin-server
https://github.com/opentracing/specification/blob/master/project_organization.md
*/
// http://localhost:9411/zipkin
