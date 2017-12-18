package main

import (
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"

	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

func main() {
	// 1) Create a opentracing.Tracer that sends data to Zipkin
	collector, _ := zipkin.NewHTTPCollector(
		fmt.Sprintf("http://%s:9411/api/v1/spans", "localhost"))

	tracer, _ := zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, "127.0.0.1:0", "trivial"),
	)

	// 2) Demonstrate simple OpenTracing instrumentation
	parent := tracer.StartSpan("Parent")

	for i := 0; i < 20; i++ {
		parent.LogEvent(fmt.Sprintf("Starting child #%d", i))
		child := tracer.StartSpan("Child", opentracing.ChildOf(parent.Context()))
		time.Sleep(50 * time.Millisecond)
		child.Finish()
	}
	parent.LogEvent("A Log")
	parent.Finish()

	// ... give Zipkin ample time to flush
	time.Sleep(2 * time.Second)
}
