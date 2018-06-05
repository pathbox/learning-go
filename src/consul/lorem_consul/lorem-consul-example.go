package main

import (
	"flag"
	"fmt"
	ilog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"context"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	ratelimitkit "github.com/go-kit/kit/ratelimit"
	"github.com/juju/ratelimit"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	lorem_consul "github.com/ru-rocker/gokit-playground/lorem-consul"
)

func main() {
	var (
		consulAddr    = flag.String("consul.addr", "", "consul address")
		consulPort    = flag.String("consul.port", "", "consul port")
		advertiseAddr = flag.String("advertise.addr", "", "advertise address")
		advertisePort = flag.String("advertise.port", "", "advertise port")
	)
	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// declare metrics
	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "namespace",
		Subsystem: "lorem_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryForm(stdprometheus.SummaryOpts{
		Namespace: "namespace",
		Subsystem: "lorem_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var svc lorem_consul.Service
	svc = lorem_consul.LoremService{}
	svc = lorem_consul.LoggingMiddleware(logger)(svc)
	svc = lorem_consul.Metrics(requestCount, requestLatency)(svc)

	rlbucket := ratelimit.NewBucket(1*time.Second, 5)
	e := lorem_consul.MakeLoremLoggingEndpoint(svc)
	e = ratelimitkit.NewTokenBucketThrottler(rlbucket, time.Sleep)(e)
	healthEndpoint := lorem_consul.MakeHealthEndpoint(svc)
	endpoint := lorem_consul.Endpoints{
		LoremEndpoint:  e,
		healthEndpoint: healthEndpoint,
	}

	r := lorem_consul.MakeHttpHandler(ctx, endpoint, logger)

	registar := lorem_consul.Register(*consulAddr,
		*consulPort,
		*advertiseAddr,
		*advertisePort,
	)

	go func() {
		ilog.Println("Starting server at port", *advertisePort)
		// register service
		registar.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+*advertisePort, handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errchan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan

	registar.Deregister() // 如果服务奔溃了，则取消注册
	ilog.Fatalln(error)

}
