package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

// error codes
type Error int

const (
	errorSleep          = time.Second * 3 // retry when error.
	waitTimeout         = time.Second * 5 // timeout to get connection from pool.
	checkParentInterval = time.Second * 1 // the interval to check the parent pid.
	version             = "3.0.50"
	server              = "BMS/3.0.50(BOCAR)"
)

const (
	errorParseResponse Error = 100 + iota
	errorQueryRequest
	errorDetailRequest
	errorCreateRequest
	errorRemoveRequest
	errorFetchTimeout
	errorSystem
)

// trace level log
var Trace *log.Logger = log.New(os.Stdout, fmt.Sprintf("[%v][trace]", os.Getpid()), log.LstdFlags)
var Warn *log.Logger = log.New(os.Stdout, fmt.Sprintf("[%v][warn]", os.Getpid()), log.LstdFlags)
var Fatal *log.Logger = log.New(os.Stdout, fmt.Sprintf("[%v][error]", os.Getpid()), log.LstdFlags)

func (v Error) Error() string {
	return fmt.Sprintf("BOCAR ERROR code=%d", int(v))
}

func createRedisConnection(redisUri string) (c *redis.Client, err error) {
	defer func() {
		if r := recover(); r != nil {
			if err == nil {
				switch r := r.(type) {
				case error:
					err = r
				default:
					err = fmt.Errorf("fatal error: %v", r)
				}
			}
			Warn.Println("BOCAR recover from", r)
			return
		}
	}()

	var u *url.URL
	if u, err = url.Parse(redisUri); err != nil {
		Fatal.Println("BOCAR parse ", redisUri, "failed. err is", err)
		return
	}

	var conn *redis.Client
	if conn, err = redis.Dial("tcp", u.Host); err != nil {
		Fatal.Println("BOCAR dial", redisUri, "failed. err is", err)
		return
	}
	Trace.Println("BOCAR connected at", redisUri)

	return conn, err
}

func createRedisPoll(redisUri string, concurrency int, kick chan bool, pool chan *redis.Client) {
	// to generate number of connections.
	for i := 0; i < concurrency; i++ {
		kick <- true
	}

	// the connection generator
	go func() {
		for {
			<-kick
			for {
				c, err := createRedisConnection(redisUri)
				if err == nil {
					pool <- c
					break
				}
				Fatal.Println("BOCAR create redis failed. err is", err)
				time.Sleep(errorSleep)
			}
		}
	}()
}

type StandardResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

func NewErrorResponse(err error) *StandardResponse {
	if err, ok := err.(Error); ok {
		return &StandardResponse{Code: int(err)}
	}
	return &StandardResponse{Code: int(errorSystem)}
}
func NewStandardResponse(data interface{}) *StandardResponse {
	return &StandardResponse{Data: data}
}

func (v *StandardResponse) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", server)
	w.Write([]byte(v.String()))
}

func (v *StandardResponse) String() string {
	if b, err := json.Marshal(v); err == nil {
		return string(b)
	} else {
		v.Code = int(errorParseResponse)
		v.Data = fmt.Sprintf("parse data %v failed. err is %v", v.Data, err)
		if b, err := json.Marshal(v); err == nil {
			return string(b)
		}
		panic(err)
	}
}

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func main() {
	var nbCpus, concurrency int
	var listen, redisUri string

	flag.StringVar(&listen, "listen", ":2016", "the http server listen at")
	flag.StringVar(&redisUri, "redis", "redis://localhost:6379", "the uri of redis server")
	flag.IntVar(&nbCpus, "cpus", 1, "the cpus to use")
	flag.IntVar(&concurrency, "concurrency", 3, "the concurrency to redis, connection pool capacity.")

	flag.Usage = func() {
		fmt.Println(fmt.Sprintf("Usage: %v [--listen=string] [--redis=string] [--cpus=int] [--concurrency=int] [-h|--help]", os.Args[0]))
		fmt.Println(fmt.Sprintf(" listen, the listen [host]:port. default :2016"))
		fmt.Println(fmt.Sprintf(" redis, the host:port of redis server. default redis://localhost:6379"))
		fmt.Println(fmt.Sprintf(" cpus, the cpus to use. default 1"))
		fmt.Println(fmt.Sprintf(" concurrency, the concurrency connection to redis. default 3"))
		fmt.Println(fmt.Sprintf(" help, show this help and exit"))
		fmt.Println(fmt.Sprintf("@remark for redis uri, read https://www.iana.org/assignments/uri-schemes/prov/redis"))
		fmt.Println(fmt.Sprintf(" redis://user:secret@localhost:6379/0?foo=bar&qux=baz"))
		fmt.Println(fmt.Sprintf("For example:"))
		fmt.Println(fmt.Sprintf(" %v --listen=:2016 --redis=redis://localhost:6379", os.Args[0]))
	}
	flag.Parse()

	p, err := pool.New("tcp", "localhost:6379", 3)
	errHndlr(err)

	help := `query`
	Trace.Println(fmt.Sprintf("handle http://%v/api/v1/%v", listen, help))
	http.HandleFunc("/api/v1/query", func(w http.ResponseWriter, r *http.Request) {
		conn, err := p.Get()
		errHndlr(err)

		ret := conn.Cmd("GET", "foo")
		errHndlr(ret.Err)
		s, _ := ret.Str()
		Trace.Println(fmt.Sprintf("GET foo %v", s))

		NewStandardResponse(s).Handle(w, r)
	})

	if err := http.ListenAndServe(listen, nil); err != nil {
		Fatal.Println("BOCAR listen at", listen, "failed. err is", err)
		os.Exit(-1)
	}
}
